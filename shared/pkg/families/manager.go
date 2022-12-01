// Copyright © 2022 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package families

import (
	log "github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/shared/pkg/families/exploits"
	_interface "github.com/openclarity/vmclarity/shared/pkg/families/interface"
	"github.com/openclarity/vmclarity/shared/pkg/families/malware"
	"github.com/openclarity/vmclarity/shared/pkg/families/misconfiguration"
	"github.com/openclarity/vmclarity/shared/pkg/families/results"
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

func (m *Manager) Run() (*results.Results, error) {
	familiesResults := results.New()
	if len(m.analyzers) > 0 {
		for _, analyzer := range m.analyzers {
			ret, err := analyzer.Run(familiesResults)
			if err != nil {
				return nil, err
			}
			familiesResults.SetResults(ret)
		}
	}

	if len(m.scanners) > 0 {
		for _, scanner := range m.scanners {
			ret, err := scanner.Run(familiesResults)
			if err != nil {
				return nil, err
			}
			familiesResults.SetResults(ret)
		}
	}

	if len(m.enrichers) > 0 {
		for _, enricher := range m.enrichers {
			ret, err := enricher.Run(familiesResults)
			if err != nil {
				return nil, err
			}
			familiesResults.SetResults(ret)
		}
	}

	return familiesResults, nil
}
