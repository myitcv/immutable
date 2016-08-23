package generator

const ImmSliceTmpl = `
type {{.Name}} struct {
	theSlice []{{.Type}}

	mutable bool
}

func {{Export "New"}}{{Capitalise .Name}}() {{.Name}} {
	return {{.Name}}{}
}

func (m {{.Name}}) {{Choose "Len" "len_"}}() int {
	return len(m.theSlice)
}

func (m {{.Name}}) {{Export "Get"}}(i int) {{.Type}} {
	return m.theSlice[i]
}

func (m {{.Name}}) {{Export "AsMutable"}}() {{.Name}} {
	res := m.dup()
	res.mutable = true

	return res
}

func (m {{.Name}}) dup() {{.Name}} {
	resSlice := make([]{{.Type}}, len(m.theSlice))

	for i := range m.theSlice {
		resSlice[i] = m.theSlice[i]
	}

	res := {{.Name}}{
		theSlice: resSlice,
	}

	return res
}

func (m {{.Name}}) {{Export "AsImmutable"}}() {{.Name}} {
	m.mutable = false

	return m
}

func (m {{.Name}}) {{Choose "Range" "range_"}}() []{{.Type}} {
	return m.theSlice
}

func (m {{.Name}}) {{Export "WithMutations"}}(f func(mi {{.Name}})) {{.Name}} {
	res := m.{{Export "AsMutable"}}()
	f(res)
	res = res.{{Export "AsImmutable"}}()

	// TODO optimise here if the maps are identical?

	return res
}

func (m {{.Name}}) {{Export "Set"}}(i int, v {{.Type}}) {{.Name}} {
	if m.mutable {
		m.theSlice[i] = v
		return m
	}

	// TODO: work out a way of enabling this
	// if m.theSlice[i] == v {
	// 	return m
	// }

	res := m.dup()
	res.theSlice[i] = v

	return res
}
`
