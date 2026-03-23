package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println(`usage: go run . "expression"`)
		os.Exit(1)
	}

	expr, err := parse(os.Args[1])
	if err != nil {
		fmt.Println("parse error:", err)
		os.Exit(1)
	}

	result := simplify(differentiate(expr))
	fmt.Println(format(result))
}
