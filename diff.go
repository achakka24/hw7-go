package main

func differentiate(e Expr) Expr {
	return differentiateVar(e, "x")
}

func deriveNth(e Expr, variable string, n int) Expr {
	if n <= 0 {
		return e
	}
	cur := e
	for i := 0; i < n; i++ {
		cur = simplify(differentiateVar(cur, variable))
	}
	return cur
}

func differentiateVar(e Expr, variable string) Expr {
	switch n := e.(type) {
	case Const:
		return Const{Value: 0}
	case Var:
		if n.Name == variable {
			return Const{Value: 1}
		}
		return Const{Value: 0}
	case Neg:
		return Neg{X: differentiateVar(n.X, variable)}
	case Add:
		return Add{
			Left:  differentiateVar(n.Left, variable),
			Right: differentiateVar(n.Right, variable),
		}
	case Sub:
		return Sub{
			Left:  differentiateVar(n.Left, variable),
			Right: differentiateVar(n.Right, variable),
		}
	case Mul:
		if c, ok := asConst(n.Left); ok {
			return Mul{
				Left:  Const{Value: c},
				Right: differentiateVar(n.Right, variable),
			}
		}
		if c, ok := asConst(n.Right); ok {
			return Mul{
				Left:  Const{Value: c},
				Right: differentiateVar(n.Left, variable),
			}
		}
		return Add{
			Left: Mul{
				Left:  differentiateVar(n.Left, variable),
				Right: n.Right,
			},
			Right: Mul{
				Left:  n.Left,
				Right: differentiateVar(n.Right, variable),
			},
		}
	case Div:
		if c, ok := asConst(n.Right); ok && c != 0 {
			return Div{
				Left:  differentiateVar(n.Left, variable),
				Right: Const{Value: c},
			}
		}
		num := Sub{
			Left: Mul{
				Left:  differentiateVar(n.Left, variable),
				Right: n.Right,
			},
			Right: Mul{
				Left:  n.Left,
				Right: differentiateVar(n.Right, variable),
			},
		}
		den := Pow{
			Base:     n.Right,
			Exponent: Const{Value: 2},
		}
		return Div{Left: num, Right: den}
	case Pow:
		if exp, ok := asConst(n.Exponent); ok {
			if exp == 0 {
				return Const{Value: 0}
			}
			if exp == 1 {
				return differentiateVar(n.Base, variable)
			}
			return Mul{
				Left: Mul{
					Left:  Const{Value: exp},
					Right: Pow{Base: n.Base, Exponent: Const{Value: exp - 1}},
				},
				Right: differentiateVar(n.Base, variable),
			}
		}
		return Const{Value: 0}
	default:
		return Const{Value: 0}
	}
}
