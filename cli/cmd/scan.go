/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ghodss/yaml"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	kubeclarityutils "github.com/openclarity/kubeclarity/shared/pkg/utils"

	"github.com/openclarity/vmclarity/cli/pkg/cli"
	"github.com/openclarity/vmclarity/cli/pkg/presenter"
	"github.com/openclarity/vmclarity/cli/pkg/state"
	"github.com/openclarity/vmclarity/shared/pkg/backendclient"
	"github.com/openclarity/vmclarity/shared/pkg/families"
	"github.com/openclarity/vmclarity/shared/pkg/families/malware"
	misconfigurationTypes "github.com/openclarity/vmclarity/shared/pkg/families/misconfiguration/types"
	"github.com/openclarity/vmclarity/shared/pkg/families/rootkits"
	"github.com/openclarity/vmclarity/shared/pkg/families/sbom"
	"github.com/openclarity/vmclarity/shared/pkg/families/secrets"
	"github.com/openclarity/vmclarity/shared/pkg/families/vulnerabilities"
	"github.com/openclarity/vmclarity/shared/pkg/log"
	"github.com/openclarity/vmclarity/shared/pkg/utils"
)

const (
	DefaultWatcherInterval = 2 * time.Minute
	DefaultMountTimeout    = 10 * time.Minute
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan",
	Long:  `Run scanner families`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Infof("Running...")

		// Main context which remains active even if the scan is aborted allowing post-processing operations
		// like updating asset scan state
		ctx := log.SetLoggerForContext(cmd.Context(), logger)

		cfgFile, err := cmd.Flags().GetString("config")
		if err != nil {
			logger.Fatalf("Unable to get asset json file name: %v", err)
		}
		server, err := cmd.Flags().GetString("server")
		if err != nil {
			logger.Fatalf("Unable to get VMClarity server address: %v", err)
		}
		output, err := cmd.Flags().GetString("output")
		if err != nil {
			logger.Fatalf("Unable to get output file name: %v", err)
		}
		assetScanID, err := cmd.Flags().GetString("asset-scan-id")
		if err != nil {
			logger.Fatalf("Unable to get asset scan ID: %v", err)
		}
		mountVolume, err := cmd.Flags().GetBool("mount-attached-volume")
		if err != nil {
			logger.Fatalf("Unable to get mount attached volume flag: %v", err)
		}

		config := loadConfig(cfgFile)
		cli, err := newCli(config, server, assetScanID, output)
		if err != nil {
			return fmt.Errorf("failed to initialize CLI: %w", err)
		}

		// Create context used to signal to operations that the scan is aborted
		abortCtx, cancel := context.WithCancel(ctx)
		defer cancel()

		// Start watching for abort event
		cli.WatchForAbort(ctx, cancel, DefaultWatcherInterval)

		if err := cli.WaitForReadyState(abortCtx); err != nil {
			err = fmt.Errorf("failed to wait for AssetScan being ready to scan: %w", err)
			if e := cli.MarkDone(ctx, []error{err}); e != nil {
				logger.Errorf("Failed to update AssetScan status to completed with errors: %v", e)
			}
			return err
		}

		if mountVolume {
			// Set timeout for mounting volumes
			mountCtx, mountCancel := context.WithTimeout(abortCtx, DefaultMountTimeout)
			defer mountCancel()

			mountPoints, err := cli.MountVolumes(mountCtx)
			if err != nil {
				err = fmt.Errorf("failed to mount attached volume: %w", err)
				if e := cli.MarkDone(ctx, []error{err}); e != nil {
					logger.Errorf("Failed to update asset scan stat to completed with errors: %v", e)
				}
				return err
			}
			setMountPointsForFamiliesInput(mountPoints, config)
		}

		err = cli.MarkInProgress(ctx)
		if err != nil {
			return fmt.Errorf("failed to inform server %v scan has started: %w", server, err)
		}

		logger.Infof("Running scanners...")
		runErrors := families.New(config).Run(abortCtx, cli)

		err = cli.MarkDone(ctx, runErrors)
		if err != nil {
			return fmt.Errorf("failed to inform the server %v the scan was completed: %w", server, err)
		}

		if len(runErrors) > 0 {
			logger.Errorf("Errors when running families: %+v", runErrors)
		}

		return nil
	},
}

