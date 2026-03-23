package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Poly map[int]int

func parse(input string) (Poly, error) {
	s := strings.ReplaceAll(input, " ", "")
	s = strings.ReplaceAll(s, "-", "+-")
	if strings.HasPrefix(s, "+") {
		s = s[1:]
	}

	p := Poly{}
	for _, term := range strings.Split(s, "+") {
		if term == "" {
			continue
		}
		coeff, exp, err := parseTerm(term)
		if err != nil {
			return nil, err
		}
		p[exp] += coeff
	}
	return p, nil
}

func parseTerm(term string) (int, int, error) {
	if !strings.Contains(term, "x") {
		n, err := strconv.Atoi(term)
		return n, 0, err
	}

	coeff := 1
	exp := 1

	parts := strings.Split(term, "x")
	left := strings.TrimSuffix(parts[0], "*")
	switch left {
	case "", "+":
		coeff = 1
	case "-":
		coeff = -1
	default:
		n, err := strconv.Atoi(left)
		if err != nil {
			return 0, 0, err
		}
		coeff = n
	}

	if len(parts) > 1 && strings.HasPrefix(parts[1], "^") {
		n, err := strconv.Atoi(parts[1][1:])
		if err != nil {
			return 0, 0, err
		}
		exp = n
	}

	return coeff, exp, nil
}

func differentiate(p Poly) Poly {
	out := Poly{}
	for exp, coeff := range p {
		if exp == 0 {
			continue
		}
		out[exp-1] += coeff * exp
	}
	return out
}

func format(p Poly) string {
	var exps []int
	for exp, coeff := range p {
		if coeff != 0 {
			exps = append(exps, exp)
		}
	}
	if len(exps) == 0 {
		return "0"
	}

	sort.Sort(sort.Reverse(sort.IntSlice(exps)))

	var parts []string
	for _, exp := range exps {
		coeff := p[exp]
		sign := "+"
		if coeff < 0 {
			sign = "-"
			coeff = -coeff
		}

		var piece string
		switch exp {
		case 0:
			piece = fmt.Sprintf("%d", coeff)
		case 1:
			if coeff == 1 {
				piece = "x"
			} else {
				piece = fmt.Sprintf("%d*x", coeff)
			}
		default:
			if coeff == 1 {
				piece = fmt.Sprintf("x^%d", exp)
			} else {
				piece = fmt.Sprintf("%d*x^%d", coeff, exp)
			}
		}

		if len(parts) == 0 {
			if sign == "-" {
				parts = append(parts, "-"+piece)
			} else {
				parts = append(parts, piece)
			}
		} else {
			parts = append(parts, sign+" "+piece)
		}
	}

	return strings.Join(parts, " ")
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println(`usage: go run . "x^3 + 2*x + 5"`)
		os.Exit(1)
	}

	p, err := parse(os.Args[1])
	if err != nil {
		fmt.Println("parse error:", err)
		os.Exit(1)
	}

	fmt.Println(format(differentiate(p)))
}

