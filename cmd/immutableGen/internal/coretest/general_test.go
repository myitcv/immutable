package coretest

import "testing"

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
	e2 := new(Embed2).SetAge(42)
	e1 := new(Embed1).WithMutable(func(e1 *Embed1) {
		e1.SetName("Paul")
		e1.SetEmbed2(e2)
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
}
