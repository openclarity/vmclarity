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

package secrets

import (
	"reflect"
	"testing"

	"github.com/openclarity/vmclarity/scanner/families/secrets/types"
)

func TestStripPathFromResult(t *testing.T) {
	type args struct {
		result []types.Finding
		path   string
	}
	tests := []struct {
		name string
		args args
		want []types.Finding
	}{
		{
			name: "sanity",
			args: args{
				result: []types.Finding{
					{
						Description: "description",
						File:        "/mnt/file1",
						Message:     "message",
						Fingerprint: "finger:/mnt:hello",
					},
					{
						Description: "description2",
						File:        "/mnt/file2",
						Message:     "message",
						Fingerprint: "finger:/mnt/foo:hello2",
					},
					{
						Description: "description3",
						File:        "file3",
						Message:     "message",
						Fingerprint: "finger:hello3",
					},
				},
				path: "/mnt",
			},
			want: []types.Finding{
				{
					Description: "description",
					File:        "/file1",
					Message:     "message",
					Fingerprint: "finger:/:hello",
				},
				{
					Description: "description2",
					File:        "/file2",
					Message:     "message",
					Fingerprint: "finger:/foo:hello2",
				},
				{
					Description: "description3",
					File:        "file3",
					Message:     "message",
					Fingerprint: "finger:hello3",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stripPathFromResult(tt.args.result, tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StripPathFromResult() = %v, want %v", got, tt.want)
			}
		})
	}
}
