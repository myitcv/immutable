package core

import (
	"text/template"
	"unicode"
	"unicode/utf8"
)

const (
	NonUtf8Initial = "Cannot decode utf8 character from the beginning of %v"
)

var ExportMap = template.FuncMap{
	"Export":     export,
	"Capitalise": capitalise,
	"Choose":     chooseFirst,
}
var UnexportMap = template.FuncMap{
	"Export":     unexport,
	"Capitalise": capitalise,
	"Choose":     chooseSecond,
}

// exporter returns a template.FuncMap with
// "Export" mapped to an appropriate function
// depending on the initial capitalisation of
// name
func exporter(name string) template.FuncMap {
	r, _ := utf8.DecodeRuneInString(name)

	// Note: we are choosing to ignore the situation where we decode utf8.RuneError
	// this situation would only happen if the source contained an invalid utf8 code point...
	// which is impossible because it won't compile

	// But we unexport for "safety", even though this doesn't mean anything

	if unicode.IsUpper(r) {
		return ExportMap
	}

	return UnexportMap
}

func export(name string) string {
	return capitalise(name)
}

func unexport(name string) string {
	return uncapitalise(name)
}

func capitalise(name string) string {
	r, n := utf8.DecodeRuneInString(name)

	// again, choosing to ignore error

	return string(unicode.ToUpper(r)) + name[n:]
}

func uncapitalise(name string) string {
	r, n := utf8.DecodeRuneInString(name)

	// again, choosing to ignore error

	return string(unicode.ToLower(r)) + name[n:]
}

func chooseFirst(s1, s2 string) string {
	return s1
}

func chooseSecond(s1, s2 string) string {
	return s2
}
