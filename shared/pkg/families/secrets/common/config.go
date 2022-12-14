package common

import (
	gitleaksconfig "github.com/openclarity/vmclarity/shared/pkg/families/secrets/gitleaks/config"
)

type ScannersConfig struct {
	Gitleaks gitleaksconfig.Config `yaml:"gitleaks" mapstructure:"gitleaks"`
}

func (ScannersConfig) IsConfig() {}
