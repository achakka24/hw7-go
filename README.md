# HW7 Go

Go implementation of Homework 7 symbolic differentiation and simplification.

Implemented pipeline:

- `parse` for expressions with `+ - * / ^`, parentheses, unary minus, constants, identifiers
- `differentiate` with respect to `x`
- `simplify` with repeated rewrite rules
- `print` to canonical-style algebra output

## Run

```bash
go run . "x^4+2*x^3-x^2+5*x-1/x"
```

## Test

```bash
go test ./...
```

## Covered Homework-style Cases

- `deriv(x^2) -> 2*x`
- `deriv((x*2*x)/x) -> 2`
- `deriv(x^4+2*x^3-x^2+5*x-1/x) -> 4*x^3+6*x^2-2*x+5+1/x^2`
- `deriv(4*x^3+6*x^2-2*x+5+1/x^2) -> 12*x^2+12*x-2-2/x^3`
- `deriv(12*x^2+12*x-2-2/x^3) -> 24*x+12+6/x^4`

Simplify checks included:

- `simplify(5-x*(3/3)+2) -> 5-x+2`
- `simplify(1*x-0/3+2) -> x+2`
