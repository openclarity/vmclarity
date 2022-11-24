package vulnerabilities

import (
	"github.com/openclarity/kubeclarity/shared/pkg/config"

	"github.com/openclarity/vmclarity/shared/pkg/families/types"
)

func (*Config) IsConfig() {}

type Config struct {
	Enabled        bool          `yaml:"enabled" mapstructure:"enabled"`
	ScannersList   []string      `yaml:"scanners_list" mapstructure:"scanners_list"`
	Inputs         []Inputs      `yaml:"inputs" mapstructure:"inputs"`
	InputFromSbom  bool          `yaml:"input_from_sbom" mapstructure:"input_from_sbom"`
	ScannersConfig config.Config `yaml:"scanners_config" mapstructure:"scanners_config"`
}

type Inputs struct {
	Input     string `yaml:"input" mapstructure:"input"`
	InputType string `yaml:"input_type" mapstructure:"input_type"`
}

type InputFromFamily struct {
	FamilyType types.FamilyType `yaml:"family_type" mapstructure:"family_type"`
}

type GrypeConfig struct {
	Scope string `yaml:"scope" mapstructure:"scope"`
}
