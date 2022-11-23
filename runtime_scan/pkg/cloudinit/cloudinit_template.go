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

// TODO example scanner script, needs to be updated
const cloudInitTmpl string = `#cloud-config
package_upgrade: true
packages:
  - jq
write_files:
  - path: /root/scanscript.sh
    permissions: "0755"
    content: |
      #!/bin/bash

      set -euo pipefail

      configjson=scanconfig.json

      if [[ ! -f "${configjson}" ]]; then
          echo "${configjson} not exists."
          exit 1
      fi

      SCANCONFIG=$(cat ${configjson})

      DIR=$(echo ${SCANCONFIG} | jq .directory_to_scan)
      SERVER=$(echo ${SCANCONFIG} | jq .server_to_report)

      validate_config() {
          if [[ ${SERVER} == null ]]; then
              echo "server not set"
              exit 1
          fi

          if [[ ${DIR} == null ]]; then
              echo "dir not set"
              exit 1
          fi
          if [[ ! -d "${DIR}" ]]; then
              echo "${DIR} not exists."
              exit 1
          fi
      }

      vulscan=$(echo ${SCANCONFIG} | jq .vulnerability_scan.vuls)
      rkscan=$(echo ${SCANCONFIG} | jq .rootkit_scan.chkrootkit)
      misconfigscan=$(echo ${SCANCONFIG} | jq .misconfig_scan.lynis)
      secretscan=$(echo ${SCANCONFIG} | jq .secret_scan.trufflehog)
      malwarescan=$(echo ${SCANCONFIG} | jq .malewre_scan.clamav)
      expcheck=$(echo ${SCANCONFIG} | jq .exploit_check.vuls)

      install_vuls() {
          echo "installing vuls..."
      }

      install_chkrootkit() {
          echo "install chkrootkit..."
      }

      install_lynis() {
          echo "install lynis..."
      }

      install_trufflehog() {
          echo "install trufflehog..."
      }

      install_clamav() {
          echo "install clamav..."
      }

      install_scanners() {
          if [[ ${vulscan} != null ]]; then
              echo "Vulnerability scan with vuls is enabled..."
              install_vuls
          fi
              if [[ ${rkscan} != null ]]; then
              echo "Vulnerability scan with vuls is enabled..."
              install_chkrootkit
          fi
          if [[ ${misconfigscan} != null ]]; then
              echo "Vulnerability scan with vuls is enabled..."
              install_lynis
          fi
          if [[ ${secretscan} != null ]]; then
              echo "Vulnerability scan with vuls is enabled..."
              install_trufflehog
          fi
          if [[ ${malwarescan} != null ]]; then
              echo "Vulnerability scan with vuls is enabled..."
              install_clamav
          fi
      }

      validate_config
      install_scanners

      run_vuls() {
          echo "run vuls..."
          config=$(echo ${vulscan} | jq .config)
          echo "runnig with config: ${config}"
      }

      run_chkrootkit() {
          echo "run chkrootkit..."
          config=$(echo ${rkscan} | jq .config)
          echo "runnig with config: ${config}"
      }

      run_lynis() {
          echo "run lynis..."
          config=$(echo ${misconfigscan} | jq .config)
          echo "runnig with config: ${config}"
      }

      run_trufflehog() {
          echo "run trufflehog..."
          config=$(echo ${secretscan} | jq .config)
          echo "runnig with config: ${config}"
      }

      run_clamav() {
          echo "run clamav..."
          config=$(echo ${malwarescan} | jq .config)
          echo "runnig with config: ${config}"
      }

      run_exploit_check() {
          echo "run vuls to check exploit..."
          config=$(echo ${expcheck} | jq .config)
          echo "runnig with config: ${config}"
      }

      run_scanners() {
          if [[ ${vulscan} != null ]]; then
              run_vuls
          fi
          if [[ ${rkscan} != null ]]; then
              run_chkrootkit
          fi
          if [[ ${misconfigscan} != null ]]; then
              run_lynis
          fi
          if [[ ${secretscan} != null ]]; then
              run_trufflehog
          fi
          if [[ ${malwarescan} != null ]]; then
              run_clamav
          fi
          if [[ ${expcheck} != null ]]; then
              run_exploit_check
          fi
      }

      report_results() {
          echo "report results to ${SERVER}..."
      }

      run_scanners
      report_results
      touch /root/finished
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
      ExecStart=/root/scanscript.sh

      [Install]
      WantedBy=multi-user.target
runcmd:
  - [ systemctl, daemon-reaload ]
  - [ systemctl, start, vmclarity-scanner.service ]
`
