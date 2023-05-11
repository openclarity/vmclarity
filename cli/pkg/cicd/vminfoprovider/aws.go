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

	azResp, err := http.Get("http://169.254.169.254/latest/meta-data/placement/availability-zone") // nolint:noctx
	if err != nil {
		return "", "", fmt.Errorf("failed to get availability-zone: %v", err)
	}
	defer azResp.Body.Close()
	availabilityZone, err := io.ReadAll(azResp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read availability-zone response body: %v", err)
	}
	az := string(availabilityZone)
	return string(instanceID), az[:len(az)-1], nil
}
