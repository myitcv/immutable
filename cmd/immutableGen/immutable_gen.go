// Copyright (c) 2016 Paul Jolly <paul@myitcv.org.uk>, all rights reserved.
// Use of this document is governed by a license found in the LICENSE document.

package main

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/myitcv/immutable/cmd/immutableGen/internal/core"
)

const (
	_GoFile    = "GOFILE"
	_GoPackage = "GOPACKAGE"
)

var (
	fLicenseFile = flag.String("licenseFile", "", "file containing an uncommented license header")
	fGoGenCmds   core.GoGenCmds
)

func init() {
	flag.Var(&fGoGenCmds, "G", "Path to search for imports (flag can be used multiple times)")
}

func main() {
	flag.Parse()

	envFile, ok := os.LookupEnv(_GoFile)
	if !ok {
		panic("Env not correct; missing " + _GoFile)
	}

	envPkg, ok := os.LookupEnv(_GoPackage)
	if !ok {
		panic("Env not correct; missing " + _GoPackage)
	}

	licenseHeader := ""

	if *fLicenseFile != "" {
		byts, err := ioutil.ReadFile(*fLicenseFile)
		if err != nil {
			panic(err)
		}

		licenseHeader = string(byts)
	}

	err := core.Execute(envFile, envPkg, licenseHeader, fGoGenCmds)
	if err != nil {
		panic(err)
	}
}
