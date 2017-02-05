// Copyright (c) 2016 Paul Jolly <paul@myitcv.org.uk>, all rights reserved.
// Use of this document is governed by a license found in the LICENSE document.

package immutable

import (
	"go/ast"
	"strings"
)

const (
	// ImmTypeTmplPrefix is the prefix used to identify immutable type templates
	ImmTypeTmplPrefix = "_Imm_"

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
// slice or map), returning the name of the type with the ImmTypeTmplPrefix removed in that case
func IsImmTmpl(ts *ast.TypeSpec) (string, bool) {
	typName := ts.Name.Name

	if !strings.HasPrefix(typName, ImmTypeTmplPrefix) {
		return "", false
	}

	valid := false

	switch typ := ts.Type.(type) {
	case *ast.MapType:
		valid = true
	case *ast.ArrayType:
		if typ.Len == nil {
			valid = true
		}
	case *ast.StructType:
		valid = true
	}

	if !valid {
		return "", false
	}

	name := strings.TrimPrefix(typName, ImmTypeTmplPrefix)

	return name, true
}
