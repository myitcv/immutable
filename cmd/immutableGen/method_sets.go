package main

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"
	"unicode"
	"unicode/utf8"

	"myitcv.io/immutable/util"
)

func (o *output) calcMethodSets() {

	for _, fts := range o.files {

		typeToString := func(t types.Type) string {
			return types.TypeString(t, func(p *types.Package) string {
				if p.Path() == o.pkgPath {
					return ""
				}

				for i := range fts.imports {
					ip := strings.Trim(i.Path.Value, "\"")
					if p.Path() == ip {
						if i.Name != nil {
							return i.Name.Name
						}
						return p.Name()
					}
				}

				newImport := &ast.ImportSpec{
					Path: &ast.BasicLit{Value: fmt.Sprintf(`"%v"`, p.Path())},
				}
				fts.imports[newImport] = struct{}{}

				return p.Name()
			})
		}
		possInvalidTypeToString := func(t types.Type) string {

			var tn string
			if typeIsInvalid(t) {
				// TODO: bit gross....
				var fte ast.Expr
				for e, tv := range o.info.Types {
					if tv.Type == t {
						if fte != nil {
							panic(fmt.Errorf("had two entries for the same invalid type; what to do?: %v and %v", fte, e))
						}
						fte = e
					}
				}
				if fte == nil {
					panic(fmt.Errorf("could not resolve expression for invalid type"))
				}
				tn = o.exprString(fte)
			} else {
				tn = typeToString(t)
			}
			return tn
		}

		for _, is := range fts.structs {
			debugf(">> calculating %v\n", is.name)

			seen := make(map[interface{}]bool)
			set := make(map[string]*field)
			possSet := make(map[string]*field)

			work := []embedded{{es: "*" + is.name}}
			var next []embedded
			var h embedded

			addPoss := func(name string, f field) {
				if _, ok := set[name]; !ok {
					if _, ok := possSet[name]; ok {
						possSet[name] = nil
					} else {
						fmt.Printf("addPoss %v\n", name)
						f.path = append(append([]string(nil), h.path...), f.path...)
						possSet[name] = &f
					}
				}
			}

			for len(work) > 0 {
				h, work = work[0], work[1:]
				debugf(" - examining %v\n", h.es)

				// what do we have?
				if typeIsInvalid(h.typ) {
					if seen[h.es] {
						continue
					}
					seen[h.es] = true
					debugf("using es check\n")
					it, ok := o.immTmpls[h.es]
					if !ok {
						panic(fmt.Errorf("failed to find generated imm type for %v", h.es))
					}

					switch it := it.(type) {
					case *immStruct:
						// here the fields do _not_ have a prefix

						for _, f := range it.fields {
							if h.typ == nil {
								// we are at the first level of a struct
								// so the paths must be the prefixed field names
								fname := fieldNamePrefix + fieldHidingPrefix + f.name

								if f.anon {
									fname = fieldAnonPrefix + fname
								}
								addPoss(f.name, field{
									path: []string{fname},
									typ:  o.exprString(f.field.Type),
									doc:  f.field.Doc,
								})
							} else {
								addPoss(f.name, field{
									path: []string{f.name + "()"},
									typ:  o.exprString(f.field.Type),
								})
							}

							if f.anon {
								next = append(next, embedded{
									es:   o.exprString(f.field.Type),
									path: append(append([]string(nil), h.path...), fieldTypeToIdent(f.field.Type).Name+"()"),
									typ:  o.info.TypeOf(f.field.Type),
								})
							}

							debugf(")) %v %v %v\n", f.name, f.field.Type, o.exprString(f.field.Type))
						}
					}
				} else {
					type ptr struct {
						types.Type
					}
					kt := h.typ
					if pt, ok := kt.(*types.Pointer); ok {
						kt = ptr{pt.Elem()}
					}
					if seen[kt] {
						continue
					}
					seen[kt] = true
					debugf("using type check on %T %v\n", h.typ, h.typ)
					if v, ok := util.IsImmType(h.typ).(util.ImmTypeStruct); ok {
						is := v.Struct
						fmt.Printf("))) %v %v\n", h.typ, is.NumFields())
						for i := 0; i < is.NumFields(); i++ {
							f := is.Field(i)
							name := f.Name()
							isAnon := false
							if strings.HasPrefix(name, "anon") {
								isAnon = true
								name = strings.TrimPrefix(name, "anon")
							}
							if !strings.HasPrefix(name, "field_") {
								continue
							}
							name = strings.TrimPrefix(name, "field_")
							// we can only consider exported fields
							if r, _ := utf8.DecodeRuneInString(name); unicode.IsLower(r) {
								continue
							}
							addPoss(name, field{
								typ:  possInvalidTypeToString(f.Type()),
								path: []string{name + "()"},
							})

							if isAnon {
								next = append(next, embedded{
									path: append(append([]string(nil), h.path...), name+"()"),
									typ:  f.Type(),
								})
							}
						}
					} else if v, ok := h.typ.Underlying().(*types.Struct); ok {
						fmt.Printf("))) %v %v\n", h.typ, v.NumFields())
						for i := 0; i < v.NumFields(); i++ {
							f := v.Field(i)
							fmt.Printf("::: %v %v %v\n", f.Name(), f.Exported(), f.Anonymous())
							if !f.Exported() {
								continue
							}
							name := f.Name()
							addPoss(name, field{
								typ:  possInvalidTypeToString(f.Type()),
								path: []string{name},
							})
							if f.Anonymous() {
								next = append(next, embedded{
									path: append(append([]string(nil), h.path...), name),
									typ:  f.Type(),
								})
							}
						}
					}
				}

				if len(work) == 0 {
					for n, f := range possSet {
						if f == nil {
							continue
						}
						set[n] = f
					}
					possSet = make(map[string]*field)
					work = next
					next = nil
				}
			}

			is.methods = set

			debugf("-----------\n")
			for n, f := range set {
				if f.path == nil {
					continue
				}

				debugf("%v() %v => %v\n", n, f.typ, strings.Join(f.path, "."))
			}
			debugf("===============\n")
		}
	}
}
