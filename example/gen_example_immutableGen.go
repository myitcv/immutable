// Copyright (c) 2016 Paul Jolly <paul@myitcv.org.uk>, all rights reserved.
// Use of this document is governed by a license found in the LICENSE document.

// File generated by immutableGen

package example

//go:generate echo "hello world"

import (
	"github.com/myitcv/immutable"
)

// MyMap will be an immutable map
//
//
// MyMap is an immutable type and has the following template:
//
// 	map[string]*MySlice
//
type MyMap struct {
	//github.com/myitcv/immutable:ImmutableType

	theMap  map[string]*MySlice
	mutable bool
	__tmpl  _Imm_MyMap
}

var _ immutable.Immutable = &MyMap{}

func NewMyMap(inits ...func(m *MyMap)) *MyMap {
	res := NewMyMapCap(0)
	if len(inits) == 0 {
		return res
	}

	return res.WithMutable(func(m *MyMap) {
		for _, i := range inits {
			i(m)
		}
	})
}

func NewMyMapCap(l int) *MyMap {
	return &MyMap{
		theMap: make(map[string]*MySlice, l),
	}
}

func (m *MyMap) Mutable() bool {
	return m.mutable
}

func (m *MyMap) Len() int {
	if m == nil {
		return 0
	}

	return len(m.theMap)
}

func (m *MyMap) Get(k string) (*MySlice, bool) {
	v, ok := m.theMap[k]
	return v, ok
}

func (m *MyMap) AsMutable() *MyMap {
	if m == nil {
		return nil
	}

	if m.Mutable() {
		return m
	}

	res := m.dup()
	res.mutable = true

	return res
}

func (m *MyMap) dup() *MyMap {
	resMap := make(map[string]*MySlice, len(m.theMap))

	for k := range m.theMap {
		resMap[k] = m.theMap[k]
	}

	res := &MyMap{
		theMap: resMap,
	}

	return res
}

func (m *MyMap) AsImmutable(v *MyMap) *MyMap {
	if m == nil {
		return nil
	}

	if v == m {
		return m
	}

	m.mutable = false
	return m
}

func (m *MyMap) Range() map[string]*MySlice {
	if m == nil {
		return nil
	}

	return m.theMap
}

func (m *MyMap) WithMutable(f func(mi *MyMap)) *MyMap {
	res := m.AsMutable()
	f(res)
	res = res.AsImmutable(m)

	return res
}

func (m *MyMap) WithImmutable(f func(mi *MyMap)) *MyMap {
	prev := m.mutable
	m.mutable = false
	f(m)
	m.mutable = prev

	return m
}

func (m *MyMap) Set(k string, v *MySlice) *MyMap {
	if m.mutable {
		m.theMap[k] = v
		return m
	}

	res := m.dup()
	res.theMap[k] = v

	return res
}

func (m *MyMap) Del(k string) *MyMap {
	if _, ok := m.theMap[k]; !ok {
		return m
	}

	if m.mutable {
		delete(m.theMap, k)
		return m
	}

	res := m.dup()
	delete(res.theMap, k)

	return res
}

// MySlice will be an immutable slice
//
//
// MySlice is an immutable type and has the following template:
//
// 	[]string
//
type MySlice struct {
	//github.com/myitcv/immutable:ImmutableType

	theSlice []string
	mutable  bool
	__tmpl   _Imm_MySlice
}

var _ immutable.Immutable = &MySlice{}

func NewMySlice(s ...string) *MySlice {
	c := make([]string, len(s))
	copy(c, s)

	return &MySlice{
		theSlice: c,
	}
}

func NewMySliceLen(l int) *MySlice {
	c := make([]string, l)

	return &MySlice{
		theSlice: c,
	}
}

func (m *MySlice) Mutable() bool {
	return m.mutable
}

func (m *MySlice) Len() int {
	if m == nil {
		return 0
	}

	return len(m.theSlice)
}

