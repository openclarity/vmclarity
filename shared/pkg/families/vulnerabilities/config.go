package vulnerabilities

import "github.com/openclarity/vmclarity/shared/pkg/families/types"

func (*Config) IsConfig() {}

type Config struct {
	Enabled         bool              `yaml:"enabled" mapstructure:"enabled"`
	ScannersList    []string          `yaml:"scanners_list" mapstructure:"scanners_list"`
	Inputs          []Inputs          `yaml:"inputs" mapstructure:"inputs"`
	InputFromFamily []InputFromFamily `yaml:"input_from_family" mapstructure:"input_from_family"`
	GrypeConfig     GrypeConfig       `yaml:"grype_config" mapstructure:"grype_config"`
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
