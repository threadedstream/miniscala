package main

type Exp interface {
}

type (
	Plus struct {
		Exp
		x Exp
		y Exp
	}

	Minus struct {
		Exp
		x Exp
		y Exp
	}

	Times struct {
		Exp
		y Exp
		x Exp
	}

	Div struct {
		Exp
		x Exp
		y Exp
	}

	Mod struct {
		Exp
		x Exp
		y Exp
	}

	Lit struct {
		Exp
		x int
	}
)
