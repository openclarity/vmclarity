// Copyright © 2023 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ghodss/yaml"
	kubeclarityutils "github.com/openclarity/kubeclarity/shared/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/cli/pkg"
	cicdinitiator "github.com/openclarity/vmclarity/cli/pkg/cicd/initiator"
	"github.com/openclarity/vmclarity/cli/pkg/cli"
	"github.com/openclarity/vmclarity/cli/pkg/presenter"
	"github.com/openclarity/vmclarity/cli/pkg/state"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/utils"
	"github.com/openclarity/vmclarity/shared/pkg/backendclient"
	"github.com/openclarity/vmclarity/shared/pkg/families"
	"github.com/openclarity/vmclarity/shared/pkg/families/malware"
	misconfigurationTypes "github.com/openclarity/vmclarity/shared/pkg/families/misconfiguration/types"
	"github.com/openclarity/vmclarity/shared/pkg/families/rootkits"
	"github.com/openclarity/vmclarity/shared/pkg/families/sbom"
	"github.com/openclarity/vmclarity/shared/pkg/families/secrets"
	"github.com/openclarity/vmclarity/shared/pkg/families/vulnerabilities"
)

const DefaultWatcherInterval = 2 * time.Minute

var (
	cfgFile string
	config  *families.Config
	logger  *logrus.Entry
	output  string

	server                string
	scanResultID          string
	mountVolume           bool
	waitForServerAttached bool

	// flags for CICD mode.
	cicdMode          bool
	scanConfigName    string
	scanConfigID      string
	input             string
	inputType         string
	exportCICDResults bool
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:          "vmclarity",
	Short:        "VMClarity",
	Long:         `VMClarity`,
	Version:      pkg.GitRevision,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Infof("Running...")

		// Main context which remains active even if the scan is aborted allowing post-processing operations
		// like updating scan result state
		ctx := cmd.Context()

		cli, err := newCli()
		if err != nil {
			return fmt.Errorf("failed to initialize CLI: %w", err)
		}

		// Create context used to signal to operations that the scan is aborted
		abortCtx, cancel := context.WithCancel(ctx)
		defer cancel()

		// Start watching for abort event
		cli.WatchForAbort(ctx, cancel, DefaultWatcherInterval)

		if waitForServerAttached {
			if err := cli.WaitForVolumeAttachment(abortCtx); err != nil {
				err = fmt.Errorf("failed to wait for block device being attached: %w", err)
				if e := cli.MarkDone(ctx, []error{err}); e != nil {
					logger.Errorf("Failed to update scan result stat to completed with errors: %v", e)
				}
				return err
			}
		}

		if mountVolume {
			mountPoints, err := cli.MountVolumes(abortCtx)
			if err != nil {
				err = fmt.Errorf("failed to mount attached volume: %w", err)
				if e := cli.MarkDone(ctx, []error{err}); e != nil {
					logger.Errorf("Failed to update scan result stat to completed with errors: %v", e)
				}
				return err
			}
			setMountPointsForFamiliesInput(mountPoints, config)
		}

		if input != "" {
			famInputType, err := getFamiliesInputType(inputType)
			if err != nil {
				return fmt.Errorf("failed to get families input type by inputType=%s: %w", inputType, err)
			}
			appendInput(input, famInputType, config)
		}

		err = cli.MarkInProgress(ctx)
		if err != nil {
			return fmt.Errorf("failed to inform server %v scan has started: %w", server, err)
		}

		logger.Infof("Running scanners...")
		res, familiesErr := families.New(logger, config).Run(abortCtx)

		logger.Infof("Exporting results...")
		if scanResultID == "" {
			// In the case of standalone mode if the scanResultID is not set
			// we get it from the status manager.
			cli.SetScanResultID(cli.GetScanResultID())
		}
		errs := cli.ExportResults(abortCtx, res, familiesErr)

		if len(familiesErr) > 0 {
			errs = append(errs, fmt.Errorf("at least one family failed to run"))
		}

		err = cli.MarkDone(ctx, errs)
		if err != nil {
			return fmt.Errorf("failed to inform the server %v the scan was completed: %w", server, err)
		}

		if len(familiesErr) > 0 {
			return fmt.Errorf("failed to run families: %+v", familiesErr)
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

// nolint: gochecknoinits
func init() {
	cobra.OnInitialize(
		initLogger,
		initConfig,
	)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.vmclarity.yaml)")
	rootCmd.PersistentFlags().StringVar(&output, "output", "", "set output directory path. Stdout is used if not set")
	rootCmd.PersistentFlags().StringVar(&server, "server", "", "VMClarity server to export scan results to, for example: http://localhost:9999/api")
	rootCmd.PersistentFlags().StringVar(&scanConfigName, "scan-config-name", "", "use an existing scan config that is defined in the VMClarity server")
	rootCmd.PersistentFlags().StringVar(&scanResultID, "scan-result-id", "", "the ScanResult ID to export the scan results to")
	rootCmd.PersistentFlags().BoolVar(&mountVolume, "mount-attached-volume", false, "discover for an attached volume and mount it before the scan")
	rootCmd.PersistentFlags().BoolVar(&waitForServerAttached, "wait-for-server-attached", false, "wait for the VMClarity server to attach the volume")
	rootCmd.PersistentFlags().StringVar(&input, "input", "", "input for families")
	rootCmd.PersistentFlags().StringVar(&inputType, "input-type", "dir", "input type for families")
	rootCmd.PersistentFlags().BoolVar(&cicdMode, "cicd-mode", false, "CICD mode")
	rootCmd.PersistentFlags().BoolVar(&exportCICDResults, "export-cicd-results", false, "export results to VMclarity server")

	validateRequiredFlagForDefinedFlag(rootCmd, "scan-result-id", "server")
	validateRequiredFlagForDefinedFlag(rootCmd, "scan-config-name", "server")
	rootCmd.MarkFlagsMutuallyExclusive("config", "scan-config-name")
	rootCmd.MarkFlagsMutuallyExclusive("mount-attached-volume", "input")
	validateRequiredFlagForDefinedFlag(rootCmd, "input-type", "input")
	validateRequiredFlagForDefinedFlag(rootCmd, "export-cicd-results", "server")
}

func validateRequiredFlagForDefinedFlag(rootCmd *cobra.Command, definedFlag, requiredFlag string) {
	if flag := rootCmd.Flag(definedFlag); flag != nil {
		if required := rootCmd.Flag(requiredFlag); required == nil {
			logrus.Fatalf("Cannot set flag '%s' alone without flag '%s'", definedFlag, requiredFlag)
		}
	}
}

func getConfigFromBackend() *families.Config {
	if server == "" {
		panic("Missing backend")
	}
	client, err := backendclient.Create(server)
	if err != nil {
		logrus.Fatalf("failed to create VMClarity API client: %v", err)
	}

	scanConfigs, err := client.GetScanConfigs(context.TODO(), models.GetScanConfigsParams{
		Filter: utils.PointerTo(fmt.Sprintf("name eq '%s'", scanConfigName)),
	})
	if err != nil {
		logrus.Fatalf("Failed to get scan config by name %v", err)
	}
	if len(*scanConfigs.Items) == 0 {
		logrus.Fatalf("There is no scan config with name=%s", scanConfigName)
	}

	scanConfig := families.CreateFamilyConfigFromModel(
		(*scanConfigs.Items)[0].ScanFamiliesConfig,
		families.LoadAddresses("localhost"),
		families.LoadPaths(),
	)
	scanConfigID = *(*scanConfigs.Items)[0].Id

	return &scanConfig
}

func getConfigFromFile() *families.Config {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory OR current directory with name ".families" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".families")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	cobra.CheckErr(err)

	// Load config
	config = &families.Config{}
	err = viper.Unmarshal(config)
	cobra.CheckErr(err)

	return config
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	logrus.Infof("init config")

	if scanConfigName != "" {
		config = getConfigFromBackend()
	} else {
		config = getConfigFromFile()
	}

	if logrus.IsLevelEnabled(logrus.InfoLevel) {
		configB, err := yaml.Marshal(config)
		cobra.CheckErr(err)
		if scanConfigName != "" {
			logrus.Infof("Using config from backend (%s):\n%s", scanConfigName, string(configB))
		} else {
			logrus.Infof("Using config file (%s):\n%s", viper.ConfigFileUsed(), string(configB))
		}
	}
}

