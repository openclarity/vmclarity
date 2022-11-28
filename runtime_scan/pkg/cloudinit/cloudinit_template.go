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

package cloudinit

const cloudInitTmpl string = `#cloud-config
package_upgrade: true
packages:
  - jq
  - curl
write_files:
  - path: /root/scanner_family_download.sh
    permissions: "0755"
    content: |
      #!/bin/bash

      set -euo pipefail

      curl -L -o /tmp/scanner_family.tar.gz https://example.com/scanner_family.tar.gz
      tar -xf /tmp/scanner_family.tar.gz -C /root
  - path: /root/scanconfig.json
    permissions: "0644"
    content: |
      {{ .Config }}
  - path: /etc/systemd/system/vmclarity-scanner.service
    permissions: "0644"
    content: |
      [Unit]
      Description=VMClarity scanner job
      After=network.target

      [Service]
      Type=oneshot
      WorkingDirectory=/root
      ExecStart=/root/scanner_family_cli --config=/root/scanconfig.json

      [Install]
      WantedBy=multi-user.target
runcmd:
  - [ /root/scanner_family_download.sh ]
  - [ systemctl, daemon-reload ]
  - [ systemctl, start, vmclarity-scanner.service ]
`
