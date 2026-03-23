package main

func simplify(e Expr) Expr {
	cur := e
	for i := 0; i < 128; i++ {
		next := simplifyOnce(cur)
		if equalExpr(cur, next) {
			return next
		}
		cur = next
	}
	return cur
}

func simplifyOnce(e Expr) Expr {
	switch n := e.(type) {
	case Const, Var:
		return e

	case Neg:
		x := simplifyOnce(n.X)
		if v, ok := asConst(x); ok {
			return Const{Value: -v}
		}
		if inner, ok := x.(Neg); ok {
			return inner.X
		}
		return Neg{X: x}

	case Add:
		l := simplifyOnce(n.Left)
		r := simplifyOnce(n.Right)

		if isZero(l) {
			return r
		}
		if isZero(r) {
			return l
		}
		if lv, lok := asConst(l); lok {
			if rv, rok := asConst(r); rok {
				return Const{Value: lv + rv}
			}
		}
		if rn, ok := r.(Neg); ok {
			return simplifyOnce(Sub{Left: l, Right: rn.X})
		}
		if ln, ok := l.(Neg); ok {
			return simplifyOnce(Sub{Left: r, Right: ln.X})
		}
		if equalExpr(l, r) {
			return simplifyOnce(Mul{Left: Const{Value: 2}, Right: l})
		}

		if lc, lt, ok := splitConstMul(l); ok {
			if rc, rt, ok2 := splitConstMul(r); ok2 && equalExpr(lt, rt) {
				return simplifyOnce(Mul{Left: Const{Value: lc + rc}, Right: lt})
			}
		}
		return Add{Left: l, Right: r}

	case Sub:
		l := simplifyOnce(n.Left)
		r := simplifyOnce(n.Right)

		if isZero(r) {
			return l
		}
		if isZero(l) {
			return simplifyOnce(Neg{X: r})
		}
		if equalExpr(l, r) {
			return Const{Value: 0}
		}
		if lv, lok := asConst(l); lok {
			if rv, rok := asConst(r); rok {
				return Const{Value: lv - rv}
			}
		}
		if rn, ok := r.(Neg); ok {
			return simplifyOnce(Add{Left: l, Right: rn.X})
		}

		if lc, lt, ok := splitConstMul(l); ok {
			if rc, rt, ok2 := splitConstMul(r); ok2 && equalExpr(lt, rt) {
				return simplifyOnce(Mul{Left: Const{Value: lc - rc}, Right: lt})
			}
		}
		return Sub{Left: l, Right: r}

	case Mul:
		l := simplifyOnce(n.Left)
		r := simplifyOnce(n.Right)

		if isZero(l) || isZero(r) {
			return Const{Value: 0}
		}
		if isOne(l) {
			return r
		}
		if isOne(r) {
			return l
		}
		if isNegOne(l) {
			return simplifyOnce(Neg{X: r})
		}
		if isNegOne(r) {
			return simplifyOnce(Neg{X: l})
		}
		if lv, lok := asConst(l); lok {
			if rv, rok := asConst(r); rok {
				return Const{Value: lv * rv}
			}
		}
		if ln, ok := l.(Neg); ok {
			return simplifyOnce(Neg{X: Mul{Left: ln.X, Right: r}})
		}
		if rn, ok := r.(Neg); ok {
			return simplifyOnce(Neg{X: Mul{Left: l, Right: rn.X}})
		}
		if equalExpr(l, r) {
			return simplifyOnce(Pow{Base: l, Exponent: Const{Value: 2}})
		}

		// Prefer constant coefficient on the left.
		if _, ok := asConst(r); ok {
			if _, ok2 := asConst(l); !ok2 {
				l, r = r, l
			}
		}

		// Combine constant factors: c1 * (c2 * t) => (c1*c2) * t
		if lc, ok := asConst(l); ok {
			if rm, ok2 := r.(Mul); ok2 {
				if rc, ok3 := asConst(rm.Left); ok3 {
					return simplifyOnce(Mul{Left: Const{Value: lc * rc}, Right: rm.Right})
				}
				if rc, ok3 := asConst(rm.Right); ok3 {
					return simplifyOnce(Mul{Left: Const{Value: lc * rc}, Right: rm.Left})
				}
			}
		}
		// Normalize c*x*x => c*x^2 so quotient cancellation can trigger.
		if lm, ok := l.(Mul); ok {
			if lc, ok2 := asConst(lm.Left); ok2 && equalExpr(lm.Right, r) {
				return simplifyOnce(Mul{
					Left:  Const{Value: lc},
					Right: Pow{Base: r, Exponent: Const{Value: 2}},
				})
			}
			if lc, ok2 := asConst(lm.Right); ok2 && equalExpr(lm.Left, r) {
				return simplifyOnce(Mul{
					Left:  Const{Value: lc},
					Right: Pow{Base: r, Exponent: Const{Value: 2}},
				})
			}
		}
		if rm, ok := r.(Mul); ok {
			if rc, ok2 := asConst(rm.Left); ok2 && equalExpr(rm.Right, l) {
				return simplifyOnce(Mul{
					Left:  Const{Value: rc},
					Right: Pow{Base: l, Exponent: Const{Value: 2}},
				})
			}
			if rc, ok2 := asConst(rm.Right); ok2 && equalExpr(rm.Left, l) {
				return simplifyOnce(Mul{
					Left:  Const{Value: rc},
					Right: Pow{Base: l, Exponent: Const{Value: 2}},
				})
			}
		}

		// x * x^n => x^(n+1), and x^a * x^b => x^(a+b)
		if rp, ok := r.(Pow); ok && equalExpr(l, rp.Base) {
			if rv, ok2 := asConst(rp.Exponent); ok2 {
				return simplifyOnce(Pow{Base: l, Exponent: Const{Value: rv + 1}})
			}
		}
		if lp, ok := l.(Pow); ok && equalExpr(r, lp.Base) {
			if lv, ok2 := asConst(lp.Exponent); ok2 {
				return simplifyOnce(Pow{Base: r, Exponent: Const{Value: lv + 1}})
			}
		}
		if lp, ok := l.(Pow); ok {
			if rp, ok2 := r.(Pow); ok2 && equalExpr(lp.Base, rp.Base) {
				lv, lok := asConst(lp.Exponent)
				rv, rok := asConst(rp.Exponent)
				if lok && rok {
					return simplifyOnce(Pow{Base: lp.Base, Exponent: Const{Value: lv + rv}})
				}
			}
		}
		return Mul{Left: l, Right: r}

	case Div:
		l := simplifyOnce(n.Left)
		r := simplifyOnce(n.Right)

		if isZero(l) {
			return Const{Value: 0}
		}
		if isOne(r) {
			return l
		}
		if equalExpr(l, r) {
			return Const{Value: 1}
		}
		if rv, ok := asConst(r); ok {
			if rv == -1 {
				return simplifyOnce(Neg{X: l})
			}
			if lv, ok2 := asConst(l); ok2 && rv != 0 && lv%rv == 0 {
				return Const{Value: lv / rv}
			}
		}
		if lv, ok := asConst(l); ok && lv < 0 {
			return simplifyOnce(Neg{
				X: Div{
					Left:  Const{Value: -lv},
					Right: r,
				},
			})
		}
		if ln, ok := l.(Neg); ok {
			return simplifyOnce(Neg{X: Div{Left: ln.X, Right: r}})
		}
		if rn, ok := r.(Neg); ok {
			return simplifyOnce(Neg{X: Div{Left: l, Right: rn.X}})
		}

		// (a*b)/b => a, including c*t / t => c
		if lm, ok := l.(Mul); ok {
			if equalExpr(lm.Left, r) {
				return lm.Right
			}
			if equalExpr(lm.Right, r) {
				return lm.Left
			}
			if c, ok2 := asConst(lm.Left); ok2 && equalExpr(lm.Right, r) {
				return Const{Value: c}
			}
			if c, ok2 := asConst(lm.Right); ok2 && equalExpr(lm.Left, r) {
				return Const{Value: c}
			}
		}

		// x^a / x^b => x^(a-b)
		if lp, ok := l.(Pow); ok {
			if rp, ok2 := r.(Pow); ok2 && equalExpr(lp.Base, rp.Base) {
				lv, lok := asConst(lp.Exponent)
				rv, rok := asConst(rp.Exponent)
				if lok && rok {
					diff := lv - rv
					switch {
					case diff == 0:
						return Const{Value: 1}
					case diff > 0:
						return simplifyOnce(Pow{Base: lp.Base, Exponent: Const{Value: diff}})
					default:
						return simplifyOnce(Div{
							Left:  Const{Value: 1},
							Right: Pow{Base: lp.Base, Exponent: Const{Value: -diff}},
						})
					}
				}
			}
		}

		// x / x^n => 1 / x^(n-1)
		if rp, ok := r.(Pow); ok && equalExpr(l, rp.Base) {
			if rv, ok2 := asConst(rp.Exponent); ok2 {
				if rv == 1 {
					return Const{Value: 1}
				}
				return simplifyOnce(Div{
					Left:  Const{Value: 1},
					Right: Pow{Base: l, Exponent: Const{Value: rv - 1}},
				})
			}
		}
		// c*x / x^n => c / x^(n-1)
		if lm, ok := l.(Mul); ok {
			if rp, ok2 := r.(Pow); ok2 {
				if c, ok3 := asConst(lm.Left); ok3 && equalExpr(lm.Right, rp.Base) {
					if rv, ok4 := asConst(rp.Exponent); ok4 {
						if rv == 1 {
							return Const{Value: c}
						}
						return simplifyOnce(Div{
							Left:  Const{Value: c},
							Right: Pow{Base: rp.Base, Exponent: Const{Value: rv - 1}},
						})
					}
				}
				if c, ok3 := asConst(lm.Right); ok3 && equalExpr(lm.Left, rp.Base) {
					if rv, ok4 := asConst(rp.Exponent); ok4 {
						if rv == 1 {
							return Const{Value: c}
						}
						return simplifyOnce(Div{
							Left:  Const{Value: c},
							Right: Pow{Base: rp.Base, Exponent: Const{Value: rv - 1}},
						})
					}
				}
				// c*x^a / x^b => c / x^(b-a)
				if c, ok3 := asConst(lm.Left); ok3 {
					if lp, ok4 := lm.Right.(Pow); ok4 && equalExpr(lp.Base, rp.Base) {
						if a, ok5 := asConst(lp.Exponent); ok5 {
							if b, ok6 := asConst(rp.Exponent); ok6 {
								diff := b - a
								switch {
								case diff == 0:
									return Const{Value: c}
								case diff > 0:
									return simplifyOnce(Div{
										Left:  Const{Value: c},
										Right: Pow{Base: rp.Base, Exponent: Const{Value: diff}},
									})
								default:
									return simplifyOnce(Mul{
										Left:  Const{Value: c},
										Right: Pow{Base: rp.Base, Exponent: Const{Value: -diff}},
									})
								}
							}
						}
					}
				}
				if c, ok3 := asConst(lm.Right); ok3 {
					if lp, ok4 := lm.Left.(Pow); ok4 && equalExpr(lp.Base, rp.Base) {
						if a, ok5 := asConst(lp.Exponent); ok5 {
							if b, ok6 := asConst(rp.Exponent); ok6 {
								diff := b - a
								switch {
								case diff == 0:
									return Const{Value: c}
								case diff > 0:
									return simplifyOnce(Div{
										Left:  Const{Value: c},
										Right: Pow{Base: rp.Base, Exponent: Const{Value: diff}},
									})
								default:
									return simplifyOnce(Mul{
										Left:  Const{Value: c},
										Right: Pow{Base: rp.Base, Exponent: Const{Value: -diff}},
									})
								}
							}
						}
					}
				}
			}
		}

		// x^n / x => x^(n-1)
		if lp, ok := l.(Pow); ok && equalExpr(lp.Base, r) {
			if lv, ok2 := asConst(lp.Exponent); ok2 {
				if lv == 1 {
					return Const{Value: 1}
				}
				return simplifyOnce(Pow{Base: r, Exponent: Const{Value: lv - 1}})
			}
		}
		return Div{Left: l, Right: r}

	case Pow:
		b := simplifyOnce(n.Base)
		e2 := simplifyOnce(n.Exponent)

		if ev, ok := asConst(e2); ok {
			if ev == 0 {
				return Const{Value: 1}
			}
			if ev == 1 {
				return b
			}
			if bv, ok2 := asConst(b); ok2 && ev >= 0 {
				return Const{Value: intPow(bv, ev)}
			}
			if ev < 0 {
				return simplifyOnce(Div{
					Left:  Const{Value: 1},
					Right: Pow{Base: b, Exponent: Const{Value: -ev}},
				})
			}
		}
		if isOne(b) {
			return Const{Value: 1}
		}
		if isZero(b) {
			if ev, ok := asConst(e2); ok && ev > 0 {
				return Const{Value: 0}
			}
		}
		if bp, ok := b.(Pow); ok {
			if ev1, ok1 := asConst(bp.Exponent); ok1 {
				if ev2, ok2 := asConst(e2); ok2 {
					return simplifyOnce(Pow{Base: bp.Base, Exponent: Const{Value: ev1 * ev2}})
				}
			}
		}
		return Pow{Base: b, Exponent: e2}

	default:
		return e
	}
}

func asConst(e Expr) (int, bool) {
	v, ok := e.(Const)
	if !ok {
		return 0, false
	}
	return v.Value, true
}

func isZero(e Expr) bool {
	v, ok := asConst(e)
	return ok && v == 0
}

func isOne(e Expr) bool {
	v, ok := asConst(e)
	return ok && v == 1
}

func isNegOne(e Expr) bool {
	v, ok := asConst(e)
	return ok && v == -1
}

func splitConstMul(e Expr) (int, Expr, bool) {
	if _, ok := asConst(e); ok {
		return 0, nil, false
	}
	if m, ok := e.(Mul); ok {
		if c, ok2 := asConst(m.Left); ok2 {
			return c, m.Right, true
		}
		if c, ok2 := asConst(m.Right); ok2 {
			return c, m.Left, true
		}
	}
	return 1, e, true
}

func intPow(base, exp int) int {
	if exp == 0 {
		return 1
	}
	result := 1
	b := base
	e := exp
	for e > 0 {
		if e%2 == 1 {
			result *= b
		}
		b *= b
		e /= 2
	}
	return result
}
