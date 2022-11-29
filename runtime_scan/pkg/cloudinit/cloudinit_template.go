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
  - curl
  - apt-transport-https
  - ca-certificates
  - software-properties-common
write_files:
  - path: /root/install_docker_ce.sh
    permissions: "0755"
    content: |
      #!/bin/bash

      set -euo pipefail

      curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
      echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null
      apt update
      apt install docker-ce -y
  - path: /root/config/scanconfig.json
    permissions: "0644"
    content: |
      {{ .Config }}
  - path: /etc/systemd/system/vmclarity-scanner.service
    permissions: "0644"
    content: |
      [Unit]
      Description=VMClarity scanner job
      Requires=docker.service
      After=network.target docker.service

      [Service]
      Type=oneshot
      WorkingDirectory=/root
      ExecStart=docker run --rm --name %n -v /root/config:/config busybox ls /config

      [Install]
      WantedBy=multi-user.target
runcmd:
  - [ /root/install_docker_ce.sh ]
  - [ systemctl, daemon-reload ]
  - [ systemctl, start, docker.service ]
  - [ docker, pull, busybox ]
  - [ systemctl, start, vmclarity-scanner.service ]
`
