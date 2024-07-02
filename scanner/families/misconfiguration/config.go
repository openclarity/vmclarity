package misconfiguration

import (
	"github.com/openclarity/vmclarity/scanner/families/misconfiguration/types"
	familiestypes "github.com/openclarity/vmclarity/scanner/families/types"
)

type Config struct {
	Enabled         bool                  `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	ScannersList    []string              `yaml:"scanners_list" mapstructure:"scanners_list"`
	StripInputPaths bool                  `yaml:"strip_input_paths" mapstructure:"strip_input_paths"`
	Inputs          []familiestypes.Input `yaml:"inputs" mapstructure:"inputs"`
	ScannersConfig  types.ScannersConfig  `yaml:"scanners_config" mapstructure:"scanners_config"`
}
