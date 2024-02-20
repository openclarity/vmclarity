package discoverer

import (
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

const (
	DefaultEnvPrefix = "VMCLARITY_GCP"

	projectID = "project_id"
)

type Config struct {
	ProjectID string `mapstructure:"project_id"`
}

func NewConfig() (*Config, error) {
	// Avoid modifying the global instance
	v := viper.New()

	v.SetEnvPrefix(DefaultEnvPrefix)
	v.AllowEmptyEnv(true)
	v.AutomaticEnv()

	_ = v.BindEnv(projectID)

	config := &Config{}
	if err := v.Unmarshal(&config, viper.DecodeHook(mapstructure.TextUnmarshallerHookFunc())); err != nil {
		return nil, fmt.Errorf("failed to parse provider configuration. Provider=GCP: %w", err)
	}
	return config, nil
}

// nolint:cyclop
func (c Config) Validate() error {
	if c.ProjectID == "" {
		return fmt.Errorf("parameter ProjectID must be provided by setting %v_%v environment variable", DefaultEnvPrefix, strings.ToUpper(projectID))
	}

	return nil
}
