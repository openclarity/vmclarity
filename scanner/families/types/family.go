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
	"github.com/openclarity/vmclarity/scanner/types"
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

// Family defines interface required to fully run a family.
type Family[T any] interface {
	GetType() FamilyType
	Run(context.Context, *Results) (T, error)
}

// Scanner defines implementation of a family scanner. It should be
// concurrently-safe as Scan can be called concurrently.
type Scanner[T any] interface {
	Scan(ctx context.Context, sourceType types.InputType, source string) (T, error)
}
