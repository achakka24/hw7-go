package main

import "strconv"

func format(e Expr) string {
	return formatWithPrec(e, 0)
}

func formatWithPrec(e Expr, parentPrec int) string {
	switch n := e.(type) {
	case Const:
		return strconv.Itoa(n.Value)
	case Var:
		return n.Name
	case Neg:
		s := formatWithPrec(n.X, precedence(n))
		if precedence(n.X) < precedence(n) {
			s = "(" + s + ")"
		}
		out := "-" + s
		if precedence(n) < parentPrec {
			return "(" + out + ")"
		}
		return out
	case Add:
		left := formatWithPrec(n.Left, precedence(n))
		right := formatWithPrec(n.Right, precedence(n))
		out := left + "+" + right
		if precedence(n) < parentPrec {
			return "(" + out + ")"
		}
		return out
	case Sub:
		left := formatWithPrec(n.Left, precedence(n))
		right := formatWithPrec(n.Right, precedence(n))
		if precedence(n.Right) <= precedence(n) {
			right = "(" + right + ")"
		}
		out := left + "-" + right
		if precedence(n) < parentPrec {
			return "(" + out + ")"
		}
		return out
	case Mul:
		left := formatWithPrec(n.Left, precedence(n))
		right := formatWithPrec(n.Right, precedence(n))
		if precedence(n.Left) < precedence(n) {
			left = "(" + left + ")"
		}
		if precedence(n.Right) < precedence(n) {
			right = "(" + right + ")"
		}
		out := left + "*" + right
		if precedence(n) < parentPrec {
			return "(" + out + ")"
		}
		return out
	case Div:
		left := formatWithPrec(n.Left, precedence(n))
		right := formatWithPrec(n.Right, precedence(n))
		if precedence(n.Left) < precedence(n) {
			left = "(" + left + ")"
		}
		if precedence(n.Right) <= precedence(n) {
			right = "(" + right + ")"
		}
		out := left + "/" + right
		if precedence(n) < parentPrec {
			return "(" + out + ")"
		}
		return out
	case Pow:
		base := formatWithPrec(n.Base, precedence(n))
		exp := formatWithPrec(n.Exponent, precedence(n))
		if precedence(n.Base) < precedence(n) {
			base = "(" + base + ")"
		}
		if precedence(n.Exponent) <= precedence(n) {
			exp = "(" + exp + ")"
		}
		out := base + "^" + exp
		if precedence(n) < parentPrec {
			return "(" + out + ")"
		}
		return out
	default:
		return ""
	}
}

func precedence(e Expr) int {
	switch e.(type) {
	case Add, Sub:
		return 1
	case Mul, Div:
		return 2
	case Pow:
		return 3
	case Neg:
		return 4
	case Const, Var:
		return 5
	default:
		return 0
	}
}
