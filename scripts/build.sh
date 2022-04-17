#!/usr/bin/env bash

set -e -o pipefail

PROJECT_PATH=`readlink -f ${0%/*}/..`

PROJECT_NAME=`grep module $PROJECT_PATH/go.mod`
PROJECT_NAME=${PROJECT_NAME##*/}
EXT=
[[ "$OS" == "Windows_NT" || "$GOOS" == "windows" ]] && EXT='.exe'

(
  cd ${0%/*}/../
  TAG=`git tag --sort=-version:refname | head -n 1`
  go build \
    -ldflags "-s -w -X 'dataflows.com/autorclone/internal/pkg/autorclone.version=$TAG'" \
    -trimpath \
    -o $PROJECT_PATH/bin/$PROJECT_NAME$EXT \
    $PROJECT_PATH/cmd/$PROJECT_NAME/main.go
)
