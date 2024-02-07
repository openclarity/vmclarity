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

package registry

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"testing"
)

// TODO(ramizpolic): add this file but strip everything from it
// TODO(ramizpolic): add user NT registry with some preinstalled apps
const drivePath = "/media/ramiz-polic/System"

func TestRegistry(t *testing.T) {
	reg, err := NewRegistry(drivePath, log.NewEntry(log.New()))
	if err != nil {
		t.Fatalf("should not error %v", err)
	}

	bom, err := reg.GetBOM()
	if err != nil {
		t.Fatalf("should not error %v", err)
	}

	prettyPrint(bom)
}

func prettyPrint(data interface{}) {
	b, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(b))
}
