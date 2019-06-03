#!/bin/bash

pushd $(dirname "${0}") > /dev/null
basedir=$(pwd -L)
# Use "pwd -P" for the path without links. man bash for more info.
popd > /dev/null

echo "Setting go path to: ${basedir}"

export GOPATH="${basedir}:$GOPATH"