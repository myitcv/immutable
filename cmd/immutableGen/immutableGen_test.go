// Copyright (c) 2016 Paul Jolly <paul@myitcv.org.uk>, all rights reserved.
// Use of this document is governed by a license found in the LICENSE document.

package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/myitcv/immutable"
)

const (
	TestFiles = "_testFiles"
)

func TestBasic(t *testing.T) {
	license := "My favourite license"
	echoCmd := `echo "hello world"` // need a command that will succeed with zero exit code

	genTarget := "gen_main_immutableGen.go"

	execute(TestFiles, "main", license, gogenCmds{echoCmd})

	genFile := filepath.Join(TestFiles, genTarget)

	genOut, err := os.Open(genFile)
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(os.Stdout, genOut)
	if err != nil {
		panic(err)
	}

	_, err = genOut.Seek(0, 0)
	if err != nil {
		panic(err)
	}

	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, genTarget, genOut, parser.AllErrors|parser.ParseComments)
	if err != nil {
		panic(err)
	}

	foundMyStruct := false
	foundMySlice := false
	foundMyMap := false

	for _, d := range f.Decls {

		name := immutable.IsImmType(f, d)

		if name != "" {
			fmt.Println(">>", name)
		}

		switch name {
		case "MyStruct":
			foundMyStruct = true
		case "MySlice":
			foundMySlice = true
		case "MyMap":
			foundMyMap = true
		}

		if name == "MyStruct" {
			foundMyStruct = true
		}
	}

	if !foundMyStruct {
		t.Errorf("did not find MyStruct in generated output")
	}
	if !foundMySlice {
		t.Errorf("did not find MySlice in generated output")
	}
	if !foundMyMap {
		t.Errorf("did not find myMap in generated output")
	}
}
