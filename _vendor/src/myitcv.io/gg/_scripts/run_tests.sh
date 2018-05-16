#!/usr/bin/env bash

# Copyright (c) 2016 Paul Jolly <paul@myitcv.org.uk>, all rights reserved.
# Use of this document is governed by a license found in the LICENSE document.

source "${BASH_SOURCE%/*}/common.bash"

export GOPATH=$PWD/_vendor:$GOPATH
export PATH=$GOPATH/bin:$PATH

go generate ./...
go test -count=1 ./...
go vet ./...
