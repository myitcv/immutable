#!/bin/sh

# Copyright (c) 2016 Paul Jolly <paul@myitcv.org.uk>, all rights reserved.
# Use of this document is governed by a license found in the LICENSE document.

set -e

set -v

go generate ./...
go vet ./...

# no tests to run here...

cd cmd/immutableGen/_testFiles/

go generate
go test
go vet
