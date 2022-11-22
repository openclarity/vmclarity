package families

import (
	"github.com/openclarity/vmclarity/shared/pkg/families/exploits"
	"github.com/openclarity/vmclarity/shared/pkg/families/malware"
	"github.com/openclarity/vmclarity/shared/pkg/families/misconfiguration"
	"github.com/openclarity/vmclarity/shared/pkg/families/rootkits"
	"github.com/openclarity/vmclarity/shared/pkg/families/sbom"
	"github.com/openclarity/vmclarity/shared/pkg/families/secrets"
	"github.com/openclarity/vmclarity/shared/pkg/families/vulnerabilities"
)

type Config struct {
	// Analyzers
	SBOM sbom.Config `json:"sbom" yaml:"sbom" mapstructure:"sbom"`

	// Scanners
	Vulnerabilities  vulnerabilities.Config  `json:"vulnerabilities" yaml:"vulnerabilities" mapstructure:"vulnerabilities"`
	Secrets          secrets.Config          `json:"secrets" yaml:"secrets" mapstructure:"secrets"`
	Rootkits         rootkits.Config         `json:"rootkits" yaml:"rootkits" mapstructure:"rootkits"`
	Malware          malware.Config          `json:"malware" yaml:"malware" mapstructure:"malware"`
	Misconfiguration misconfiguration.Config `json:"misconfiguration" yaml:"misconfiguration" mapstructure:"misconfiguration"`

	// Enrichers
	Exploits exploits.Config `json:"exploits" yaml:"exploits" mapstructure:"exploits"`
}
