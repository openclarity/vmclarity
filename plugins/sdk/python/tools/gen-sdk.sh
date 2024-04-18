#!/usr/bin/env bash

## Consts
OPENAPI_GENERATOR_IMAGE="openapitools/openapi-generator-cli:v7.5.0"

## Script vars
SCRIPT_PATH="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
PLUGIN_PATH="$(realpath $SCRIPT_PATH/../../..)"
SDK_PATH="$(realpath $SCRIPT_PATH/..)"

## Generate code
echo "Generating Python server code from $PLUGIN_PATH"
docker run --rm \
  -u $(id -u ${USER}):$(id -g ${USER}) \
  -v ${PLUGIN_PATH}:/src \
  $OPENAPI_GENERATOR_IMAGE \
  generate \
  -i /src/openapi.yaml \
  -g python-flask \
  -o /src/sdk/python/generated \
  --package-name "plugin" \
  --additional-properties "legacyDiscriminatorBehavior=false"

## Do not overwrite base model
rm -rf $SDK_PATH/generated/plugin/models/base_model.py

## Move rest of autogenerated code
cp -R $SDK_PATH/generated/plugin/models/* $SDK_PATH/plugin/models

## Cleanup
rm -rf $SDK_PATH/generated
