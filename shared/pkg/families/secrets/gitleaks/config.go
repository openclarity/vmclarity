package gitleaks

type Config struct {
	BinaryPath string `yaml:"binary_path" mapstructure:"binary_path"`
}
