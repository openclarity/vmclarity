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

package aws

import (
	"errors"
	"time"

	"github.com/aws/smithy-go"

	"github.com/openclarity/vmclarity/runtime_scan/pkg/provider"
)

type (
	FatalError     = provider.FatalError
	RetryableError = provider.RetryableError
)

const (
	DefaultRetryAfter     = 15 * time.Second
	RetryServerErrorAfter = time.Minute

	ErrorCodeUnauthorizedOperation = "UnauthorizedOperation"
)

func WrapError(err error) error {
	var fatalError FatalError
	if errors.As(err, &fatalError) {
		return err
	}

	var retryableError RetryableError
	if errors.As(err, &retryableError) {
		return err
	}

	var awsApiError smithy.APIError
	if errors.As(err, &awsApiError) {
		if awsApiError.ErrorCode() == ErrorCodeUnauthorizedOperation {
			return FatalError{err}
		}

		switch awsApiError.ErrorFault() {
		case smithy.FaultServer:
			return RetryableError{err, RetryServerErrorAfter}
		case smithy.FaultClient, smithy.FaultUnknown:
			return FatalError{err}
		}
	}

	return RetryableError{err, DefaultRetryAfter}
}
