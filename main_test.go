package main

import "testing"

func TestDifferentiate(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"x", "1"},
		{"5", "0"},
		{"x^3 + 2*x + 5", "3*x^2 + 2"},
		{"4*x^2 - 3*x + 7", "8*x - 3"},
		{"9*x^5", "45*x^4"},
		{"-x^2 + x - 1", "-2*x + 1"},
		{"x^2 - x^2 + 5", "0"},
		{"3*x^4 + 2*x^3 - x + 8", "12*x^3 + 6*x^2 - 1"},
	}

	for _, tc := range tests {
		p, err := parse(tc.input)
		if err != nil {
			t.Fatalf("parse(%q) failed: %v", tc.input, err)
		}

		got := format(simplify(differentiate(p)))
		if got != tc.want {
			t.Fatalf("input %q: got %q want %q", tc.input, got, tc.want)
		}
	}
}

