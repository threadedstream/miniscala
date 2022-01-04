package main

const (
	MEMSIZE = 1024
)

var (
	memory [MEMSIZE]int
)

func eval(e Exp, sp int, env map[string]int) {
	switch e.(type) {
	case Lit:
		var literal = e.(Lit)
		memory[sp] = literal.x
	// TODO(threadedstream): make a single huge Prim case
	case Prim:
		var prim = e.(Prim)
		switch prim.op {
		case "+":
			eval(prim.xs[0], sp, env)
			eval(prim.xs[1], sp+1, env)
			memory[sp] += memory[sp+1]
		case "-":
			eval(prim.xs[0], sp, env)
			eval(prim.xs[1], sp+1, env)
			memory[sp] -= memory[sp+1]
		case "*":
			eval(prim.xs[0], sp, env)
			eval(prim.xs[1], sp+1, env)
			memory[sp] *= memory[sp+1]
		case "/":
			eval(prim.xs[0], sp, env)
			eval(prim.xs[1], sp+1, env)
			memory[sp] /= memory[sp+1]
		case "%":
			eval(prim.xs[0], sp, env)
			eval(prim.xs[1], sp, env)
			memory[sp] %= memory[sp+1]
		}

	case Let:
		var let = e.(Let)
		eval(let.rhs, sp, env)
		env[let.name] = sp
		eval(let.body, sp+1, env)
		memory[sp] = memory[sp+1]
	case Var:
		var variable = e.(Var)
		memory[sp] = memory[env[variable.name]]
	}
}
