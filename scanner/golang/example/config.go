package main

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

const (
	DefaultEnvPrefix     = "VMCLARITY_SCANNER"
	DefaultListenAddress = "0.0.0.0:8765"
	DefaultLogLevel      = "info"
)

type Config struct {
	ListenAddress string `json:"listen-address,omitempty" mapstructure:"listen_address"`
	LogLevel      string `json:"log-level,omitempty" mapstructure:"log_level"`
}

func NewConfig() (*Config, error) {
	v := viper.NewWithOptions(
		viper.KeyDelimiter("."),
		viper.EnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_")),
	)

	v.SetEnvPrefix(DefaultEnvPrefix)
	v.AllowEmptyEnv(true)
	v.AutomaticEnv()

	_ = v.BindEnv("listen_address")
	v.SetDefault("listen_address", DefaultListenAddress)

	_ = v.BindEnv("log_level")
	v.SetDefault("log_level", DefaultLogLevel)

	config := &Config{}
	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to load API Server configuration: %w", err)
	}

	return config, nil
}
