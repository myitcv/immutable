package main

import (
	"os"

	"github.com/myitcv/immutable/cmd/immutableGen/internal/generator"
)

const (
	_GoFile    = "GOFILE"
	_GoPackage = "GOPACKAGE"
)

func main() {
	envFile, ok := os.LookupEnv(_GoFile)
	if !ok {
		panic("Env not correct; missing " + _GoFile)
	}

	envPkg, ok := os.LookupEnv(_GoPackage)
	if !ok {
		panic("Env not correct; missing " + _GoPackage)
	}

	err := generator.DoIt(envFile, envPkg)
	if err != nil {
		panic(err)
	}
}
