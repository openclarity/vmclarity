package sbom

func (*Config) IsConfig() {}

type Config struct {
	Enabled       bool        `yaml:"enabled" mapstructure:"enabled"`
	AnalyzersList []string    `yaml:"analyzers_list" mapstructure:"analyzers_list"`
	Inputs        []Inputs    `yaml:"inputs" mapstructure:"inputs"`
	MergeWith     []MergeWith `yaml:"merge_with" mapstructure:"merge_with"`
	SyftConfig    SyftConfig  `yaml:"syft_config" mapstructure:"syft_config"`
}

type Inputs struct {
	Input     string `yaml:"input" mapstructure:"input"`
	InputType string `yaml:"input_type" mapstructure:"input_type"`
}

type MergeWith struct {
	SbomPath string `yaml:"sbom_path" mapstructure:"sbom_path"`
}

type SyftConfig struct {
	Scope string `yaml:"scope" mapstructure:"scope"`
}
