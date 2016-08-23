package main

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/myitcv/immutable/cmd/immutableGen/internal/generator"
)

const (
	TestFiles = "_testFiles"
)

func TestBasic(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "immutableGenTest")
	if err != nil {
		panic(err)
	}
	// fmt.Printf("TmpDir: %v\n", tmpDir)
	defer os.RemoveAll(tmpDir)

	cp := exec.Command("cp", "-r", TestFiles, tmpDir)
	res, err := cp.CombinedOutput()
	if err != nil {
		panic(string(res))
	}

	target := filepath.Join(tmpDir, TestFiles, "main.go")

	err = generator.DoIt(target, "main")

	if err != nil {
		t.Fatalf("Err should have been nil: %v\n", err)
	}

	genFile := filepath.Join(tmpDir, TestFiles, "gen_main_immutable.go")

	genOut, err := os.Open(genFile)
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(os.Stdout, genOut)
	if err != nil {
		panic(err)
	}

}
