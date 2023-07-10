package external

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	DefaultEnvPrefix = "VMCLARITY_EXTERNAL"
)

type Config struct {
	ProviderPluginAddress string `mapstructure:"provider_plugin_address"`
}

func (c *Config) Validate() error {
	if c.ProviderPluginAddress == "" {
		return fmt.Errorf("parameter ProviderPluginAddress must be provided")
	}

	return nil
}

func NewConfig() (*Config, error) {
	// Avoid modifying the global instance
	v := viper.New()

	v.SetEnvPrefix(DefaultEnvPrefix)
	v.AllowEmptyEnv(true)
	v.AutomaticEnv()

	_ = v.BindEnv("provider_plugin_address")

	config := &Config{}
	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to parse provider configuration. Provider=External: %w", err)
	}

	return config, nil
}
