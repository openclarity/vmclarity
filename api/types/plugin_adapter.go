package types

import plugintypes "github.com/openclarity/vmclarity/plugins/sdk/types"

// DefaultPluginAdapter is used to convert latest version Plugin API models to VMClarity
var DefaultPluginAdapter PluginAdapter = &pluginAdapter{}

// PluginAdapter is responsible for converting Plugin security findings to
// low-level VMClarity findings.
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

type pluginAdapter struct{}

func (p pluginAdapter) Result(data plugintypes.Result) ([]Finding_FindingInfo, error) {
	var findings []Finding_FindingInfo

	// Convert misconfigurations
	if misconfigurations := data.Vmclarity.Misconfigurations; misconfigurations != nil {
		for _, misconfiguration := range *misconfigurations {
			misconfiguration, err := p.Misconfiguration(misconfiguration)
			if err != nil {
				return nil, err
			}
			if misconfiguration == nil {
				continue
			}

			var finding Finding_FindingInfo
			_ = finding.FromMisconfigurationFindingInfo(*misconfiguration)
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
	severityMapping := map[plugintypes.MisconfigurationSeverity]MisconfigurationSeverity{
		plugintypes.MisconfigurationSeverityHigh:   MisconfigurationHighSeverity,
		plugintypes.MisconfigurationSeverityMedium: MisconfigurationMediumSeverity,
		plugintypes.MisconfigurationSeverityLow:    MisconfigurationLowSeverity,
		plugintypes.MisconfigurationSeverityInfo:   MisconfigurationInfoSeverity,
	}

	severity := MisconfigurationInfoSeverity
	if data.Severity != nil {
		if s, ok := severityMapping[*data.Severity]; ok {
			severity = s
		}
	}

	return &MisconfigurationFindingInfo{
		Category:    data.Category,
		Description: data.Description,
		Id:          data.Id,
		Location:    data.Location,
		Message:     data.Message,
		Remediation: data.Remediation,
		// TODO(ramizpolic): Remove ScannerName property from Misconfiguration API.
		// TODO(ramizpolic): This data is available on higher Finding object.
		ScannerName: nil,
		Severity:    &severity,
	}, nil
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