func initLogger() {
	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)
	logger = log.WithField("app", "vmclarity")
}

// nolint:gocognit,cyclop
func newCli() (*cli.CLI, error) {
	var manager state.Manager
	var presenters []presenter.Presenter

	var err error

	if config == nil {
		return nil, errors.New("families config must not be nil")
	}

	//	if (server != "" && !cicdMode) || (cicdMode && exportCICDResults) {
	if server != "" {
		var client *backendclient.BackendClient
		var p presenter.Presenter

		client, err = backendclient.Create(server)
		if err != nil {
			return nil, fmt.Errorf("failed to create VMClarity API client: %w", err)
		}

		if cicdMode {
			if exportCICDResults {
				cicdInitiatorConfig := cicdinitiator.CreateConfig(client, config, scanConfigID, scanConfigName, input, inputType)
				manager, err = state.NewCICDState(client, scanResultID, cicdInitiatorConfig)
				if err != nil {
					return nil, fmt.Errorf("failed to create CICD state manager: %w", err)
				}

				p, err = presenter.NewVMClarityPresenter(client, scanResultID)
				if err != nil {
					return nil, fmt.Errorf("failed to create VMClarity presenter: %w", err)
				}
			}
		} else {
			manager, err = state.NewVMClarityState(client, scanResultID)
			if err != nil {
				return nil, fmt.Errorf("failed to create VMClarity state manager: %w", err)
			}

			p, err = presenter.NewVMClarityPresenter(client, scanResultID)
			if err != nil {
				return nil, fmt.Errorf("failed to create VMClarity presenter: %w", err)
			}
		}
		if p != nil {
			presenters = append(presenters, p)
		}
	} else {
		manager, err = state.NewLocalState()
		if err != nil {
			return nil, fmt.Errorf("failed to create local state: %w", err)
		}
	}

	if output != "" {
		presenters = append(presenters, presenter.NewFilePresenter(output, config))
	} else {
		presenters = append(presenters, presenter.NewConsolePresenter(os.Stdout, config))
	}

	var p presenter.Presenter
	if len(presenters) == 1 {
		p = presenters[0]
	} else {
		p = &presenter.MultiPresenter{Presenters: presenters}
	}

	return &cli.CLI{Manager: manager, Presenter: p, FamiliesConfig: config}, nil
}

