package types

import plugintypes "github.com/openclarity/vmclarity/plugins/sdk/types"

// DefaultPluginAdapter is used to convert Plugin to VMClarity objects
var DefaultPluginAdapter PluginAdapter = &pluginAdapter{}

type PluginAdapter interface {
	Result(data plugintypes.Result) ([]Finding_FindingInfo, error)

	Exploit(data plugintypes.Exploit) (*ExploitFindingInfo, error)
	InfoFinder(data plugintypes.InfoFinder) (*InfoFinderFindingInfo, error)
	Malware(data plugintypes.Malware) (*MalwareFindingInfo, error)
	Misconfiguration(data plugintypes.Misconfiguration) (*MisconfigurationFindingInfo, error)
	Package(data plugintypes.Package) (*PackageFindingInfo, error)
	Rootkit(data plugintypes.Rootkit) (*RootkitFindingInfo, error)
	Secret(data plugintypes.Secret) (*SecretFindingInfo, error)
	Vulnerability(data plugintypes.Vulnerability) (*VulnerabilityFindingInfo, error)
}

// implement PluginAdapter
type pluginAdapter struct{}

func (p pluginAdapter) Result(data plugintypes.Result) ([]Finding_FindingInfo, error) {
	var findings []Finding_FindingInfo

	// Convert secrets
	if secrets := data.Vmclarity.Secrets; secrets != nil {
		for _, secret := range *secrets {
			secret, err := p.Secret(secret)
			if err != nil {
				return nil, err
			}

			var finding Finding_FindingInfo
			_ = finding.FromSecretFindingInfo(*secret)
			findings = append(findings, finding)
		}
	}
	// Convert others...

	return findings, nil
}

func (p pluginAdapter) Exploit(data plugintypes.Exploit) (*ExploitFindingInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (p pluginAdapter) InfoFinder(data plugintypes.InfoFinder) (*InfoFinderFindingInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (p pluginAdapter) Malware(data plugintypes.Malware) (*MalwareFindingInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (p pluginAdapter) Misconfiguration(data plugintypes.Misconfiguration) (*MisconfigurationFindingInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (p pluginAdapter) Package(data plugintypes.Package) (*PackageFindingInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (p pluginAdapter) Rootkit(data plugintypes.Rootkit) (*RootkitFindingInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (p pluginAdapter) Secret(data plugintypes.Secret) (*SecretFindingInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (p pluginAdapter) Vulnerability(data plugintypes.Vulnerability) (*VulnerabilityFindingInfo, error) {
	//TODO implement me
	panic("implement me")
}
