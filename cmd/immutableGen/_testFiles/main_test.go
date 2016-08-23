package main

import (
	"testing"
)

func TestBasic(t *testing.T) {
	m1 := newMyMap()

	if m1.Len() != 0 {
		t.Fatalf("Expected m1 length to be 0, got %v", m1.Len())
	}

	m2 := m1.Set("test", 5)

	if m1.Len() != 0 {
		t.Fatalf("Expected m1 length to be 0, got %v", m1.Len())
	}
	if m2.Len() != 1 {
		t.Fatalf("Expected m2 length to be 1, got %v", m2.Len())
	}

	l1 := NewSlice()

	if l1.Len() != 0 {
		t.Fatalf("Expected l1 length to be 0, got %v", l1.Len())
	}

	s1 := "test"

	l2 := l1.Append(&s1)

	if l1.Len() != 0 {
		t.Fatalf("Expected l1 length to be 0, got %v", l1.Len())
	}
	if l2.Len() != 1 {
		t.Fatalf("Expected l2 length to be 1, got %v", l2.Len())
	}

	ms1 := newMyStruct()

	if ms1.Name() != "" {
		t.Fatalf("Expected ms1.Name() to be \"\", got %v", ms1.Name())
	}

	ms2 := ms1.SetName("paul")

	if ms1.Name() != "" {
		t.Fatalf("Expected ms1.Name() to be \"\", got %v", ms1.Name())
	}
	if ms2.Name() != "paul" {
		t.Fatalf("Expected ms2.Name() to be \"paul\", got %v", ms2.Name())
	}
}
