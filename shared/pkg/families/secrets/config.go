package secrets

func (*Config) IsConfig() {}

type Config struct {
	Enabled      bool     `yaml:"enabled" mapstructure:"enabled"`
	ScannersList []string `yaml:"scanners_list" mapstructure:"scanners_list"`
	Inputs       []Inputs `yaml:"inputs" mapstructure:"inputs"`
}

type Inputs struct {
	Input     string `yaml:"input" mapstructure:"input"`
	InputType string `yaml:"input_type" mapstructure:"input_type"`
}