func setMountPointsForFamiliesInput(mountPoints []string, familiesConfig *families.Config) *families.Config {
	// update families inputs with the mount point as rootfs
	for _, mountDir := range mountPoints {
		familiesConfig = appendInput(mountDir, string(kubeclarityutils.ROOTFS), familiesConfig)
	}
	return familiesConfig
}

func appendInput(input, inputType string, familiesConfig *families.Config) *families.Config {
	if familiesConfig.SBOM.Enabled {
		familiesConfig.SBOM.Inputs = append(familiesConfig.SBOM.Inputs, sbom.Input{
			Input:     input,
			InputType: inputType,
		})
	}

	if familiesConfig.Vulnerabilities.Enabled {
		if familiesConfig.SBOM.Enabled {
			familiesConfig.Vulnerabilities.InputFromSbom = true
		} else {
			familiesConfig.Vulnerabilities.Inputs = append(familiesConfig.Vulnerabilities.Inputs, vulnerabilities.Input{
				Input:     input,
				InputType: inputType,
			})
		}
	}

	if familiesConfig.Secrets.Enabled {
		familiesConfig.Secrets.Inputs = append(familiesConfig.Secrets.Inputs, secrets.Input{
			StripPathFromResult: utils.PointerTo(true),
			Input:               input,
			InputType:           inputType,
		})
	}

	if familiesConfig.Malware.Enabled {
		familiesConfig.Malware.Inputs = append(familiesConfig.Malware.Inputs, malware.Input{
			StripPathFromResult: utils.PointerTo(true),
			Input:               input,
			InputType:           inputType,
		})
	}

	if familiesConfig.Rootkits.Enabled {
		familiesConfig.Rootkits.Inputs = append(familiesConfig.Rootkits.Inputs, rootkits.Input{
			StripPathFromResult: utils.PointerTo(true),
			Input:               input,
			InputType:           inputType,
		})
	}

	if familiesConfig.Misconfiguration.Enabled {
		familiesConfig.Misconfiguration.Inputs = append(
			familiesConfig.Misconfiguration.Inputs,
			misconfigurationTypes.Input{
				StripPathFromResult: utils.PointerTo(true),
				Input:               input,
				InputType:           inputType,
			},
		)
	}

	return familiesConfig
}

func getFamiliesInputType(inputType string) (string, error) {
	switch inputType {
	case "dir", "DIR":
		return string(kubeclarityutils.DIR), nil
	case "vm", "VM":
		return string(kubeclarityutils.ROOTFS), nil
	default:
		return "", errors.New("input type is not supported")
	}
}
