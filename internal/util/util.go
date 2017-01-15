package util

import (
	"go/ast"
	"go/token"
	"strings"
)

const (
	ImmTypeTemplPrefix = "_Imm_"
)

// IsImmTmpl determines whether the supplied declaration is an immutable template type (either a struct,
// slice or map), returning the declaration, type spec and name of the type with the ImmTypeTemplPrefix
// removed in that case
func IsImmTmpl(d ast.Decl) (*ast.GenDecl, *ast.TypeSpec, string) {
	gd, ok := d.(*ast.GenDecl)
	if !ok || gd.Tok != token.TYPE {
		return nil, nil, ""
	}

	if len(gd.Specs) != 1 {
		panic("@myitcv needs to better understand go/ast")
	}

	ts := gd.Specs[0].(*ast.TypeSpec)

	typName := ts.Name.Name

	if !strings.HasPrefix(typName, ImmTypeTemplPrefix) {
		return nil, nil, ""
	}

	name := strings.TrimPrefix(typName, ImmTypeTemplPrefix)

	return gd, ts, name
}
