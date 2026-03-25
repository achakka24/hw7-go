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

With extensions:

```bash
go run . --var y "x^3+y^3"
go run . --nth 2 "x^4"
go run . --steps "x^2+1/x"
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

## Extra Credit Extensions

The project includes three small extensions beyond baseline HW7 behavior:

1. `--var <name>`: choose differentiation variable (example: `--var y`).
2. `--nth <n>`: compute nth derivative (example: `--nth 2`).
3. `--steps`: print intermediate raw and simplified derivatives for demo/explanation.
