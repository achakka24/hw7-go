package main

import "testing"

func TestDerivHomeworkCases(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"x^2", "2*x"},
		{"(x*2*x)/x", "2"},
		{"x^4+2*x^3-x^2+5*x-1/x", "4*x^3+6*x^2-2*x+5+1/x^2"},
		{"4*x^3+6*x^2-2*x+5+1/x^2", "12*x^2+12*x-2-2/x^3"},
		{"12*x^2+12*x-2-2/x^3", "24*x+12+6/x^4"},
	}

	for _, tc := range tests {
		expr, err := parse(tc.input)
		if err != nil {
			t.Fatalf("parse(%q) failed: %v", tc.input, err)
		}

		got := format(simplify(differentiate(expr)))
		if got != tc.want {
			t.Fatalf("input %q: got %q want %q", tc.input, got, tc.want)
		}
	}
}

func TestSimplifyHomeworkCases(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"5-x*(3/3)+2", "5-x+2"},
		{"1*x-0/3+2", "x+2"},
	}

	for _, tc := range tests {
		expr, err := parse(tc.input)
		if err != nil {
			t.Fatalf("parse(%q) failed: %v", tc.input, err)
		}
		got := format(simplify(expr))
		if got != tc.want {
			t.Fatalf("input %q: got %q want %q", tc.input, got, tc.want)
		}
	}
}
