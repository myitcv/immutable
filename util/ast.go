package util // import "myitcv.io/immutable/util"

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"path"
	"strings"

	"myitcv.io/immutable"
)

const (
	debug = false
)

// IsImmTmplAst determines whether the supplied type spec is an immutable template type (either a struct,
// slice or map), returning the name of the type with the ImmTypeTmplPrefix removed in that case
func IsImmTmplAst(ts *ast.TypeSpec) (string, bool) {
	typName := ts.Name.Name

	if !strings.HasPrefix(typName, immutable.ImmTypeTmplPrefix) {
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

	name := strings.TrimPrefix(typName, immutable.ImmTypeTmplPrefix)

	return name, true
}

type ImmTypeAst interface {
	isImmTypeAst()
}

type (
	ImmTypeAstImplsIntf struct{}
	ImmTypeAstExtIntf   struct{}
	ImmTypeAstBasic     struct{}
	ImmTypeAstSpecial   struct{}
	ImmTypeAstStruct    struct{}
	ImmTypeAstMap       struct {
		Key  ast.Expr
		Elem ast.Expr
	}
	ImmTypeAstSlice struct {
		Elem ast.Expr
	}
)

func (i ImmTypeAstImplsIntf) isImmTypeAst() {}
func (i ImmTypeAstExtIntf) isImmTypeAst()   {}
func (i ImmTypeAstBasic) isImmTypeAst()     {}
func (i ImmTypeAstSpecial) isImmTypeAst()   {}
func (i ImmTypeAstStruct) isImmTypeAst()    {}
func (i ImmTypeAstMap) isImmTypeAst()       {}
func (i ImmTypeAstSlice) isImmTypeAst()     {}

type PkgLookup func(imp *ast.ImportSpec) (*ast.Package, error)

type Checker struct {
	astTypeCache map[string]ImmTypeAst
	pkgLookup    PkgLookup
	fset         *token.FileSet
}

func NewChecker(pkgLookup PkgLookup, fset *token.FileSet) Checker {
	return Checker{
		astTypeCache: make(map[string]ImmTypeAst),
		pkgLookup:    pkgLookup,
		fset:         fset,
	}
}

// IsImmTypeAst determines by syntax tree analysis alone whether the supplied
// ast.Expr represents an immutable type. In case a type is immutable, a value
// of type ImmTypeAstStruct, ImmTypeAstSlice or ImmTypeAstMap. In case the type
// is "implements" the full immutable "interface" but neither of the
// aforementioned instances, ImmTypeAstImplsInt is returned. For special types
// like time.Time, ImmTypeAstSpecial is returned. For basic types,
// ImmTypeAstBasic is returned.  If a type is a reference to an interface type
// that extends immutable.Immutable then ImmTypeAstExtIntf is returned.  If a
// type is not immutable then nil is returned.
func (c Checker) IsImmTypeAst(ts ast.Expr, file *ast.File, pkg string) (ImmTypeAst, error) {
	// The only way the provided expression can "be" an immutable type is when
	// it is a Type reference (per the spec) and that type "implements" the
	// immutable "interface"
	// (https://myitcv.io/immutable/wiki/immutableGen)

	isPointer := false
	pkgStr := ""
	typStr := ""

	if v, ok := ts.(*ast.ParenExpr); ok {
		ts = v.X
	}

	if v, ok := ts.(*ast.StarExpr); ok {
		isPointer = true
		ts = v.X
	}

	imps := file.Imports

	switch ts := ts.(type) {
	case *ast.Ident:
		switch ts.Name {
		case "bool", "int", "int8", "int16", "int32", "int64",
			"uint", "uint8", "uint16", "uint32", "uint64", "uintptr",
			"float32", "float64", "complex64", "complex128", "string":

			return ImmTypeAstBasic{}, nil
		}

		pkgStr = pkg
		typStr = ts.Name

	case *ast.SelectorExpr:
		if x, ok := ts.X.(*ast.Ident); ok {

			typStr = ts.Sel.Name

			foundImp := false

			for _, i := range imps {
				p := strings.Trim(i.Path.Value, "\"")

				toCheck := path.Base(p)

				if i.Name != nil {
					toCheck = i.Name.Name
				}

				if x.Name == toCheck {
					pkgStr = p
					foundImp = true

					break
				}
			}

			if foundImp {
				break
			}
		}

		// we failed to properly resolve the selector expr
		return nil, nil

	default:
		return nil, nil
	}

	// special cases...
	key := buildKey(pkgStr, typStr, isPointer)

	if key == "time.Time" {
		return ImmTypeAstSpecial{}, nil
	}

	return c.isAstTypeImm(pkgStr, typStr, isPointer)
}

func (c Checker) isAstTypeImm(pkgStr, typStr string, isPointer bool) (ImmTypeAst, error) {

	key := buildKey(pkgStr, typStr, isPointer)

	if v, ok := c.astTypeCache[key]; ok {
		return v, nil
	}

	pkg, err := c.pkgLookup(pkgStr)
	if err != nil {
		return nil, err
	}

	// set initially to allow for early return
	// when false
	var res ImmTypeAst
	defer func() {
		c.astTypeCache[key] = res
	}()

	var types []*ast.TypeSpec
	meths := make(map[string]*ast.FuncDecl)

	for _, f := range pkg.Files {
		for _, d := range f.Decls {
			switch d := d.(type) {
			case *ast.FuncDecl:
				if d.Recv == nil {
					continue
				}

				f := d.Recv.List[0]

				switch v := f.Type.(type) {
				case *ast.StarExpr:
					id, ok := v.X.(*ast.Ident)
					if !ok {
						continue
					}

					if isPointer && id.Name == typStr {
						meths[d.Name.Name] = d
					}
				case *ast.Ident:
					if !isPointer && v.Name == typStr {
						meths[d.Name.Name] = d
					}
				}

			case *ast.GenDecl:
				if d.Tok != token.TYPE {
					continue
				}

				for _, ts := range d.Specs {
					ts := ts.(*ast.TypeSpec)

					if ts.Name.Name != typStr {
						continue
					}

					types = append(types, ts)

					// now we can quickly check whether this type is one of
					// immutableGen's results

					var st *ast.StructType

					switch t := ts.Type.(type) {
					case *ast.InterfaceType:

						ok, err := c.interfaceExtendsImmutable(pkgStr, t, f.Imports)
						if err != nil {
							return nil, err
						}
						if ok {
							res = ImmTypeAstExtIntf{}
						}

						return res, nil

					case *ast.StructType:
						st = t
						// continues below...

					default:
						continue
					}

					var key, val ast.Expr

					foundTheMap := false
					foundTheSlice := false
					foundMutable := false
					foundTmpl := false

				NextField:
					for _, f := range st.Fields.List {
						if len(f.Names) != 1 {
							continue NextField
						}

						fn := f.Names[0].Name

						switch fn {
						case "theMap":
							switch t := f.Type.(type) {
							case *ast.MapType:
								key = t.Key
								val = t.Value
							default:
								continue NextField
							}

							foundTheMap = true

						case "theSlice":
							switch t := f.Type.(type) {
							case *ast.ArrayType:
								if t.Len != nil {
									continue NextField
								}

								val = t.Elt
							default:
								continue NextField
							}

							foundTheSlice = true

						case "__tmpl":
							n, ok := f.Type.(*ast.Ident)
							if ok && n.Name == immutable.ImmTypeTmplPrefix+typStr {
								foundTmpl = true
							}

						case "mutable":
							n, ok := f.Type.(*ast.Ident)
							if ok && n.Name == "bool" {
								foundMutable = true
							}
						}
					}

					if !foundMutable || !foundTmpl {
						continue
					}

					switch {
					case foundTheMap:
						res = ImmTypeAstMap{
							Key:  key,
							Elem: val,
						}

					case foundTheSlice:
						res = ImmTypeAstSlice{
							Elem: val,
						}
					default:
						res = ImmTypeAstStruct{}
					}

					return res, nil
				}
			}
		}
	}

	if len(types) != 1 {
		return nil, nil
	}

	fullTypStr := typStr

	if isPointer {
		fullTypStr = "*" + fullTypStr
	}

	res = c.astImplsImm(fullTypStr, meths)

	return res, nil
}

func (c Checker) lookupIntf(pkgStr, intf string) (*ast.InterfaceType, []*ast.ImportSpec, error) {
	pkg, err := c.pkgLookup(pkgStr)
	if err != nil {
		return nil, nil, err
	}

	for _, f := range pkg.Files {
		for _, d := range f.Decls {
			switch t := d.(type) {
			case *ast.GenDecl:
				if t.Tok != token.TYPE {
					continue
				}

				for _, s := range t.Specs {
					ts := s.(*ast.TypeSpec)

					switch v := ts.Type.(type) {
					case *ast.InterfaceType:
						if ts.Name.Name == intf {
							return v, f.Imports, nil
						}
					}
				}
			}
		}
	}

	return nil, nil, nil
}

func (c Checker) interfaceExtendsImmutable(pkgStr string, intf *ast.InterfaceType, imps []*ast.ImportSpec) (bool, error) {

	for _, f := range intf.Methods.List {
		if len(f.Names) != 0 {
			// method
			continue
		}

		// embedded
		switch v := f.Type.(type) {
		case *ast.Ident:
			ni, _, err := c.lookupIntf(pkgStr, v.Name)
			if err != nil {
				return false, err
			}

			ok, err := c.interfaceExtendsImmutable(pkgStr, ni, imps)
			if err != nil {
				return false, err
			}
			if ok {
				return true, nil
			}

		case *ast.SelectorExpr:
			pn := v.X.(*ast.Ident).Name

			for _, i := range imps {
				p := strings.Trim(i.Path.Value, "\"")

				toCheck := path.Base(p)

				if i.Name != nil {
					toCheck = i.Name.Name
				}

				if pn == toCheck {
					ni, nimps, err := c.lookupIntf(p, v.Sel.Name)
					if err != nil {
						return false, err
					}

					if ni != nil {
						if p == immutable.PkgImportPath && v.Sel.Name == "Immutable" {
							return true, nil
						}

						ok, err := c.interfaceExtendsImmutable(p, ni, nimps)
						if err != nil {
							return false, err
						}

						if ok {
							return true, nil
						}
					}
				}
			}

		default:
			panic(fmt.Errorf("Check %v; it contains an interface we cannot walk", pkgStr))
		}

	}

	return false, nil
}

func buildKey(pkgStr, typStr string, isPointer bool) string {

	key := typStr

	if pkgStr != "" {
		key = pkgStr + "." + key
	}

	if isPointer {
		key = "*" + key
	}

	return key
}

func (c Checker) astString(node interface{}) string {
	b := bytes.NewBuffer(nil)

	err := printer.Fprint(b, c.fset, node)
	if err != nil {
		panic(fmt.Errorf("failed to print node %v: %v", node, err))
	}

	return b.String()
}

func (c Checker) astImplsImm(typStr string, meths map[string]*ast.FuncDecl) ImmTypeAst {
	// Need to check for the presence of the methods defined here:
	// https://myitcv.io/immutable/wiki/immutableGen

	// TODO use shared constants with immutableGen

	// Mutable() bool
	{
		m, ok := meths["Mutable"]
		if !ok {
			return nil
		}

		sig := m.Type

		if sig.Params.NumFields() != 0 || sig.Results.NumFields() != 1 {
			return nil
		}

		f, ok := sig.Results.List[0].Type.(*ast.Ident)
		if !ok || f.Name != "bool" {
			return nil
		}
	}

	// AsMutable() *T
	{
		m, ok := meths["AsMutable"]
		if !ok {
			return nil
		}

		sig := m.Type

		if sig.Params.NumFields() != 0 || sig.Results.NumFields() != 1 {
			return nil
		}

		f := c.astString(sig.Results.List[0].Type)
		if f != typStr {
			return nil
		}
	}

	// AsImmutable() *T
	{
		m, ok := meths["AsImmutable"]
		if !ok {
			return nil
		}
		sig := m.Type

		if sig.Params.NumFields() != 1 || sig.Results.NumFields() != 1 {
			return nil
		}

		p := c.astString(sig.Params.List[0].Type)
		if p != typStr {
			return nil
		}

		f := c.astString(sig.Results.List[0].Type)
		if f != typStr {
			return nil
		}
	}

	// WithMutable(f func(t *T)) *T
	{
		m, ok := meths["WithMutable"]
		if !ok {
			return nil
		}
		sig := m.Type

		if sig.Params.NumFields() != 1 || sig.Results.NumFields() != 1 {
			return nil
		}

		pf, ok := sig.Params.List[0].Type.(*ast.FuncType)
		if !ok || pf.Params.NumFields() != 1 || pf.Results.NumFields() != 0 {
			return nil
		}

		p := c.astString(pf.Params.List[0].Type)
		if p != typStr {
			return nil
		}

		f := c.astString(sig.Results.List[0].Type)
		if f != typStr {
			return nil
		}
	}

	// WithImmutable(f func(t *T)) *T
	{
		m, ok := meths["WithImmutable"]
		if !ok {
			return nil
		}
		sig := m.Type

		if sig.Params.NumFields() != 1 || sig.Results.NumFields() != 1 {
			return nil
		}

		pf, ok := sig.Params.List[0].Type.(*ast.FuncType)
		if !ok || pf.Params.NumFields() != 1 || pf.Results.NumFields() != 0 {
			return nil
		}

		p := c.astString(pf.Params.List[0].Type)
		if p != typStr {
			return nil
		}

		f := c.astString(sig.Results.List[0].Type)
		if f != typStr {
			return nil
		}
	}

	return ImmTypeAstImplsIntf{}
}

func debugf(format string, args ...interface{}) {
	if debug {
		fmt.Printf(format, args...)
	}
}
