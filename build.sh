#!/bin/bash

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
export GOPATH="$SCRIPTPATH"
go build -o grm grm
chmod +x grm