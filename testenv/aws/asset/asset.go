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

package asset

// Create a new asset to be scanned in the test environment.
// This asset will be an EC2 instance tagged with {"scanconfig": "test"}.
// The asset will be created outside the VMClarity Server region.
type Asset struct {
	InstanceID  string
	InstanceTag string
}

func NewAsset(id, instanceID, instanceIP, instanceTag string) *Asset {
	return &Asset{
		InstanceID:  instanceID,
		InstanceTag: instanceTag,
	}
}

func (a *Asset) Create() error {
	// Create an EC2 instance
	return nil
}

func (a *Asset) Delete() error {
	// Delete the EC2 instance
	return nil
}
