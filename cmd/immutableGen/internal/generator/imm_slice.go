package generator

const ImmSliceTmpl = `
type {{.Name}} struct {
	theSlice []{{.Type}}

	mutable bool
}

func {{Export "New"}}{{Capitalise .Name}}() {{.Name}} {
	return {{.Name}}{}
}

func (m {{.Name}}) Len() int {
	return len(m.theSlice)
}

func (m {{.Name}}) Get(i int) {{.Type}} {
	return m.theSlice[i]
}

func (m {{.Name}}) AsMutable() {{.Name}} {
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

func (m {{.Name}}) AsImmutable() {{.Name}} {
	m.mutable = false

	return m
}

func (m {{.Name}}) Range() []{{.Type}} {
	return m.theSlice
}

func (m {{.Name}}) WithMutations(f func(mi {{.Name}})) {{.Name}} {
	res := m.AsMutable()
	f(res)
	res = res.AsImmutable()

	// TODO optimise here if the maps are identical?

	return res
}

func (m {{.Name}}) Set(i int, v {{.Type}}) {{.Name}} {
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

func (m {{.Name}}) Append(v ...{{.Type}}) {{.Name}} {
	if m.mutable {
		m.theSlice = append(m.theSlice, v...)
		return m
	}

	// TODO: work out a way of enabling this
	// if m.theSlice[i] == v {
	// 	return m
	// }

	res := m.dup()
	res.theSlice = append(res.theSlice, v...)

	return res
}
`
