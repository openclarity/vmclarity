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

package scan

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/api/types"
)

type ScanReconcileEvent struct {
	ScanID types.ScanID
}

func (e ScanReconcileEvent) ToFields() log.Fields {
	return log.Fields{
		"ScanID": e.ScanID,
	}
}

func (e ScanReconcileEvent) String() string {
	return fmt.Sprintf("ScanID=%s", e.ScanID)
}

func (e ScanReconcileEvent) Hash() string {
	return e.ScanID
}