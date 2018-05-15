// Copyright (c) 2016 Paul Jolly <paul@myitcv.org.uk>, all rights reserved.
// Use of this document is governed by a license found in the LICENSE document.

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/printer"
	"go/token"
	"go/types"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"myitcv.io/gogenerate"
	"myitcv.io/hybridimporter"
	"myitcv.io/immutable/util"
)

const (
	fieldHidingPrefix = "_"
)

func execute(dir string, envPkg string, licenseHeader string, cmds gogenCmds) {

	absDir, err := filepath.Abs(dir)
	if err != nil {
		fatalf("could not make absolute path from %v: %v", dir, err)
	}

	bpkg, err := build.ImportDir(absDir, 0)
	if err != nil {
		fatalf("could not resolve package from dir %v: %v", dir, err)
	}

	fset := token.NewFileSet()

	pkgs, err := parser.ParseDir(fset, dir, nil, parser.AllErrors|parser.ParseComments)
	if err != nil {
		fatalf("could not parse dir %v: %v", dir, err)
	}

	pkg, ok := pkgs[envPkg]

	if !ok {
		pps := make([]string, 0, len(pkgs))
		for k := range pkgs {
			pps = append(pps, k)
		}
		fatalf("expected to have parsed %v, instead parsed %v", envPkg, pps)
	}

	var files []*ast.File
	var fns []string
	for fn, f := range pkg.Files {
		files = append(files, f)
		fns = append(fns, fn)
	}

	sort.Strings(fns)

	imp, err := hybridimporter.New(&build.Default, fset, bpkg.ImportPath, ".")
	if err != nil {
		fatalf("failed to create importer for %v: %v", bpkg.ImportPath, err)
	}

	info := &types.Info{
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
		Types: make(map[ast.Expr]types.TypeAndValue),
	}

	conf := types.Config{
		IgnoreFuncBodies: true,
		Importer:         imp,
		Error:            func(err error) {},
	}

	_, err = conf.Check(bpkg.ImportPath, fset, files, info)
	if err != nil {
		if _, ok := err.(types.Error); !ok {
			fatalf("failed to type check %v: %v", bpkg.ImportPath, err)
		}
	}

	out := &output{
		dir:       dir,
		info:      info,
		fset:      fset,
		pkgName:   envPkg,
		pkgPath:   bpkg.ImportPath,
		license:   licenseHeader,
		goGenCmds: cmds,
		files:     make(map[*ast.File]*fileTmpls),
		cms:       make(map[*ast.File]ast.CommentMap),
		immTypes:  make(map[string]util.ImmType),
	}

	for _, fn := range fns {

		// skip files that we generated
		if gogenerate.FileGeneratedBy(fn, immutableGenCmd) {
			continue
		}

		f := pkg.Files[fn]
		out.curFile = f

		out.cms[f] = ast.NewCommentMap(fset, f, f.Comments)
		out.gatherImmTypes()

	}

	out.genImmTypes()
}

type output struct {
	dir       string
	pkgName   string
	pkgPath   string
	fset      *token.FileSet
	license   string
	goGenCmds gogenCmds
	info      *types.Info

	output *bytes.Buffer

	curFile *ast.File

	// a convenience map of all the imm types we will
	// be generating in this package
	immTypes map[string]util.ImmType

	files map[*ast.File]*fileTmpls
	cms   map[*ast.File]ast.CommentMap
}

type fileTmpls struct {
	imports map[*ast.ImportSpec]struct{}

	maps    []immMap
	slices  []immSlice
	structs []immStruct
}

func (o *output) isImm(t types.Type, exp string) util.ImmType {
	ct := t
	switch v := ct.(type) {
	case *types.Pointer:
		ct = v.Elem()
	case *types.Named:
		ct = v.Underlying()
	}

	// we might have an invalid type because it refers to a yet-to-be-generated
	// immutable type in this package. If that is the case we fall back to a
	// comparison of the string representation of the type (which will be a
	// pointer).
	if tb, ok := ct.(*types.Basic); ok && tb.Kind() == types.Invalid {
		return o.immTypes[exp]
	}

	return util.IsImmType(t)
}

