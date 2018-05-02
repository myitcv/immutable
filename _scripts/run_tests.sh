#!/usr/bin/env bash

# Copyright (c) 2016 Paul Jolly <paul@myitcv.org.uk>, all rights reserved.
# Use of this document is governed by a license found in the LICENSE document.

source "${BASH_SOURCE%/*}/common.bash"

export GOPATH=$PWD/_vendor:$GOPATH
export PATH=$GOPATH/bin:$PATH

rm -f !(_vendor)/**/gen_*.go

go install myitcv.io/immutable/cmd/immutableGen
go install myitcv.io/immutable/cmd/immutableVet

pushd cmd/immutableVet/_testFiles > /dev/null
go generate
popd > /dev/null

go generate ./...
go install ./...
go test myitcv.io/immutable/cmd/immutableGen/internal/coretest
go test ./...
immutableVet myitcv.io/immutable/example
