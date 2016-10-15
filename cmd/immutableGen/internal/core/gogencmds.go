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
