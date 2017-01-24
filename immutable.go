// Copyright (c) 2016 Paul Jolly <paul@myitcv.org.uk>, all rights reserved.
// Use of this document is governed by a license found in the LICENSE document.

package immutable

import (
	"fmt"
	"go/ast"
	"strings"
)

const (
	// ImmTypeTemplPrefix is the prefix used to identify immutable type templates
	ImmTypeTemplPrefix = "_Imm_"

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

// IsImmTmpl determines whether the supplied type spec is an immutable template type (either a struct,
// slice or map), returning the name of the type with the ImmTypeTemplPrefix removed in that case
func IsImmTmpl(ts *ast.TypeSpec) (string, bool) {
	typName := ts.Name.Name

	if !strings.HasPrefix(typName, ImmTypeTemplPrefix) {
		return "", false
	}

	name := strings.TrimPrefix(typName, ImmTypeTemplPrefix)

	return name, true
}

// IsImmType confirms whether the supplied declaration, which has to be found within the supplied
// file, is the result of immutable generation
func IsImmType(file *ast.File, ts *ast.TypeSpec) bool {
	if ts.Pos() > file.End() || ts.Pos() < file.Pos() {
		panic(fmt.Errorf("Declaration within supplied file"))
	}

	st, ok := ts.Type.(*ast.StructType)
	if !ok {
		return false
	}

	// we need to find the first comment group in the struct
	// before the first field

	// all our implementations have fields
	if st.Fields.NumFields() == 0 {
		return false
	}

	ffPos := st.Fields.List[0].Pos()

	for _, cg := range file.Comments {
		// is the comment group after the first field?
		if cg.Pos() > ffPos {
			break
		}

		if cg.Pos() > st.Pos() {
			for _, c := range cg.List {
				if c.Text == "//"+ImmTypeIdentifier {
					return true
				}
			}

			// now we can break because by definition we've exhausted
			// the first comment group (of which there can be only one)
			break
		}
	}

	return false
}
