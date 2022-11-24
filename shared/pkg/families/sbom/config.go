package sbom

import "github.com/openclarity/kubeclarity/shared/pkg/config"

func (*Config) IsConfig() {}

type Config struct {
	Enabled         bool          `yaml:"enabled" mapstructure:"enabled"`
	AnalyzersList   []string      `yaml:"analyzers_list" mapstructure:"analyzers_list"`
	Inputs          []Inputs      `yaml:"inputs" mapstructure:"inputs"`
	MergeWith       []MergeWith   `yaml:"merge_with" mapstructure:"merge_with"`
	AnalyzersConfig config.Config `yaml:"analyzers_config" mapstructure:"analyzers_config"`
}

type Inputs struct {
	Input     string `yaml:"input" mapstructure:"input"`
	InputType string `yaml:"input_type" mapstructure:"input_type"`
}

type MergeWith struct {
	SbomPath string `yaml:"sbom_path" mapstructure:"sbom_path"`
}
