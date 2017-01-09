// Copyright (c) 2016 Paul Jolly <paul@myitcv.org.uk>, all rights reserved.
// Use of this document is governed by a license found in the LICENSE document.

package gogenerate

import (
	"bufio"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

// These constants correspond in name and value to the details given in
// go generate --help
const (
	GOARCH    = "GOARCH"
	GOOS      = "GOOS"
	GOFILE    = "GOFILE"
	GOLINE    = "GOLINE"
	GOPACKAGE = "GOPACKAGE"

	GoGeneratePrefix = "//go:generate"
)

const (
	FlagLog = "gglog"

	LogInfo    = "info"
	LogWarning = "warning"
	LogError   = "error"
	LogFatal   = "fatal"
)

func LogFlag() *string {
	return flag.String(FlagLog, LogFatal, "log level; one of info, warning, error, fatal")
}

func commentRegex(command string) (*regexp.Regexp, error) {
	// notice we make the trailing space or newline optional here.... because
	// when we read a file line by line using a scanner, the read line is stripped
	// of its \n
	return regexp.Compile(`\A` + GoGeneratePrefix + ` +` + command + `(?:\n| .*\n?)?\z`)
}

func FilesContainingCmd(path string, command string) ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	// files is sorted....

	matching, err := commentRegex(command)
	if err != nil {
		return nil, err
	}

	var matches []os.FileInfo

	for _, fi := range files {
		if fi.IsDir() {
			continue
		}

		fn := filepath.Join(path, fi.Name())

		file, err := os.Open(fn)
		if err != nil {
			return nil, err
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if matching.MatchString(scanner.Text()) {
				matches = append(matches, fi)
				break
			}
		}

		file.Close()

		if err := scanner.Err(); err != nil {
			return nil, err
		}
	}

	return matches, nil
}