// nolint: gochecknoinits
func init() {
	rootCmd.AddCommand(scanCmd)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	scanCmd.Flags().String("config", "", "config file (default is $HOME/.vmclarity.yaml)")
	scanCmd.Flags().String("output", "", "set output directory path. Stdout is used if not set.")
	scanCmd.Flags().String("server", "", "VMClarity server to export asset scans to, for example: http://localhost:9999/api")
	scanCmd.Flags().String("asset-scan-id", "", "the AssetScan ID to monitor and report results to")
	scanCmd.Flags().Bool("mount-attached-volume", false, "discover for an attached volume and mount it before the scan")

	// TODO(sambetts) we may have to change this to our own validation when
	// we add the CI/CD scenario and there isn't an existing asset-scan-id
	// in the backend to PATCH
	scanCmd.MarkFlagsRequiredTogether("server", "asset-scan-id")
}

// loadConfig reads in config file and ENV variables if set.
func loadConfig(cfgFile string) *families.Config {
	logger.Infof("Initializing configuration...")
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
	config := &families.Config{}
	err = viper.Unmarshal(config)
	cobra.CheckErr(err)

	if logrus.IsLevelEnabled(logrus.InfoLevel) {
		configB, err := yaml.Marshal(config)
		cobra.CheckErr(err)
		logger.Infof("Using config file (%s):\n%s", viper.ConfigFileUsed(), string(configB))
	}

	return config
}

func newCli(config *families.Config, server, assetScanID, output string) (*cli.CLI, error) {
	var manager state.Manager
	var presenters []presenter.Presenter
	var err error

	if config == nil {
		return nil, errors.New("families config must not be nil")
	}

	if server != "" {
		var client *backendclient.BackendClient
		var p presenter.Presenter

		client, err = backendclient.Create(server)
		if err != nil {
			return nil, fmt.Errorf("failed to create VMClarity API client: %w", err)
		}

		manager, err = state.NewVMClarityState(client, assetScanID)
		if err != nil {
			return nil, fmt.Errorf("failed to create VMClarity state manager: %w", err)
		}

		p, err = presenter.NewVMClarityPresenter(client, assetScanID)
		if err != nil {
			return nil, fmt.Errorf("failed to create VMClarity presenter: %w", err)
		}
		presenters = append(presenters, p)
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
		if familiesConfig.SBOM.Enabled {
			familiesConfig.SBOM.Inputs = append(familiesConfig.SBOM.Inputs, sbom.Input{
				Input:     mountDir,
				InputType: string(kubeclarityutils.ROOTFS),
			})
		}

		if familiesConfig.Vulnerabilities.Enabled {
			if familiesConfig.SBOM.Enabled {
				familiesConfig.Vulnerabilities.InputFromSbom = true
			} else {
				familiesConfig.Vulnerabilities.Inputs = append(familiesConfig.Vulnerabilities.Inputs, vulnerabilities.Input{
					Input:     mountDir,
					InputType: string(kubeclarityutils.ROOTFS),
				})
			}
		}

		if familiesConfig.Secrets.Enabled {
			familiesConfig.Secrets.Inputs = append(familiesConfig.Secrets.Inputs, secrets.Input{
				StripPathFromResult: utils.PointerTo(true),
				Input:               mountDir,
				InputType:           string(kubeclarityutils.ROOTFS),
			})
		}

		if familiesConfig.Malware.Enabled {
			familiesConfig.Malware.Inputs = append(familiesConfig.Malware.Inputs, malware.Input{
				StripPathFromResult: utils.PointerTo(true),
				Input:               mountDir,
				InputType:           string(kubeclarityutils.ROOTFS),
			})
		}

		if familiesConfig.Rootkits.Enabled {
			familiesConfig.Rootkits.Inputs = append(familiesConfig.Rootkits.Inputs, rootkits.Input{
				StripPathFromResult: utils.PointerTo(true),
				Input:               mountDir,
				InputType:           string(kubeclarityutils.ROOTFS),
			})
		}

		if familiesConfig.Misconfiguration.Enabled {
			familiesConfig.Misconfiguration.Inputs = append(
				familiesConfig.Misconfiguration.Inputs,
				misconfigurationTypes.Input{
					StripPathFromResult: utils.PointerTo(true),
					Input:               mountDir,
					InputType:           string(kubeclarityutils.ROOTFS),
				},
			)
		}
	}
	return familiesConfig
}
