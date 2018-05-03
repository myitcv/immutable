package main

import (
	"fmt"
	"go/types"
	"strings"
)

func (o *output) calcMethodSets() {
	for _, fts := range o.files {
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
									path:   []string{fname},
									typ:    f.field.Type,
									setter: true,
									doc:    f.field.Doc,
								})
							} else {
								addPoss(f.name, field{
									path: []string{f.name + "()"},
									typ:  f.field.Type,
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
