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

	licenseHeader := ""

	if *fLicenseFile != "" {
		byts, err := ioutil.ReadFile(*fLicenseFile)
		if err != nil {
			panic(err)
		}

		licenseHeader = string(byts)
	}

	err := core.Execute(envFile, envPkg, licenseHeader)
	if err != nil {
		panic(err)
	}
}
