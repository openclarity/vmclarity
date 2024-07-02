// Copyright © 2023 Cisco Systems, Inc. and its affiliates.
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

package types

import (
	"context"
	"errors"
	"sync"
)

type FamilyType string

const (
	SBOM             FamilyType = "sbom"
	Vulnerabilities  FamilyType = "vulnerabilities"
	Secrets          FamilyType = "secrets"
	Rootkits         FamilyType = "rootkits"
	Malware          FamilyType = "malware"
	Misconfiguration FamilyType = "misconfiguration"
	InfoFinder       FamilyType = "infofinder"
	Exploits         FamilyType = "exploits"
	Plugins          FamilyType = "plugins"
)

type FamilyResult interface {
	IsResult()
}

type Family interface {
	GetType() FamilyType
	Run(context.Context, *FamiliesResults) (FamilyResult, error)
}

// FamiliesResults stores results from all families. Safe for concurrent usage.
type FamiliesResults struct {
	mu      sync.RWMutex
	results []FamilyResult
}

func NewFamiliesResults() *FamiliesResults {
	return &FamiliesResults{}
}

func (r *FamiliesResults) SetFamilyResult(result FamilyResult) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.results = append(r.results, result)
}

// GetFamilyResult returns results for a specific family from the given results.
func GetFamilyResult[familyType FamilyResult](r *FamiliesResults) (familyType, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, result := range r.results {
		res, ok := result.(familyType)
		if ok {
			return res, nil
		}
	}

	var res familyType
	return res, errors.New("missing result")
}