func (m *MySlice) Get(i int) string {
	return m.theSlice[i]
}

func (m *MySlice) AsMutable() *MySlice {
	if m == nil {
		return nil
	}

	if m.Mutable() {
		return m
	}

	res := m.dup()
	res.mutable = true

	return res
}

func (m *MySlice) dup() *MySlice {
	resSlice := make([]string, len(m.theSlice))

	for i := range m.theSlice {
		resSlice[i] = m.theSlice[i]
	}

	res := &MySlice{
		theSlice: resSlice,
	}

	return res
}

func (m *MySlice) AsImmutable(v *MySlice) *MySlice {
	if m == nil {
		return nil
	}

	if v == m {
		return m
	}

	m.mutable = false
	return m
}

func (m *MySlice) Range() []string {
	if m == nil {
		return nil
	}

	return m.theSlice
}

func (m *MySlice) WithMutable(f func(mi *MySlice)) *MySlice {
	res := m.AsMutable()
	f(res)
	res = res.AsImmutable(m)

	return res
}

func (m *MySlice) WithImmutable(f func(mi *MySlice)) *MySlice {
	prev := m.mutable
	m.mutable = false
	f(m)
	m.mutable = prev

	return m
}

func (m *MySlice) Set(i int, v string) *MySlice {
	if m.mutable {
		m.theSlice[i] = v
		return m
	}

	res := m.dup()
	res.theSlice[i] = v

	return res
}

func (m *MySlice) Append(v ...string) *MySlice {
	if m.mutable {
		m.theSlice = append(m.theSlice, v...)
		return m
	}

	res := m.dup()
	res.theSlice = append(res.theSlice, v...)

	return res
}

// MyStruct will be an immutable struct
//
//
// MyStruct is an immutable type and has the following template:
//
// 	struct {
// 		Name	string
//
// 		surname	string
//
// 		age	int
// 	}
//
type MyStruct struct {
	//github.com/myitcv/immutable:ImmutableType
	//

	_Name    string `tag:"value"`
	_surname string
	_age     int `tag:"age"`

	mutable bool
	__tmpl  _Imm_MyStruct
}

var _ immutable.Immutable = &MyStruct{}

func (s *MyStruct) AsMutable() *MyStruct {
	if s.Mutable() {
		return s
	}

	res := *s
	res.mutable = true
	return &res
}

func (s *MyStruct) AsImmutable(v *MyStruct) *MyStruct {
	if s == nil {
		return nil
	}

	if s == v {
		return s
	}

	s.mutable = false
	return s
}

func (s *MyStruct) Mutable() bool {
	return s.mutable
}

func (s *MyStruct) WithMutable(f func(si *MyStruct)) *MyStruct {
	res := s.AsMutable()
	f(res)
	res = res.AsImmutable(s)

	return res
}

func (s *MyStruct) WithImmutable(f func(si *MyStruct)) *MyStruct {
	prev := s.mutable
	s.mutable = false
	f(s)
	s.mutable = prev

	return s
}

// Name will be exported
//
func (s *MyStruct) Name() string {
	return s._Name
}

// SetName is the setter for Name()
func (s *MyStruct) SetName(n string) *MyStruct {
	if s.mutable {
		s._Name = n
		return s
	}

	res := *s
	res._Name = n
	return &res
}

// surname will not be exported
//
func (s *MyStruct) surname() string {
	return s._surname
}

// setSurname is the setter for Surname()
func (s *MyStruct) setSurname(n string) *MyStruct {
	if s.mutable {
		s._surname = n
		return s
	}

	res := *s
	res._surname = n
	return &res
}

// age will not be exported
//
func (s *MyStruct) age() int {
	return s._age
}

// setAge is the setter for Age()
func (s *MyStruct) setAge(n int) *MyStruct {
	if s.mutable {
		s._age = n
		return s
	}

	res := *s
	res._age = n
	return &res
}
