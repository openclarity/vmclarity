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

	"github.com/openclarity/vmclarity/shared/pkg/families/interface"
	"github.com/openclarity/vmclarity/shared/pkg/families/types"
)

type Results struct {
	SBOM             _interface.IsResults `json:"sbom"`
	Vulnerabilities  _interface.IsResults `json:"vulnerabilities"`
	Secrets          _interface.IsResults `json:"secrets"`
	Rootkits         _interface.IsResults `json:"rootkits"`
	Malware          _interface.IsResults `json:"malware"`
	Misconfiguration _interface.IsResults `json:"misconfiguration"`
	Exploits         _interface.IsResults `json:"exploits"`
}

func (r *Results) SetResults(tpe types.FamilyType, result _interface.IsResults) {
	switch tpe {
	case types.SBOM:
		r.SBOM = result
	case types.Vulnerabilities:
		r.Vulnerabilities = result
	case types.Secrets:
		r.Secrets = result
	case types.Rootkits:
		r.Rootkits = result
	case types.Malware:
		r.Malware = result
	case types.Misconfiguration:
		r.Misconfiguration = result
	case types.Exploits:
		r.Exploits = result
	default:
		log.Fatalf("unknown family type %v", tpe)
	}
}

func (r *Results) GetResults(tpe types.FamilyType) _interface.IsResults {
	switch tpe {
	case types.SBOM:
		return r.SBOM
	case types.Vulnerabilities:
		return r.Vulnerabilities
	case types.Secrets:
		return r.Secrets
	case types.Rootkits:
		return r.Rootkits
	case types.Malware:
		return r.Malware
	case types.Misconfiguration:
		return r.Misconfiguration
	case types.Exploits:
		return r.Exploits
	}

	log.Fatalf("unknown family type %v", tpe)
	return nil
}
