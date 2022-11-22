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
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/openclarity/vmclarity/shared/pkg/families"
)

var (
	cfgFile string
	config  *families.Config
	logger  *logrus.Entry
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "vmclarity",
	Short: "VMClarity",
	Long:  `VMClarity`,
	//Version: pkg.GitRevision,
	RunE: func(cmd *cobra.Command, args []string) error {
		logrus.Infof("Running...")
		_, err := families.New(config).Run()
		if err != nil {
			return err
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
		initConfig,
		//initAppConfig,
		initLogger,
	)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.vmclarity.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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

		// Search config in home directory with name ".families" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".families")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	cobra.CheckErr(err)

	config = &families.Config{}
	err = viper.Unmarshal(config)
	cobra.CheckErr(err)

	logrus.Infof("Using config file (%s): %+v", viper.ConfigFileUsed(), config)
}

//
//func initAppConfig() {
//	config = config.LoadConfig()
//}

func initLogger() {
	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)
	//if level, err := logrus.ParseLevel(config.LogLevel); err != nil {
	//	log.SetLevel(level)
	//}
	//if config.EnableJSONLog {
	//	log.SetFormatter(&logrus.JSONFormatter{})
	//}
	logger = log.WithField("app", "vmclarity")
}
