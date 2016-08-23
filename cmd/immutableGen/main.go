package main

import (
	"flag"
	"os"

	"github.com/myitcv/immutable/cmd/immutableGen/internal/generator"
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

	var licenseFile *os.File

	if *fLicenseFile != "" {
		lf, err := os.Open(*fLicenseFile)
		if err != nil {
			panic(err)
		}

		licenseFile = lf
	}

	err := generator.DoIt(envFile, envPkg, licenseFile)
	if err != nil {
		panic(err)
	}
}
