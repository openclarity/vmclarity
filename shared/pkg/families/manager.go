package families

import (
	log "github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/shared/pkg/families/exploits"
	_interface "github.com/openclarity/vmclarity/shared/pkg/families/interface"
	"github.com/openclarity/vmclarity/shared/pkg/families/malware"
	"github.com/openclarity/vmclarity/shared/pkg/families/misconfiguration"
	"github.com/openclarity/vmclarity/shared/pkg/families/rootkits"
	"github.com/openclarity/vmclarity/shared/pkg/families/sbom"
	"github.com/openclarity/vmclarity/shared/pkg/families/secrets"
	"github.com/openclarity/vmclarity/shared/pkg/families/types"
	"github.com/openclarity/vmclarity/shared/pkg/families/vulnerabilities"
)

type Manager struct {
	config    *Config
	analyzers map[types.FamilyType]_interface.Family
	scanners  map[types.FamilyType]_interface.Family
	enrichers map[types.FamilyType]_interface.Family
}

func New(logger *log.Entry, config *Config) *Manager {
	manager := &Manager{
		config:    config,
		analyzers: make(map[types.FamilyType]_interface.Family),
		scanners:  make(map[types.FamilyType]_interface.Family),
		enrichers: make(map[types.FamilyType]_interface.Family),
	}

	// Analyzers
	if config.SBOM.Enabled {
		manager.analyzers[types.SBOM] = sbom.New(logger, config.SBOM)
	}

	// Scanners
	if config.Vulnerabilities.Enabled {
		manager.scanners[types.Vulnerabilities] = vulnerabilities.New(logger, config.Vulnerabilities)
	}
	if config.Secrets.Enabled {
		manager.scanners[types.Secrets] = secrets.New(logger, config.Secrets)
	}
	if config.Rootkits.Enabled {
		manager.scanners[types.Rootkits] = rootkits.New(logger, config.Rootkits)
	}
	if config.Malware.Enabled {
		manager.scanners[types.Malware] = malware.New(logger, config.Malware)
	}
	if config.Misconfiguration.Enabled {
		manager.scanners[types.Misconfiguration] = misconfiguration.New(logger, config.Misconfiguration)
	}

	// Enrichers
	if config.Exploits.Enabled {
		manager.enrichers[types.Exploits] = exploits.New(logger, config.Exploits)
	}

	return manager
}

func (m *Manager) Run() (*Results, error) {
	results := &Results{}
	if len(m.analyzers) > 0 {
		for typ, analyzer := range m.analyzers {
			ret, err := analyzer.Run(results)
			if err != nil {
				return nil, err
			}
			results.SetResults(typ, ret)
		}
	}

	if len(m.scanners) > 0 {
		for typ, scanner := range m.scanners {
			ret, err := scanner.Run(results)
			if err != nil {
				return nil, err
			}
			results.SetResults(typ, ret)
		}
	}

	if len(m.enrichers) > 0 {
		for typ, enricher := range m.enrichers {
			ret, err := enricher.Run(results)
			if err != nil {
				return nil, err
			}
			results.SetResults(typ, ret)
		}
	}

	return results, nil
}
