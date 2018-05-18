package coretest

import (
	"testing"

	"myitcv.io/immutable/cmd/immutableGen/internal/coretest/pkga"
)

func TestAnonFields(t *testing.T) {
	m := new(MyStruct)

	if v := m.string(); v != "" {
		t.Fatalf("expected zero value to be %v", "")
	}

	val := "test"
	m = m.setString(val)

	if v := m.string(); v == "" || v != val {
		t.Fatalf("expected set value to be %q", val)
	}
}

func TestEmbedAccess(t *testing.T) {
	a := new(pkga.PkgA).SetAddress("home")
	e2 := new(Embed2).SetAge(42)
	e1 := new(Embed1).WithMutable(func(e1 *Embed1) {
		e1.SetName("Paul")
		e1.SetEmbed2(e2)
		e1.SetPkgA(a)
	})

	{
		want := 42
		if got := e2.Age(); want != got {
			t.Fatalf("e2.Age(): want %v, got %v", want, got)
		}
	}
	{
		want := 42
		if got := e1.Age(); want != got {
			t.Fatalf("e1.Age(): want %v, got %v", want, got)
		}
	}
	{
		want := "home"
		if got := e1.Address(); want != got {
			t.Fatalf("e1.Address(): want %v, got %v", want, got)
		}
	}
}
