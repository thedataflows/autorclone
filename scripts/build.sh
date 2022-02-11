#!/bin/env bash

set -e -o pipefail

PROJECT_NAME=`grep module ${0%/*}/../go.mod`
PROJECT_NAME=${PROJECT_NAME##*/}
EXT=
[[ "$OS" == "Windows_NT" ]] && EXT='.exe'

(
  cd ${0%/*}/../
  go build -ldflags "-s -w -X main.version=`git tag --sort=-version:refname | head -n 1`" -trimpath -o $PROJECT_NAME$EXT cmd/$PROJECT_NAME/main.go
)
