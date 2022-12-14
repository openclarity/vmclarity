// Copyright © 2022 Cisco Systems, Inc. and its affiliates.
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
	"encoding/json"
	"fmt"
	"os"

	"github.com/ghodss/yaml"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/openclarity/vmclarity/cli/pkg"
	"github.com/openclarity/vmclarity/shared/pkg/families"
	"github.com/openclarity/vmclarity/shared/pkg/families/results"
	"github.com/openclarity/vmclarity/shared/pkg/families/sbom"
	"github.com/openclarity/vmclarity/shared/pkg/families/secrets"
	"github.com/openclarity/vmclarity/shared/pkg/families/vulnerabilities"
)

var (
	cfgFile string
	config  *families.Config
	logger  *logrus.Entry
	output  string

	server       string
	scanResultId string
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:     "vmclarity",
	Short:   "VMClarity",
	Long:    `VMClarity`,
	Version: pkg.GitRevision,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Infof("Running...")
		if server != "" {
			err := MarkScanResultInProgress()
			if err != nil {
				return fmt.Errorf("failed to inform server %v scan has started: %w", server, err)
			}
		}

		res, famerr := families.New(logger, config).Run()

		if server != "" {
			errors := []error{famerr}

			exportErrors := ExportResults(res)
			errors = append(errors, exportErrors...)

			err := MarkScanResultDone(errors)
			if err != nil {
				return fmt.Errorf("failed to inform the server %v the scan was completed: %w", server, err)
			}
		}

		if famerr != nil {
			return fmt.Errorf("failed to run families: %v", famerr)
		}

		if config.SBOM.Enabled {
			sbomResults, err := results.GetResult[*sbom.Results](res)
			if err != nil {
				return fmt.Errorf("failed to get sbom results: %v", err)
			}

			outputFormat := config.SBOM.AnalyzersConfig.Analyzer.OutputFormat
			sbomBytes, err := sbomResults.EncodeToBytes(outputFormat)
			if err != nil {
				return fmt.Errorf("failed to encode sbom results to bytes: %w", err)
			}

			// TODO: Need to implement a better presenter
			err = Output(sbomBytes, "sbom")
			if err != nil {
				return fmt.Errorf("failed to output sbom results: %v", err)
			}
		}

		if config.Vulnerabilities.Enabled {
			vulnerabilitiesResults, err := results.GetResult[*vulnerabilities.Results](res)
			if err != nil {
				return fmt.Errorf("failed to get sbom results: %v", err)
			}

			bytes, _ := json.Marshal(vulnerabilitiesResults.MergedResults)
			err = Output(bytes, "vulnerabilities")
			if err != nil {
				return fmt.Errorf("failed to output vulnerabilities results: %v", err)
			}
		}

		if config.Secrets.Enabled {
			secretsResults, err := results.GetResult[*secrets.Results](res)
			if err != nil {
				return fmt.Errorf("failed to get secrets results: %v", err)
			}

			bytes, _ := json.Marshal(secretsResults)
			err = Output(bytes, "secrets")
			if err != nil {
				return fmt.Errorf("failed to output secrets results: %v", err)
			}
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
	rootCmd.PersistentFlags().StringVar(&output, "output", "", "set file path output (default: stdout)")
	rootCmd.PersistentFlags().StringVar(&server, "server", "", "VMClarity server to export scan results to, for example: http://localhost:9999/api")
	rootCmd.PersistentFlags().StringVar(&scanResultId, "scan-result-id", "", "the ScanResult ID to export the scan results to")

	// TODO(sambetts) we may have to change this to our own validation when
	// we add the CI/CD scenario and there isn't an existing scan-result-id
	// in the backend to PATCH
	rootCmd.MarkFlagsRequiredTogether("server", "scan-result-id")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	logrus.Infof("init config")
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

	if logrus.IsLevelEnabled(logrus.InfoLevel) {
		configB, err := yaml.Marshal(config)
		cobra.CheckErr(err)
		logrus.Infof("Using config file (%s):\n%s", viper.ConfigFileUsed(), string(configB))
	}
}

func initLogger() {
	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)
	logger = log.WithField("app", "vmclarity")
}

func Output(bytes []byte, outputPrefix string) error {
	if output == "" {
		os.Stdout.Write([]byte(fmt.Sprintf("%s results:\n", outputPrefix)))
		os.Stdout.Write(bytes)
		os.Stdout.Write([]byte("\n=================================================\n"))
		return nil
	}

	filePath := outputPrefix + "." + output
	logger.Infof("Writing results to %v...", filePath)

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666) // nolint:gomnd,gofumpt
	if err != nil {
		return fmt.Errorf("failed open file %s: %v", filePath, err)
	}
	defer file.Close()

	_, err = file.Write(bytes)
	if err != nil {
		return fmt.Errorf("failed to write bytes to file %s: %v", filePath, err)
	}

	return nil
}
