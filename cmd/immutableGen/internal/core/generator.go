package core

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	gen "github.com/myitcv/immutable/gen"
)

const (
	fieldHidingPrefix = "_"
)

func Execute(file string, pkg string, licenseHeader string) error {

	path := filepath.Dir(file)
	base := filepath.Base(file)
	basename := strings.TrimSuffix(base, ".go")

	license := commentString(licenseHeader)

	g := &generator{
		path:    path,
		envFile: file,
		envPkg:  pkg,

		envBase: basename,

		license: license,

		fset: token.NewFileSet(),

		output: bytes.NewBuffer(nil),

		imports: make(map[*ast.ImportSpec]struct{}),
	}

	f, err := parser.ParseFile(g.fset, file, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	g.file = f

	g.commentMap = ast.NewCommentMap(g.fset, f, f.Comments)

	return g.gen()
}

type generator struct {
	path string

	envFile string
	envPkg  string

	output *bytes.Buffer

	// the envFile without its .go suffix
	envBase string

	license string

	fset *token.FileSet
	file *ast.File

	commentMap ast.CommentMap

	imports map[*ast.ImportSpec]struct{}

	immMaps    []immMap
	immSlices  []immSlice
	immStructs []immStruct
}

func (g *generator) addImports(exp ast.Expr) {
	finder := &importFinder{
		imports: g.file.Imports,
		matches: make(map[*ast.ImportSpec]struct{}),
	}

	ast.Walk(finder, exp)

	for i, v := range finder.matches {
		g.imports[i] = v
	}
}

type immMap struct {
	name   string
	dec    *ast.GenDecl
	typ    ast.Expr
	keyTyp ast.Expr
	valTyp ast.Expr
}

type immSlice struct {
	name string
	typ  ast.Expr
	dec  *ast.GenDecl
}

type immStruct struct {
	name string
	dec  *ast.GenDecl
	st   *ast.StructType
}

func (g *generator) gen() error {
	// 1. parse the envFile
	// 2. gather the maps, slices and structs we need to make immutable
	// 3. calculate from 2 the imports required
	// 4. generate gen_$(basename $GOFILE .go)_immutable.go file

	err := g.gatherImmTypes()
	if err != nil {
		return err
	}

	err = g.genImmTypes()
	if err != nil {
		return err
	}

	return nil
}

func (g *generator) gatherImmTypes() error {

	for _, d := range g.file.Decls {

		gd, ok := d.(*ast.GenDecl)
		if !ok || gd.Tok != token.TYPE {
			continue
		}

		if len(gd.Specs) != 1 {
			panic("myitcv needs to better understand go/ast")
		}

		ts := gd.Specs[0].(*ast.TypeSpec)

		typName := ts.Name.Name

		if !strings.HasPrefix(typName, gen.ImmTypeIdPrefix) {
			continue
		}

		name := strings.TrimPrefix(typName, gen.ImmTypeIdPrefix)

		switch typ := ts.Type.(type) {
		case *ast.MapType:
			g.immMaps = append(g.immMaps, immMap{
				name:   name,
				dec:    gd,
				typ:    typ,
				keyTyp: typ.Key,
				valTyp: typ.Value,
			})

			g.addImports(ts.Type)

		case *ast.ArrayType:
			if typ.Len == nil {
				g.immSlices = append(g.immSlices, immSlice{
					name: name,
					dec:  gd,
					typ:  typ.Elt,
				})
			}

			g.addImports(ts.Type)

		case *ast.StructType:
			g.immStructs = append(g.immStructs, immStruct{
				name: name,
				dec:  gd,
				st:   typ,
			})

			g.addImports(ts.Type)

		}
	}

	return nil
}

func (g *generator) genImmTypes() error {

	if len(g.immStructs) == 0 && len(g.immSlices) == 0 && len(g.immMaps) == 0 {
		return nil
	}

	g.pf(g.license)

	g.pf("package %v\n", g.envPkg)

	if len(g.imports) > 0 {
		g.pln("import (")
		for i := range g.imports {
			if i.Name != nil {
				g.pfln("%v %v", i.Name.Name, i.Path.Value)
			} else {
				g.pfln("%v", i.Path.Value)
			}
		}
		g.pln(")")
	}

	err := g.genImmMaps()
	if err != nil {
		return err
	}
	err = g.genImmSlices()
	if err != nil {
		return err
	}
	err = g.genImmStructs()
	if err != nil {
		return err
	}

	source := g.output.Bytes()

	toWrite := source

	formatted, err := format.Source(source)
	if err == nil {
		toWrite = formatted
	} else {
		fmt.Printf("Failed to format: %v\n", err)
	}

	ofName := filepath.Join(g.path, gen.GenFilePrefix+g.envBase+gen.GenFileSuffix)
	of, err := os.Create(ofName)
	if err != nil {
		return err
	}

	_, err = of.Write(toWrite)
	if err != nil {
		return err
	}

	return nil
}

func (g *generator) genImmMaps() error {

	for _, m := range g.immMaps {
		blanks := struct {
			Name    string
			KeyType string
			ValType string
		}{
			Name:    m.name,
			KeyType: g.exprString(m.keyTyp),
			ValType: g.exprString(m.valTyp),
		}

		fm := exporter(m.name)

		tmpl := template.New("immmap")
		tmpl.Funcs(fm)
		_, err := tmpl.Parse(immMapTmpl)
		if err != nil {
			return err
		}

		err = tmpl.Execute(g.output, blanks)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *generator) genImmSlices() error {

	for _, s := range g.immSlices {
		blanks := struct {
			Name string
			Type string
		}{
			Name: s.name,
			Type: g.exprString(s.typ),
		}

		fm := exporter(s.name)

		tmpl := template.New("immslice")
		tmpl.Funcs(fm)
		_, err := tmpl.Parse(immSliceTmpl)
		if err != nil {
			return err
		}

		err = tmpl.Execute(g.output, blanks)
		if err != nil {
			return err
		}
	}

	return nil
}

type genField struct {
	Name       string
	Type       string
	DocComment string
}

func (g *generator) commentTextFor(n ast.Node) string {
	res := ""
	comms := g.commentMap[n]
	for _, cg := range comms {
		for _, c := range cg.List {
			res = res + c.Text + "\n"
		}
	}

	return res
}

func (g *generator) genImmStructs() error {
	for _, s := range g.immStructs {

		comms := g.commentTextFor(s.dec)

		g.pf(comms)

		g.pfln("type %v struct {", s.name)

		var fields []genField

		for _, f := range s.st.Fields.List {
			names := ""
			sep := ""

			typ := g.exprString(f.Type)

			tag := ""

			if f.Tag != nil {
				tag = f.Tag.Value
			}

			for _, n := range f.Names {
				names = names + sep + fieldHidingPrefix + n.Name
				sep = ", "

				fields = append(fields, genField{
					Name:       n.Name,
					Type:       typ,
					DocComment: g.commentTextFor(f),
				})

			}
			g.pfln("%v %v %v", names, typ, tag)
		}

		g.pln("")
		g.pln("mutable bool")

		g.pfln("}")

		exp := exporter(s.name)

		g.pt(`
		func {{Export "New"}}{{Capitalise .}}() *{{.}} {
			return &{{.}}{}
		}

		func (s *{{.}}) AsMutable() *{{.}} {
			res := *s
			res.mutable = true
			return &res
		}

		func (s *{{.}}) AsImmutable() *{{.}} {
			s.mutable = false
			return s
		}

		func (s *{{.}}) WithMutations(f func(si *{{.}})) *{{.}} {
			res := s.AsMutable()
			f(res)
			res = res.AsImmutable()
			if *res == *s {
				return s
			}

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

			g.pt(`
			{{.Field.DocComment -}}
			func (s *{{.TypeName}}) {{.Field.Name}}() {{.Field.Type}} {
				return s.`+fieldHidingPrefix+`{{.Field.Name}}
			}

			func (s *{{.TypeName}}) {{Export "Set"}}{{Capitalise .Field.Name}}(n {{.Field.Type}}) *{{.TypeName}} {
				// TODO: see if we can make this work
				// if n == s.{{.Field.Name}} {
				// 	return s
				// }

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
	return nil
}

func (g *generator) exprString(e ast.Expr) string {
	var buf bytes.Buffer

	err := printer.Fprint(&buf, g.fset, e)
	if err != nil {
		panic(err)
	}

	return buf.String()
}

func (g *generator) pln(i ...interface{}) {
	fmt.Fprintln(g.output, i...)
}

func (g *generator) pf(format string, i ...interface{}) {
	fmt.Fprintf(g.output, format, i...)
}

func (g *generator) pfln(format string, i ...interface{}) {
	g.pf(format+"\n", i...)
}

func (g *generator) pt(tmpl string, fm template.FuncMap, val interface{}) {

	t := template.New("tmp")
	t.Funcs(fm)

	_, err := t.Parse(tmpl)
	if err != nil {
		panic(err)
	}

	err = t.Execute(g.output, val)
	if err != nil {
		panic(err)
	}
}

func commentString(r string) string {
	res := ""

	buf := bytes.NewBuffer([]byte(r))

	lastLineEmpty := false
	scanner := bufio.NewScanner(buf)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			lastLineEmpty = true
		}
		res = res + fmt.Sprintln("//", line)
	}

	if err := scanner.Err(); err != nil {
		// this really would be exceptional... because we passed in a string
		panic(err)
	}

	// ensure we have a space before package
	if !lastLineEmpty {
		res = res + "\n"
	}

	return res
}
