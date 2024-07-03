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

package job_manager // nolint:revive,stylecheck

import (
	"fmt"
	"github.com/openclarity/vmclarity/scanner/families/types"

	"github.com/sirupsen/logrus"
)

type CreateScannerFunc[CT any, RT types.Result[RT]] func(string, CT, *logrus.Entry) types.Scanner[RT]

type Factory[CT any, RT types.Result[RT]] struct {
	scanners map[string]CreateScannerFunc[CT, RT]
}

func NewFactory[CT any, RT types.Result[RT]]() *Factory[CT, RT] {
	return &Factory[CT, RT]{
		scanners: make(map[string]CreateScannerFunc[CT, RT]),
	}
}

func (f *Factory[CT, RT]) Register(name string, createJobFunc CreateScannerFunc[CT, RT]) {
	if f.scanners == nil {
		f.scanners = make(map[string]CreateScannerFunc[CT, RT])
	}

	if _, ok := f.scanners[name]; ok {
		logrus.Fatalf("%q already registered", name)
	}

	f.scanners[name] = createJobFunc
}

func (f *Factory[CT, RT]) CreateJob(name string, config CT, logger *logrus.Entry) (types.Scanner[RT], error) {
	createFunc, ok := f.scanners[name]
	if !ok {
		return nil, fmt.Errorf("%v not a registered job", name)
	}

	return createFunc(name, config, logger), nil
}
