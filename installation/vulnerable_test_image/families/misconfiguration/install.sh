#!/bin/sh

set -xeo pipefail

echo installing misconfiguration tests...
mkdir -p ${1}/families/misconfiguration
cp families/misconfiguration/misconfiguration.example ${1}/families/misconfiguration/misgonfiguration.example
