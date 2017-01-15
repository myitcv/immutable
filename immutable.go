// Copyright (c) 2016 Paul Jolly <paul@myitcv.org.uk>, all rights reserved.
// Use of this document is governed by a license found in the LICENSE document.

package immutable

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/myitcv/immutable/internal/util"
)

const (
	// ImmTypeTemplPrefix is the prefix used to identify immutable type templates
	ImmTypeTemplPrefix = util.ImmTypeTemplPrefix

	// Pkg is the import path of this package
	PkgImportPath = "github.com/myitcv/immutable"

	// ImmTypeIdentifier should not be used; instead consider using IsImmType
	ImmTypeIdentifier = PkgImportPath + ":ImmutableType"
)

// Immutable is the interface implemented by all immutable types. If Go had generics the interface would
// be defined, assuming a generic type parameter T, as follows:
//
// 	type Immutable<T> interface {
// 		AsMutable() T
// 		AsImmutable() T
// 		WithMutations(f func(v T)) T
// 		Mutable() bool
// 	}
//
// Because we do not have such a type parameter we can only define the Mutable() method in the interface
type Immutable interface {
	Mutable() bool
}

// IsImmTmpl determines whether the supplied declaration is an immutable template type (either a struct,
// slice or map), returning the name of the type with the ImmTypeTemplPrefix removed in that case
func IsImmTmpl(d ast.Decl) bool {
	gd, _, _ := util.IsImmTmpl(d)

	return gd != nil
}

// IsImmType confirms whether the supplied declaration, which has to be found within the supplied
// file, is the result of immutable generation
func IsImmType(file *ast.File, d ast.Decl) string {
	if d.Pos() > file.End() || d.Pos() < file.Pos() {
		panic(fmt.Errorf("Declaration within supplied file"))
	}

	gd, ok := d.(*ast.GenDecl)
	if !ok || gd.Tok != token.TYPE {
		return ""
	}

	if len(gd.Specs) != 1 {
		panic("@myitcv needs to better understand go/ast")
	}

	ts := gd.Specs[0].(*ast.TypeSpec)

	st, ok := ts.Type.(*ast.StructType)
	if !ok {
		return ""
	}

	// at this point we need to find the first comment in the struct
	// body
	typName := ts.Name.Name

	// we need to find the first comment group in the struct
	// before the first field

	// all our implementations have fields
	if st.Fields.NumFields() == 0 {
		return ""
	}

	ffPos := st.Fields.List[0].Pos()

	for _, cg := range file.Comments {
		if cg.Pos() > ffPos {
			break
		}

		if cg.Pos() > st.Pos() {
			// the comment will have a trailing newline

			cl := strings.TrimRight(cg.Text(), "\n")

			if cl == ImmTypeIdentifier {
				return typName
			}
		}
	}

	return ""
}
