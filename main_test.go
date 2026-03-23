package main

import "testing"

func TestDifferentiate(t *testing.T) {
	p, err := parse("x^3 + 2*x + 5")
	if err != nil {
		t.Fatal(err)
	}

	got := format(simplify(differentiate(p)))
	want := "3*x^2 + 2"

	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

