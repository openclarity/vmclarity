package config

import (
	familiestypes "github.com/openclarity/vmclarity/scanner/families/types"
	"time"
)

type Config struct {
	Timeout  time.Duration           `yaml:"timeout" mapstructure:"timeout"`
	Registry *familiestypes.Registry `yaml:"registry" mapstructure:"registry"`
}
