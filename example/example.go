// Copyright (c) 2016 Paul Jolly <paul@myitcv.org.uk>, all rights reserved.
// Use of this document is governed by a license found in the LICENSE document.

package example

import "io"

// The following directive will result in a generated file that includes the
// directive //go:generate echo "hello world"
//
//go:generate immutableGen -licenseFile license_header.txt -G "echo \"hello world\""

// MyMap will be an immutable map
//
type _Imm_MyMap map[string]*MySlice

// MySlice will be an immutable slice
//
type _Imm_MySlice []string

// MyStruct will be an immutable struct
//
type _Imm_MyStruct struct {

	// Name will be exported
	//
	Name string `tag:"value"`

	// surname will not be exported
	//
	surname string

	// age will not be exported
	//
	age int `tag:"age"`
}

// *****************
// Example 1
//
// Immutable structs don't have constructors automatically generated; you have to create them.
// The zero value of an immutable struct however is usable
//
func NewMyStructExample1_1() *MyStruct {

	//
	// don't do this; pointless copy
	//
	res := new(MyStruct).AsMutable()
	return res.AsImmutable(nil)
}

func NewMyStructExample1_2() *MyStruct {

	//
	// instead just use the zero value (which actually makes this constructor pointless)
	//
	return new(MyStruct)
}

// *****************
// Example 2
//
// Often it's cleaner to use WithMutable instead of topping and tailing with AsImmutable() and
// AsImmutable(...)
//
func NewMyStructExample2_1(n string) *MyStruct {

	//
	// Compare
	//
	res := new(MyStruct).AsMutable()

	res.SetName(n)
	res.setAge(42)

	return res.AsImmutable(nil)
}

func NewMyStructExample2_2(n string) *MyStruct {

	//
	// vs
	//
	return new(MyStruct).WithMutable(func(b *MyStruct) {
		b.SetName(n)
		b.setAge(42)
	})
}

// *****************
// Example 3
//
// Not everything has to be a method; sometimes functions which return instances of immutable
// types are much cleaner
//
// Consider
//
func (m *MyStruct) Parse(r io.Reader) *MyStruct {

	res := m.AsMutable()
	defer res.AsImmutable(m)

	// ...

	return res
}

//
// vs
//

func NewMyStructReader(r io.Reader) *MyStruct {
	res := new(MyStruct).AsMutable()
	defer res.AsImmutable(nil)

	// ...

	return res
}

//
// in particular how you have to use the two
//

func getStruct() {
	var r io.Reader

	_ = new(MyStruct).WithMutable(func(s *MyStruct) {
		s.Parse(r)
	})

	//
	// vs
	//

	_ = NewMyStructReader(r)
}

// *****************
// Example 4
//
// Think about what the most common cases for creating instances of immutable structs will be. Sometimes it's useful
// to pass in everything as arguments, at other times simply accept an initialiser function
//
func NewMyStructWithArgs(name string, age int) *MyStruct {

	return new(MyStruct).WithMutable(func(b *MyStruct) {
		b.SetName(name)
		b.setAge(age)
	})
}

func NewMyStructInitialiser(inits ...func(b *MyStruct)) *MyStruct {

	return new(MyStruct).WithMutable(func(b *MyStruct) {
		b.setAge(42)

		for _, i := range inits {
			i(b)
		}
	})
}

// *****************
// Example 5
//
// A method/function should leave the immutable state of any value/receiver in the
// state it received it; AsImmutable has been adapted to help make this easier. WithMutable
// also does the right thing. Defer statements should follow immediately after the
// AsMutable if that approach is followed
//
func (t *MyStruct) Update_Bad1(age int) *MyStruct {

	//
	// Bad; received as either mutable/immutable, return mutable
	//

	resBad := t.AsMutable()
	return resBad
}

func (t *MyStruct) Update_Good1(age int) *MyStruct {

	//
	// Good eg 1
	//

	resGood1 := t.AsMutable()
	defer resGood1.AsImmutable(t)

	return resGood1
}

func (t *MyStruct) Update_Good2(age int) *MyStruct {

	//
	// Good eg 2
	//

	return t.WithMutable(func(t *MyStruct) {
		// ...
	})
}

// *****************
// Example 6
//
// When creating instances of immutable types in the course of a method function,
// those values should always be returned in an immutable state (because this is
// the default for zero values of any immutable type). Constructors are the best example
// of this
//
func NewMyStructExample6_1(name string, age int) *MyStruct {

	//
	// Bad example; creating a new MyStruct and returning it mutable
	//
	res := new(MyStruct).AsMutable()
	res.SetName(name)
	return res
}

func NewMyStructExample6_2(name string, age int) *MyStruct {

	//
	// Good example; returning a new value that is immutable
	//
	return new(MyStruct).WithMutable(func(t *MyStruct) {
		t.setSurname("Jones")
		t.SetName(name)
		t.setAge(age)
	})
}

// *****************
// Example 7
//
// When creating new instances of immutable types and using AsMutable, then the
// call to AsImmutable takes nil
//
func NewMyStructExample7(name string, age int) *MyStruct {

	res := new(MyStruct).AsMutable()
	defer res.AsImmutable(nil)

	res.SetName(name)

	return res
}

// *****************
// Example 8
//
// The zero value of an immutable slice can be appended to. This is the immutable
// equivalent of being able to append on a nil regular slice
//
func Example8() *MySlice {

	return new(MySlice).WithMutable(func(s *MySlice) {
		s.Append("test")
	})
}

// *****************
// Example 9
//
// The zero value of an immutable map is not usable; instead you need to create an instance
// via one of the generated constructors
//
func Example9() *MyMap {

	return NewMyMap(func(m *MyMap) {
		m.Set("eg 1", nil)
	})
}
