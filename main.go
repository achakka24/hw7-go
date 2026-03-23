package main

import (
	"fmt"
	"os"
)

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

	result := simplify(differentiate(p))
	fmt.Println(format(result))
}

