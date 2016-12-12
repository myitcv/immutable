// Copyright (c) 2016 Paul Jolly <paul@myitcv.org.uk>, all rights reserved.
// Use of this document is governed by a license found in the LICENSE document.

package core

import "fmt"

type GoGenCmds []string

func (g *GoGenCmds) Set(value string) error {
	*g = append(*g, value)
	return nil
}

func (g *GoGenCmds) String() string {
	return fmt.Sprint(*g)
}
