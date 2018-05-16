// Code generated by immutableGen. DO NOT EDIT.

// Copyright (c) 2016 Paul Jolly <paul@myitcv.org.uk>, all rights reserved.
// Use of this document is governed by a license found in the LICENSE document.

package example

//go:generate echo "hello world"
//immutableVet:skipFile

import (
	"myitcv.io/immutable"
)

// via go generate, this template is code generated into the immutable Person
// struct within the same package
//
// Person is an immutable type and has the following template:
//
// 	struct {
// 		Name	string
// 		Age	int
// 	}
//
type Person struct {
	field_Name string
	field_Age  int

	mutable bool
	__tmpl  *_Imm_Person
}

var _ immutable.Immutable = new(Person)
var _ = new(Person).__tmpl

func (s *Person) AsMutable() *Person {
	if s.Mutable() {
		return s
	}

	res := *s
	res.mutable = true
	return &res
}

func (s *Person) AsImmutable(v *Person) *Person {
	if s == nil {
		return nil
	}

	if s == v {
		return s
	}

	s.mutable = false
	return s
}

func (s *Person) Mutable() bool {
	return s.mutable
}

func (s *Person) WithMutable(f func(si *Person)) *Person {
	res := s.AsMutable()
	f(res)
	res = res.AsImmutable(s)

	return res
}

func (s *Person) WithImmutable(f func(si *Person)) *Person {
	prev := s.mutable
	s.mutable = false
	f(s)
	s.mutable = prev

	return s
}

func (s *Person) IsDeeplyNonMutable(seen map[interface{}]bool) bool {
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
	return true
}
func (s *Person) Age() int {
	return s.field_Age
}

// SetAge is the setter for Age()
func (s *Person) SetAge(n int) *Person {
	if s.mutable {
		s.field_Age = n
		return s
	}

	res := *s
	res.field_Age = n
	return &res
}
func (s *Person) Name() string {
	return s.field_Name
}

// SetName is the setter for Name()
func (s *Person) SetName(n string) *Person {
	if s.mutable {
		s.field_Name = n
		return s
	}

	res := *s
	res.field_Name = n
	return &res
}
