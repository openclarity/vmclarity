#!/usr/bin/env bash
# This script is meant to build and compile every protocolbuffer for each
# service declared in this repository (as defined by sub-directories).
# It compiles using docker containers based on Namely's protoc image
# seen here: https://github.com/namely/docker-protoc

set -e

# Helper for adding a directory to the stack and echoing the result
function enterDir {
  echo "Entering $1"
  pushd $1 > /dev/null
}

# Helper for popping a directory off the stack and echoing the result
function leaveDir {
  echo "Leaving `pwd`"
  popd > /dev/null
}

# Enters the directory and starts the build / compile process for the services
# protobufs
function buildDir {
  currentDir="$1"
  echo "Building directory \"$currentDir\""

  enterDir $currentDir

  buildProtoForTypes $currentDir

  leaveDir
}

# Iterates through all of the languages listed in the services .protolangs file
# and compiles them individually
function buildProtoForTypes {
  target=${1%/}

  if [ -f .protolangs ]; then
    while read lang; do
      folder="$target-$lang"

      # Use the docker container for the language we care about and compile
      docker run -v `pwd`:/defs namely/protoc-$lang

      # Copy the generated files out of the pb-* path into the repository
      # that we care about and cleanup
      mkdir -p $folder
      cp -R pb-$lang/* $folder/
      rm -rf pb-$lang
    done < .protolangs
  fi
}

# Finds all directories in the repository and iterates through them calling the
# compile process for each one
function buildAll {
  echo "Buidling service's protocol buffers"
  for d in */; do
    buildDir $d
  done
}

enterDir grpc
buildAll
leaveDir
