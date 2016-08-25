package core

const immMapTmpl = `
var _ immutable.Immutable = &{{.Name}}{}

type {{.Name}} struct {
	theMap map[{{.KeyType}}]{{.ValType}}

	mutable bool
}

func {{Export "New"}}{{Capitalise .Name}}() {{.Name}} {
	return {{.Name}}{
		theMap: make(map[{{.KeyType}}]{{.ValType}}),
	}
}

func (m {{.Name}})Mutable() bool {
	return m.mutable
}

func (m {{.Name}}) Len() int {
	return len(m.theMap)
}

func (m {{.Name}}) Get(k {{.KeyType}}) ({{.ValType}}, bool) {
	v, ok := m.theMap[k]
	return v, ok
}

func (m {{.Name}}) AsMutable() {{.Name}} {
	res := m.dup()
	res.mutable = true

	return res
}

func (m {{.Name}}) dup() {{.Name}} {
	resMap := make(map[{{.KeyType}}]{{.ValType}}, len(m.theMap))

	for k := range m.theMap {
		resMap[k] = m.theMap[k]
	}

	res := {{.Name}}{
		theMap: resMap,
	}

	return res
}

func (m {{.Name}}) AsImmutable() {{.Name}} {
	m.mutable = false

	return m
}

func (m {{.Name}}) Range() map[{{.KeyType}}]{{.ValType}} {
	return m.theMap
}

func (m {{.Name}}) WithMutations(f func(mi {{.Name}})) {{.Name}} {
	res := m.AsMutable()
	f(res)
	res = res.AsImmutable()

	// TODO optimise here if the maps are identical?

	return res
}

func (m {{.Name}}) Set(k {{.KeyType}}, v {{.ValType}}) {{.Name}} {
	if m.mutable {
		m.theMap[k] = v
		return m
	}

	// TODO: work out a way of enabling this
	// if n, ok := m.theMap[k]; ok && n == v {
	// 	return m
	// }

	res := m.dup()
	res.theMap[k] = v

	return res
}
`
