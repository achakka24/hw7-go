package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	variable := flag.String("var", "x", "differentiate with respect to variable name")
	nth := flag.Int("nth", 1, "number of derivative applications")
	steps := flag.Bool("steps", false, "print intermediate derivative/simplification steps")
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Println(`usage: go run . [--var x] [--nth 1] [--steps] "expression"`)
		os.Exit(1)
	}
	if *nth < 1 {
		fmt.Println("error: --nth must be >= 1")
		os.Exit(1)
	}
	if *variable == "" {
		fmt.Println("error: --var must not be empty")
		os.Exit(1)
	}

	expr, err := parse(flag.Arg(0))
	if err != nil {
		fmt.Println("parse error:", err)
		os.Exit(1)
	}

	if *steps {
		fmt.Println("parsed:", format(expr))
	}

	result := expr
	for i := 1; i <= *nth; i++ {
		raw := differentiateVar(result, *variable)
		if *steps {
			fmt.Printf("d%d raw: %s\n", i, format(raw))
		}
		result = simplify(raw)
		if *steps {
			fmt.Printf("d%d simplified: %s\n", i, format(result))
		}
	}

	fmt.Println(format(result))
}
