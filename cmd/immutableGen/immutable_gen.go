package main

import (
	"flag"
	"io"
	"os"

	"github.com/myitcv/immutable/cmd/immutableGen/internal/core"
)

const (
	_GoFile    = "GOFILE"
	_GoPackage = "GOPACKAGE"
)

var fLicenseFile = flag.String("licenseFile", "", "file containing an uncommented license header")

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

	var licenseFile io.Reader

	if *fLicenseFile != "" {
		lf, err := os.Open(*fLicenseFile)
		if err != nil {
			panic(err)
		}

		licenseFile = lf
	}

	err := core.Execute(envFile, envPkg, licenseFile)
	if err != nil {
		panic(err)
	}
}
