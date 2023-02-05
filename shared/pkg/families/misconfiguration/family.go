// Copyright © 2022 Cisco Systems, Inc. and its affiliates.
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

package misconfiguration

import (
	"github.com/openclarity/vmclarity/shared/pkg/families/types"
	log "github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/shared/pkg/families/interfaces"
	"github.com/openclarity/vmclarity/shared/pkg/families/results"
)

type Misconfiguration struct {
	conf   Config
	logger *log.Entry
}

func (m Misconfiguration) Run(res *results.Results) (interfaces.IsResults, error) {
	// TODO implement me
	m.logger.Info("Misconfiguration Run...")
	m.logger.Info("Misconfiguration Done...")
	return &Results{}, nil
}

func (m Misconfiguration) GetType() types.FamilyType {
	return types.Misconfiguration
}

// ensure types implement the requisite interfaces.
var _ interfaces.Family = &Misconfiguration{}

func New(logger *log.Entry, conf Config) *Misconfiguration {
	return &Misconfiguration{
		conf:   conf,
		logger: logger.Dup().WithField("family", "misconfiguration"),
	}
}
