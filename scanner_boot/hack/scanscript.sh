#!/bin/bash

set -euo pipefail

configjson=scanconfig.json

if [[ ! -f "${configjson}" ]]; then
    echo "${configjson} not exists."
    exit 1
fi

SCANCONFIG=$(cat ${configjson})

DIR=$(echo ${SCANCONFIG} | jq .directory_to_scan ${configjson})
SERVER=$(echo ${SCANCONFIG} | jq .server_to_report ${configjson})

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
        #exit 1
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
