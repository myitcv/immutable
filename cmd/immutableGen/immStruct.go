package main

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/myitcv/immutable"
)

type immStruct struct {
	fset *token.FileSet

	name string
	dec  *ast.GenDecl
	st   *ast.StructType
}

func (o *output) genImmStructs(structs []immStruct) {
	type genField struct {
		Name string
		Type string
		f    *ast.Field
	}

	for _, s := range structs {

		o.printCommentGroup(s.dec.Doc)
		o.printImmPreamble(s.name, s.st)

		// start of struct
		o.pfln("type %v struct {", s.name)
		o.pfln("\t//%v", immutable.ImmTypeIdentifier)

		o.printLeadSpecCommsFor(s.st)

		o.pln("")

		var fields []genField

		for _, f := range s.st.Fields.List {
			names := ""
			sep := ""

			typ := o.exprString(f.Type)

			tag := ""

			if f.Tag != nil {
				tag = f.Tag.Value
			}

			for _, n := range f.Names {
				names = names + sep + fieldHidingPrefix + n.Name
				sep = ", "

				fields = append(fields, genField{
					Name: n.Name,
					Type: typ,
					f:    f,
				})

			}
			o.pfln("%v %v %v", names, typ, tag)
		}

		o.pln("")
		o.pln("mutable bool")

		// end of struct
		o.pfln("}")

		o.pln()

		o.pfln("var _ immutable.Immutable = &%v{}", s.name)
		o.pln()

		exp := exporter(s.name)

		o.pt(`
		func (s *{{.}}) AsMutable() *{{.}} {
			res := *s
			res.mutable = true
			return &res
		}

		func (s *{{.}}) AsImmutable() *{{.}} {
			s.mutable = false
			return s
		}

		func (s *{{.}}) Mutable() bool {
			return s.mutable
		}

		func (s *{{.}}) WithMutations(f func(si *{{.}})) *{{.}} {
			res := s.AsMutable()
			f(res)
			res = res.AsImmutable()

			return res
		}
		`, exp, s.name)

		for _, f := range fields {
			tmpl := struct {
				TypeName string
				Field    genField
			}{
				TypeName: s.name,
				Field:    f,
			}

			exp := exporter(f.Name)

			o.printCommentGroup(f.f.Doc)

			o.pt(`
			func (s *{{.TypeName}}) {{.Field.Name}}() {{.Field.Type}} {
				return s.`+fieldHidingPrefix+`{{.Field.Name}}
			}

			// {{Export "Set"}}{{Capitalise .Field.Name}} is the setter for {{Capitalise .Field.Name}}()
			func (s *{{.TypeName}}) {{Export "Set"}}{{Capitalise .Field.Name}}(n {{.Field.Type}}) *{{.TypeName}} {
				if s.mutable {
					s.`+fieldHidingPrefix+`{{.Field.Name}} = n
					return s
				}

				res := *s
				res.`+fieldHidingPrefix+`{{.Field.Name}} = n
				return &res
			}
			`, exp, tmpl)
		}
	}
}

func (o *output) printLeadSpecCommsFor(st *ast.StructType) {

	var end token.Pos

	// we are looking for comments before the first field (if there is one)

	if f := st.Fields; f != nil && len(f.List) > 0 {
		end = f.List[0].End()
	} else {
		end = st.End()
	}

	for _, cg := range o.curFile.Comments {
		if cg.Pos() > st.Pos() && cg.End() < end {
			for _, c := range cg.List {
				if strings.HasPrefix(c.Text, "//") && !strings.HasPrefix(c.Text, "// ") {
					o.pln(c.Text)
				}
			}
		}
	}

}
