// Copyright (c) 2016 Paul Jolly <paul@myitcv.org.uk>, all rights reserved.
// Use of this document is governed by a license found in the LICENSE document.

// Package gogenerate exposes some of the unexported internals of the go generate command as a convenience
// for the authors of go generate generators. See https://github.com/myitcv/gogenerate/wiki/Go-Generate-Notes
// for further notes on such generators.
package gogenerate

import (
	"flag"
	"go/build"
	"path/filepath"
	"sort"
	"strings"
)

// These constants correspond in name and value to the details given in
// go generate --help
const (
	GOARCH    = "GOARCH"
	GOFILE    = "GOFILE"
	GOLINE    = "GOLINE"
	GOOS      = "GOOS"
	GOPACKAGE = "GOPACKAGE"
	GOPATH    = "GOPATH"

	GoGeneratePrefix = "//go:generate"
)

const (
	// FlagLog is the name of the common flag shared between go generate generators
	// to control logging verbosity
	FlagLog = "gglog"
)

// The various log levels supported by the flag returned by LogFlag()
const (
	LogInfo    = "info"
	LogWarning = "warning"
	LogError   = "error"
	LogFatal   = "fatal"
)

// LogFlag defines a string flag named according to the constant FlagLog and returns a
// pointer to the string the flag sets
func LogFlag() *string {
	return flag.String(FlagLog, LogFatal, "log level; one of info, warning, error, fatal")
}

// FilesContainingCmd returns a sorted list of Go file names (defined by go list as
// GoFiles + CgoFiles + TestGoFiles + XTestGoFiles) in the directory dir that
// contain a command matching any of the provided commands after quote and
// variable expansion (as described by go generate -help). The file names will,
// by definition, be relative to dir
func FilesContainingCmd(dir string, commands ...string) ([]string, error) {

	// clean our commands
	nonZero := false
	for i, c := range commands {
		c = strings.TrimSpace(c)
		if c != "" {
			nonZero = true
		}
		commands[i] = c
	}
	if !nonZero {
		return nil, nil
	}

	pkg, err := build.ImportDir(dir, build.IgnoreVendor)
	if err != nil {
		return nil, err
	}

	// GoFiles+CgoFiles+TestGoFiles+XTestGoFiles per go list
	// these are all relative to path
	gofiles := make([]string, 0, len(pkg.GoFiles)+len(pkg.CgoFiles)+len(pkg.TestGoFiles)+len(pkg.XTestGoFiles))
	gofiles = append(gofiles, pkg.GoFiles...)
	gofiles = append(gofiles, pkg.CgoFiles...)
	gofiles = append(gofiles, pkg.TestGoFiles...)
	gofiles = append(gofiles, pkg.XTestGoFiles...)

	matchMap := make(map[string]struct{})

	for _, f := range gofiles {
		checkMatch := func(args []string) {
			for _, cmd := range commands {
				if args[0] == cmd {
					matchMap[f] = struct{}{}
				}
			}
		}

		fp := filepath.Join(dir, f)

		err = DirFunc(pkg.ImportPath, fp, checkMatch)

		if err != nil {
			return nil, err
		}
	}

	matches := make([]string, 0, len(matchMap))
	for k := range matchMap {
		matches = append(matches, k)
	}
	sort.Sort(byBase(matches))
	return matches, nil
}

type byBase []string

func (f byBase) Len() int           { return len(f) }
func (f byBase) Less(i, j int) bool { return filepath.Base(f[i]) < filepath.Base(f[j]) }
func (f byBase) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
