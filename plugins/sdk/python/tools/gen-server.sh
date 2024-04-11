#!/usr/bin/env bash

## Script vars
SCRIPTPATH="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
PLUGINPATH="$(realpath $SCRIPTPATH/../../..)"
PYTHONSDKPATH="$(realpath $SCRIPTPATH/..)"

## Generate
echo "Generating Python server code from $PLUGINPATH"
docker run --rm \
  -v ${PLUGINPATH}:/src openapitools/openapi-generator-cli generate \
  -i /src/openapi.yaml \
  -g python-flask \
  -o /src/sdk/python/generated \
  --package-name "plugin"

## Move files
mv $PYTHONSDKPATH/generated/plugin $PYTHONSDKPATH/plugin
mv $PYTHONSDKPATH/generated/.dockerignore $PYTHONSDKPATH/.dockerignore
mv $PYTHONSDKPATH/generated/.gitignore $PYTHONSDKPATH/.gitignore
mv $PYTHONSDKPATH/generated/requirements.txt $PYTHONSDKPATH/requirements.txt
mv $PYTHONSDKPATH/generated/setup.py $PYTHONSDKPATH/setup.py

## Cleanup
rm -rf $PYTHONSDKPATH/generated
rm -rf $PYTHONSDKPATH/plugin/controllers
rm -rf $PYTHONSDKPATH/plugin/test
rm -rf $PYTHONSDKPATH/plugin/openapi
rm -rf $PYTHONSDKPATH/plugin/__main__.py
rm -rf $PYTHONSDKPATH/plugin/encoder.py
