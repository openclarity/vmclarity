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

package vminfoprovider

import (
	"fmt"
	"io"
	"net/http"
)

type AWSInfoProvider struct{}

func CreateNewAWSInfoProvider() *AWSInfoProvider {
	return &AWSInfoProvider{}
}

func (a *AWSInfoProvider) GetVMInfo() (string, string, error) {
	instanceIDResp, err := http.Get("http://169.254.169.254/latest/meta-data/instance-id") // nolint:noctx
	if err != nil {
		return "", "", fmt.Errorf("failed to get instance-id: %v", err)
	}
	defer instanceIDResp.Body.Close()
	instanceID, err := io.ReadAll(instanceIDResp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read instance-id response body: %v", err)
	}

	regionResp, err := http.Get("http://169.254.169.254/latest/meta-data/placement/region") // nolint:noctx
	if err != nil {
		return "", "", fmt.Errorf("failed to get region: %v", err)
	}
	defer regionResp.Body.Close()
	region, err := io.ReadAll(regionResp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read region response body: %v", err)
	}

	// cut the last character from the availability-zone in order to return only the region
	return string(instanceID), string(region), nil
}
