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
	"context"
	"github.com/openclarity/vmclarity/scanner/types"
)

// Scanner defines the actual scanner implementation. This should be implemented
// on the example side.
type Scanner interface {
	// GetInfo returns scanner metadata details
	GetInfo(ctx context.Context) types.ScannerInfo

	// Scan performs a scan for a given input. ScanFinding.FindingInfo,
	// ScanFinding.ScanID and ScanFinding.Input are required to be populated. Note
	// that ScanFinding.Id should not be populated as this will be done by the server
	// itself. All other fields are optional.
	Scan(ctx context.Context, scanID string, input types.ScanInput) ([]types.ScanFinding, error)
}
