// Code generated by immutableGen. DO NOT EDIT.

// Copyright (c) 2016 Paul Jolly <paul@myitcv.org.uk>, all rights reserved.
// Use of this document is governed by a license found in the LICENSE document.

package example

//go:generate echo "hello world"
//immutableVet:skipFile

import (
	"myitcv.io/immutable"
)

// MyMap will be exported
//
// MyMap is an immutable type and has the following template:
//
// 	map[string]*MySlice
//
type MyMap struct {
	theMap  map[string]*MySlice
	mutable bool
	__tmpl  *_Imm_MyMap
}

var _ immutable.Immutable = new(MyMap)
var _ = new(MyMap).__tmpl

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

func (mr *MyMap) WithMutable(f func(m *MyMap)) *MyMap {
	res := mr.AsMutable()
	f(res)
	res = res.AsImmutable(mr)

	return res
}

func (mr *MyMap) WithImmutable(f func(m *MyMap)) *MyMap {
	prev := mr.mutable
	mr.mutable = false
	f(mr)
	mr.mutable = prev

	return mr
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
func (s *MyMap) IsDeeplyNonMutable(seen map[interface{}]bool) bool {
	if s == nil {
		return true
	}

	if s.Mutable() {
		return false
	}
	if s.Len() == 0 {
		return true
	}

	if seen == nil {
		return s.IsDeeplyNonMutable(make(map[interface{}]bool))
	}

	if seen[s] {
		return true
	}

	seen[s] = true

	for _, v := range s.theMap {
		if v != nil && !v.IsDeeplyNonMutable(seen) {
			return false
		}
	}
	return true
}

// MySlice will be exported
//
// MySlice is an immutable type and has the following template:
//
// 	[]*MyMap
//
type MySlice struct {
	theSlice []*MyMap
	mutable  bool
	__tmpl   *_Imm_MySlice
}

var _ immutable.Immutable = new(MySlice)
var _ = new(MySlice).__tmpl

func NewMySlice(s ...*MyMap) *MySlice {
	c := make([]*MyMap, len(s))
	copy(c, s)

	return &MySlice{
		theSlice: c,
	}
}

func NewMySliceLen(l int) *MySlice {
	c := make([]*MyMap, l)

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

func (m *MySlice) Get(i int) *MyMap {
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
	resSlice := make([]*MyMap, len(m.theSlice))

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

func (m *MySlice) Range() []*MyMap {
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

func (m *MySlice) Set(i int, v *MyMap) *MySlice {
	if m.mutable {
		m.theSlice[i] = v
		return m
	}

	res := m.dup()
	res.theSlice[i] = v

	return res
}

func (m *MySlice) Append(v ...*MyMap) *MySlice {
	if m.mutable {
		m.theSlice = append(m.theSlice, v...)
		return m
	}

	res := m.dup()
	res.theSlice = append(res.theSlice, v...)

	return res
}
func (s *MySlice) IsDeeplyNonMutable(seen map[interface{}]bool) bool {
	if s == nil {
		return true
	}

	if s.Mutable() {
		return false
	}
	if s.Len() == 0 {
		return true
	}

	if seen == nil {
		return s.IsDeeplyNonMutable(make(map[interface{}]bool))
	}

	if seen[s] {
		return true
	}

	seen[s] = true

	for _, v := range s.theSlice {
		if v != nil && !v.IsDeeplyNonMutable(seen) {
			return false
		}
	}
	return true
}

// MyStruct will be exported.
//
// It is a special type.
//
// MyStruct is an immutable type and has the following template:
//
// 	struct {
// 		Name	string
//
// 		surname	string
//
// 		self	*MyStruct
//
// 		age	int
// 	}
//
type MyStruct struct {
	field_Name    string `tag:"value"`
	field_surname string
	field_self    *MyStruct
	field_age     int `tag:"age"`

	mutable bool
	__tmpl  *_Imm_MyStruct
}

var _ immutable.Immutable = new(MyStruct)
var _ = new(MyStruct).__tmpl

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

func (s *MyStruct) IsDeeplyNonMutable(seen map[interface{}]bool) bool {
	if s == nil {
		return true
	}

	if s.Mutable() {
		return false
	}

	if seen == nil {
		return s.IsDeeplyNonMutable(make(map[interface{}]bool))
	}

	if seen[s] {
		return true
	}

	seen[s] = true
	{
		v := s.field_self

		if v != nil && !v.IsDeeplyNonMutable(seen) {
			return false
		}
	}
	return true
}

// Name is a field in MyStruct
func (s *MyStruct) Name() string {
	return s.field_Name
}

// SetName is the setter for Name()
func (s *MyStruct) SetName(n string) *MyStruct {
	if s.mutable {
		s.field_Name = n
		return s
	}

	res := *s
	res.field_Name = n
	return &res
}

// age will not be exported
func (s *MyStruct) age() int {
	return s.field_age
}

// setAge is the setter for Age()
func (s *MyStruct) setAge(n int) *MyStruct {
	if s.mutable {
		s.field_age = n
		return s
	}

	res := *s
	res.field_age = n
	return &res
}
func (s *MyStruct) self() *MyStruct {
	return s.field_self
}

// setSelf is the setter for Self()
func (s *MyStruct) setSelf(n *MyStruct) *MyStruct {
	if s.mutable {
		s.field_self = n
		return s
	}

	res := *s
	res.field_self = n
	return &res
}

// surname will not be exported
func (s *MyStruct) surname() string {
	return s.field_surname
}

// setSurname is the setter for Surname()
func (s *MyStruct) setSurname(n string) *MyStruct {
	if s.mutable {
		s.field_surname = n
		return s
	}

	res := *s
	res.field_surname = n
	return &res
}
