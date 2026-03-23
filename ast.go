package main

type Expr interface {
	isExpr()
}

type Const struct {
	Value int
}

func (Const) isExpr() {}

type Var struct {
	Name string
}

func (Var) isExpr() {}

type Add struct {
	Left  Expr
	Right Expr
}

func (Add) isExpr() {}

type Sub struct {
	Left  Expr
	Right Expr
}

func (Sub) isExpr() {}

type Mul struct {
	Left  Expr
	Right Expr
}

func (Mul) isExpr() {}

type Div struct {
	Left  Expr
	Right Expr
}

func (Div) isExpr() {}

type Pow struct {
	Base     Expr
	Exponent Expr
}

func (Pow) isExpr() {}

type Neg struct {
	X Expr
}

func (Neg) isExpr() {}

func equalExpr(a, b Expr) bool {
	switch av := a.(type) {
	case Const:
		bv, ok := b.(Const)
		return ok && av.Value == bv.Value
	case Var:
		bv, ok := b.(Var)
		return ok && av.Name == bv.Name
	case Add:
		bv, ok := b.(Add)
		return ok && equalExpr(av.Left, bv.Left) && equalExpr(av.Right, bv.Right)
	case Sub:
		bv, ok := b.(Sub)
		return ok && equalExpr(av.Left, bv.Left) && equalExpr(av.Right, bv.Right)
	case Mul:
		bv, ok := b.(Mul)
		return ok && equalExpr(av.Left, bv.Left) && equalExpr(av.Right, bv.Right)
	case Div:
		bv, ok := b.(Div)
		return ok && equalExpr(av.Left, bv.Left) && equalExpr(av.Right, bv.Right)
	case Pow:
		bv, ok := b.(Pow)
		return ok && equalExpr(av.Base, bv.Base) && equalExpr(av.Exponent, bv.Exponent)
	case Neg:
		bv, ok := b.(Neg)
		return ok && equalExpr(av.X, bv.X)
	default:
		return false
	}
}
