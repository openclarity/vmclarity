// Copyright © 2024 Cisco Systems, Inc. and its affiliates.
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

package gorm

import (
	"testing"

	"github.com/openclarity/vmclarity/api/types"
	"github.com/openclarity/vmclarity/core/to"
)

func TestFindings(t *testing.T) {
	// create test database
	db := createTestDB(t)

	// create vulnerability
	var finfo types.FindingInfo
	_ = finfo.FromVulnerabilityFindingInfo(types.VulnerabilityFindingInfo{
		Package: &types.Package{
			Name:    to.Ptr("Name"),
			Version: to.Ptr("Version"),
		},
		Severity:          to.Ptr(types.HIGH),
		VulnerabilityName: to.Ptr("Name"),
	})

	_, err := db.FindingsTable().CreateFinding(types.Finding{
		FindingInfo: &finfo,
	})
	if err != nil {
		t.Errorf("failed to create vulnerability finding: %v", err)
	}

	// create package
	_ = finfo.FromPackageFindingInfo(types.PackageFindingInfo{
		Name:    to.Ptr("Name"),
		Version: to.Ptr("Version"),
	})
	_, err = db.FindingsTable().CreateFinding(types.Finding{
		FindingInfo: &finfo,
	})
	if err != nil {
		t.Errorf("failed to create package finding: %v", err)
	}

	// get all findings
	_, err = db.FindingsTable().GetFindings(types.GetFindingsParams{})
	if err != nil {
		t.Errorf("failed to get all findings: %v", err)
	}
}
