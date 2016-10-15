package main

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/myitcv/immutable/cmd/immutableGen/internal/core"
)

const (
	TestFiles = "_testFiles"
)

func TestBasic(t *testing.T) {
	license := "My favourite license"
	trueCmd := "true" // need a command that will succeed with zero exit code

	target := filepath.Join(TestFiles, "main.go")

	err := core.Execute(target, "main", license, core.GoGenCmds{trueCmd})

	if err != nil {
		t.Fatalf("Err should have been nil: %v\n", err)
	}

	genFile := filepath.Join(TestFiles, "gen_main_immutable.go")

	genOut, err := os.Open(genFile)
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(os.Stdout, genOut)
	if err != nil {
		panic(err)
	}

}
