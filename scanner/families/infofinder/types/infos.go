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

package types

import "github.com/openclarity/vmclarity/scanner/common"

type InfoType string

const (
	SSHKnownHostFingerprint     InfoType = "sshKnownHostFingerprint"
	SSHAuthorizedKeyFingerprint InfoType = "sshAuthorizedKeyFingerprint"
	SSHPrivateKeyFingerprint    InfoType = "sshPrivateKeyFingerprint"
	SSHDaemonKeyFingerprint     InfoType = "sshDaemonKeyFingerprint"
)

type Info struct {
	Type InfoType `json:"type"`
	Path string   `json:"path"`
	Data string   `json:"data"`
}

type FlattenedInfo struct {
	Info
	ScannerName string `json:"ScannerName"`
}

type Infos struct {
	Metadata       common.ScanMetadata `json:"Metadata"`
	FlattenedInfos []FlattenedInfo     `json:"Infos"`
}

func NewInfos() *Infos {
	return &Infos{
		FlattenedInfos: []FlattenedInfo{},
	}
}

func (r *Infos) Merge(meta common.ScanInputMetadata, infos []Info) {
	r.Metadata.Merge(meta)

	for _, info := range infos {
		r.FlattenedInfos = append(r.FlattenedInfos, FlattenedInfo{
			ScannerName: meta.ScannerName,
			Info:        info,
		})
	}
}