func (o *output) gatherImmTypes() {
	file := o.curFile
	fset := o.fset
	pkgPath := o.pkgPath

	g := &fileTmpls{
		imports: make(map[*ast.ImportSpec]struct{}),
	}

	impf := &importFinder{
		imports: file.Imports,
		matches: g.imports,
	}

	for _, d := range file.Decls {

		gd, ok := d.(*ast.GenDecl)
		if !ok || gd.Tok != token.TYPE {
			continue
		}

		for _, s := range gd.Specs {
			ts := s.(*ast.TypeSpec)

			name, ok := util.IsImmTmpl(ts)
			if !ok {
				continue
			}

			typ := o.info.Defs[ts.Name].Type().(*types.Named)

			infof("found immutable declaration at %v: %v", fset.Position(gd.Pos()), typ)

			comm := commonImm{
				fset: fset,
				file: file,
				pkg:  pkgPath,
				dec:  gd,
			}

			switch u := typ.Underlying().(type) {
			case *types.Map:
				g.maps = append(g.maps, immMap{
					commonImm: comm,
					name:      name,
					typ:       u,
					syn:       ts.Type.(*ast.MapType),
				})
				o.immTypes["*"+name] = util.ImmTypeMap{}

				ast.Walk(impf, ts.Type)

			case *types.Slice:
				// TODO support for arrays

				g.slices = append(g.slices, immSlice{
					commonImm: comm,
					name:      name,
					typ:       u,
					syn:       ts.Type.(*ast.ArrayType),
				})
				o.immTypes["*"+name] = util.ImmTypeSlice{}

				ast.Walk(impf, ts.Type)

			case *types.Struct:
				g.structs = append(g.structs, immStruct{
					commonImm: comm,
					name:      name,
					typ:       u,
					syn:       ts.Type.(*ast.StructType),
					special:   isSpecialStruct(name, u),
				})
				o.immTypes["*"+name] = util.ImmTypeStruct{}

				ast.Walk(impf, ts.Type)
			}

		}
	}

	o.files[o.curFile] = g
}

func isSpecialStruct(name string, st *types.Struct) bool {
	// work out whether this is a special struct with a Key field
	// pattern is:
	//
	// 1. struct field has a field called Key of type {{.StructName}}Key (non pointer)
	//
	// later checks will include:
	//
	// 2. said type has two fields, Uuid and Version, of type {{.StructName}}Uuid and uint64 respectively
	// 3. the underlying type of {{.StructName}}Uuid is uint64 (we might be able to relax these two
	// two underlying type restrictions)

	if st.NumFields() == 0 {
		return false
	}

	for i := 0; i < st.NumFields(); i++ {
		f := st.Field(i)

		if f.Name() != "Key" {
			continue
		}

		kst, ok := f.Type().Underlying().(*types.Struct)
		if !ok {
			continue
		}

		if kst.NumFields() != 2 {
			continue
		}

		uuid := kst.Field(0)
		if uuid.Name() != "Uuid" {
			continue
		}

		ver := kst.Field(1)
		if ver.Name() != "Version" {
			continue
		}

		// we found it
		return true
	}

	return false
}

func (o *output) genImmTypes() {
	for f, v := range o.files {
		o.curFile = f

		if len(v.maps) == 0 && len(v.slices) == 0 && len(v.structs) == 0 {
			continue
		}

		o.output = bytes.NewBuffer(nil)

		o.pfln("// Code generated by %v. DO NOT EDIT.", immutableGenCmd)
		o.pln("")

		o.pf(o.license)

		o.pf("package %v\n", o.pkgName)

		// is there a "standard" place for //go:generate comments?
		for _, v := range o.goGenCmds {
			o.pf("//go:generate %v\n", v)
		}

		o.pln("//immutableVet:skipFile")
		o.pln("")

		o.pln("import (")

		o.pln("\"myitcv.io/immutable\"")
		o.pln()

		for i := range v.imports {
			if i.Name != nil {
				o.pfln("%v %v", i.Name.Name, i.Path.Value)
			} else {
				o.pfln("%v", i.Path.Value)
			}
		}

		o.pln(")")

		o.pln("")

		o.genImmMaps(v.maps)
		o.genImmSlices(v.slices)
		o.genImmStructs(v.structs)

		source := o.output.Bytes()

		toWrite := source

		fn := o.fset.Position(f.Pos()).Filename

		// this is the file path
		offn, ok := gogenerate.NameFileFromFile(fn, immutableGenCmd)
		if !ok {
			fatalf("could not name file from %v", fn)
		}

		out := bytes.NewBuffer(nil)
		cmd := exec.Command("gofmt", "-s")
		cmd.Stdin = o.output
		cmd.Stdout = out

		err := cmd.Run()
		if err == nil {
			toWrite = out.Bytes()
		} else {
			infof("failed to format %v: %v", fn, err)
		}

		if err := ioutil.WriteFile(offn, toWrite, 0644); err != nil {
			fatalf("could not write %v: %v", offn, err)
		}
	}
}

