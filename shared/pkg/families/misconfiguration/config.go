package misconfiguration

func (*Config) IsConfig() {}

type Config struct {
	Enabled bool `json:"enabled" yaml:"enabled"`
}
