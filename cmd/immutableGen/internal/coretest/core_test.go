package coretest_test

import (
	"testing"

	"github.com/myitcv/immutable/cmd/immutableGen/internal/coretest"
)

func TestBasic(t *testing.T) {
	m1 := coretest.NewMyMap()

	if m1.Len() != 0 {
		t.Fatalf("Expected m1 length to be 0, got %v", m1.Len())
	}

	m2 := m1.Set("test", 5)

	if m1 == m2 {
		t.Fatalf("Expected m1 and m2 to be different objects")
	}
	if m1.Len() != 0 {
		t.Fatalf("Expected m1 length to be 0, got %v", m1.Len())
	}
	if m2.Len() != 1 {
		t.Fatalf("Expected m2 length to be 1, got %v", m2.Len())
	}

	l1 := coretest.NewMySlice()

	if l1.Len() != 0 {
		t.Fatalf("Expected l1 length to be 0, got %v", l1.Len())
	}

	s1 := "test"

	l2 := l1.Append(&s1)

	if l1 == l2 {
		t.Fatalf("Expected l1 and l2 to be different objects")
	}
	if l1.Len() != 0 {
		t.Fatalf("Expected l1 length to be 0, got %v", l1.Len())
	}
	if l2.Len() != 1 {
		t.Fatalf("Expected l2 length to be 1, got %v", l2.Len())
	}

	ms1 := new(coretest.MyStruct)

	if ms1.Name() != "" {
		t.Fatalf("Expected ms1.Name() to be \"\", got %v", ms1.Name())
	}

	ms2 := ms1.SetName("paul")

	if ms1 == ms2 {
		t.Fatalf("Expected ms1 and ms2 to be different objects")
	}
	if ms1.Name() != "" {
		t.Fatalf("Expected ms1.Name() to be \"\", got %v", ms1.Name())
	}
	if ms2.Name() != "paul" {
		t.Fatalf("Expected ms2.Name() to be \"paul\", got %v", ms2.Name())
	}
}

func TestAsMutable(t *testing.T) {
	m1 := coretest.NewMyMap().AsMutable()

	if m1.Len() != 0 {
		t.Fatalf("Expected m1 length to be 0, got %v", m1.Len())
	}

	m2 := m1.Set("test", 5)

	if m1 != m2 {
		t.Fatalf("Expected m1 and m2 to be the same object")
	}
	if m2.Len() != 1 {
		t.Fatalf("Expected m2 length to be 1, got %v", m2.Len())
	}

	l1 := coretest.NewMySlice().AsMutable()

	if l1.Len() != 0 {
		t.Fatalf("Expected l1 length to be 0, got %v", l1.Len())
	}

	s1 := "test"

	l2 := l1.Append(&s1)

	if l1 != l2 {
		t.Fatalf("Expected l1 and l2 to be the same object")
	}
	if l2.Len() != 1 {
		t.Fatalf("Expected l2 length to be 1, got %v", l2.Len())
	}

	ms1 := new(coretest.MyStruct).AsMutable()

	if ms1.Name() != "" {
		t.Fatalf("Expected ms1.Name() to be \"\", got %v", ms1.Name())
	}

	ms2 := ms1.SetName("paul")

	if ms1 != ms2 {
		t.Fatalf("Expected ms1 and ms2 to be the same object")
	}
	if ms2.Name() != "paul" {
		t.Fatalf("Expected ms2.Name() to be \"paul\", got %v", ms2.Name())
	}
}

func TestTestTypes(t *testing.T) {
	m := coretest.NewMyTestMap()
	if m == nil {
		t.Fatalf("could not created instance of MyTestMap")
	}
}
