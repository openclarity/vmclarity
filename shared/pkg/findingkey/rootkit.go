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

package findingkey

import (
	"fmt"

	"github.com/openclarity/vmclarity/api/models"
)

type RootkitKey struct {
	Name        string
	RootkitType string
	Path        string
}

func (k RootkitKey) String() string {
	return fmt.Sprintf("%s.%s.%s", k.Name, k.RootkitType, k.Path)
}

func GenerateRootkitKey(info models.RootkitFindingInfo) RootkitKey {
	return RootkitKey{
		Name:        *info.RootkitName,
		RootkitType: string(*info.RootkitType),
		Path:        *info.Path,
	}
}
