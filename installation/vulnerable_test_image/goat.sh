#!/bin/sh

set -xeo pipefail

for d in families/*; do
 sh ${d}/install.sh $1
done