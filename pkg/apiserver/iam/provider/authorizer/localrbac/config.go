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

package localrbac

import (
	"fmt"

	"github.com/spf13/viper"
)

const authzLocalRbacRuleFilePathEnvVar = "AUTHZ_LOCAL_RBAC_RULE_FILEPATH"

type Config struct {
	RuleFilePath string `json:"rule-filepath"`
}

func LoadConfig() (*Config, error) {
	return &Config{
		RuleFilePath: viper.GetString(authzLocalRbacRuleFilePathEnvVar),
	}, nil
}

func (c *Config) Validate() error {
	if c.RuleFilePath == "" {
		return fmt.Errorf("must specify issuer for Authorizer=localrbac")
	}

	return nil
}
