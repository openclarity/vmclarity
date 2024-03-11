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

package scanner

import (
	dockle_types "github.com/Portshift/dockle/pkg/types"

	"github.com/openclarity/vmclarity/scanner/types"
)

var CISDockerImpactCategory = "best-practice"

func parseDockleReport(input types.AssetScanInput, assessmentMap dockle_types.AssessmentMap) []types.Finding {
	ret := make([]types.Finding, 0, len(assessmentMap))

	for _, codeInfo := range assessmentMap {
		severity := convertDockleSeverity(codeInfo.Level)
		if severity == "" {
			// skip when no severity
			continue
		}

		description := ""
		for _, assessment := range codeInfo.Assessments {
			description += assessment.Desc + "\n"
		}
		message := dockle_types.TitleMap[codeInfo.Code]

		findingInfo := new(types.Finding_FindingInfo)
		_ = findingInfo.FromMisconfiguration(types.Misconfiguration{
			Category:    &CISDockerImpactCategory,
			Description: &description,
			Id:          &codeInfo.Code,
			Location:    &input.Path,
			Message:     &message,
			Severity:    &severity,
		})

		ret = append(ret, types.Finding{
			FindingInfo: findingInfo,
			Input:       &input,
		})
	}

	return ret
}

func convertDockleSeverity(level int) types.MisconfigurationSeverity {
	switch level {
	case dockle_types.FatalLevel:
		return types.MisconfigurationHighSeverity
	case dockle_types.WarnLevel:
		return types.MisconfigurationMediumSeverity
	case dockle_types.InfoLevel:
		return types.MisconfigurationLowSeverity
	default: // ignore PassLevel, IgnoreLevel, SkipLevel
		return ""
	}
}
