package main

type (
	Exp interface {
	}

	Lit struct {
		Exp
		x int
	}

	Prim struct {
		Exp
		op string
		xs []Exp
	}

	Var struct {
		Exp
		name string
	}

	Let struct {
		Exp
		name string
		rhs  Exp
		body Exp
	}

	EmptyExp struct {
		Exp
	}
)