func (o *output) exprString(e ast.Expr) string {
	var buf bytes.Buffer

	err := printer.Fprint(&buf, o.fset, e)
	if err != nil {
		panic(err)
	}

	return buf.String()
}

func (o *output) printCommentGroup(d *ast.CommentGroup) {
	if d != nil {
		for _, c := range d.List {
			o.pfln("%v", c.Text)
		}
	}
}

func (o *output) printImmPreamble(name string, node ast.Node) {
	fset := o.fset

	if st, ok := node.(*ast.StructType); ok {

		// we need to do some manipulation

		buf := bytes.NewBuffer(nil)

		fmt.Fprintf(buf, "struct {\n")

		if st.Fields != nil && st.Fields.NumFields() > 0 {
			line := o.fset.Position(st.Fields.List[0].Pos()).Line

			for _, f := range st.Fields.List {
				curLine := o.fset.Position(f.Pos()).Line

				if line != curLine {
					// catch up
					fmt.Fprintln(buf, "")
					line = curLine
				}

				ids := make([]string, 0, len(f.Names))
				for _, n := range f.Names {
					ids = append(ids, n.Name)
				}
				fmt.Fprintf(buf, "%v %v\n", strings.Join(ids, ","), o.exprString(f.Type))

				line++
			}
		}

		fmt.Fprintf(buf, "}")

		exprStr := buf.String()

		fset = token.NewFileSet()
		newnode, err := parser.ParseExprFrom(fset, "", exprStr, 0)
		if err != nil {
			fatalf("could not parse documentation struct from %v: %v", exprStr, err)
		}

		node = newnode
	}

	o.pln("//")
	o.pfln("// %v is an immutable type and has the following template:", name)
	o.pln("//")

	tmplBuf := bytes.NewBuffer(nil)

	err := printer.Fprint(tmplBuf, fset, node)
	if err != nil {
		fatalf("could not printer template declaration: %v", err)
	}

	sc := bufio.NewScanner(tmplBuf)
	for sc.Scan() {
		o.pfln("// \t%v", sc.Text())
	}
	if err := sc.Err(); err != nil {
		fatalf("could not scan printed template: %v", err)
	}

	o.pln("//")
}

func (o *output) pln(i ...interface{}) {
	fmt.Fprintln(o.output, i...)
}

func (o *output) pf(format string, i ...interface{}) {
	fmt.Fprintf(o.output, format, i...)
}

func (o *output) pfln(format string, i ...interface{}) {
	o.pf(format+"\n", i...)
}

func (o *output) pt(tmpl string, fm template.FuncMap, val interface{}) {

	// on the basis most templates are for convenience define inline
	// as raw string literals which start the ` on one line but then start
	// the template on the next (for readability) we strip the first leading
	// \n if one exists
	tmpl = strings.TrimPrefix(tmpl, "\n")

	t := template.New("tmp")
	t.Funcs(fm)

	_, err := t.Parse(tmpl)
	if err != nil {
		panic(err)
	}

	err = t.Execute(o.output, val)
	if err != nil {
		panic(err)
	}
}

func fatalf(format string, args ...interface{}) {
	panic(fmt.Errorf(format, args...))
}

func infoln(args ...interface{}) {
	if *fGoGenLog == string(gogenerate.LogInfo) {
		log.Println(args...)
	}
}

func infof(format string, args ...interface{}) {
	if *fGoGenLog == string(gogenerate.LogInfo) {
		log.Printf(format, args...)
	}
}
