#!/usr/bin/env bash

set -xeo pipefail

for d in families/*; do
 sh ${d}/install.sh
done