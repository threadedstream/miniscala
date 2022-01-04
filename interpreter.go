package main

import "fmt"

const (
	MEMSIZE = 1024
)

var (
	memory [MEMSIZE]int
)

func eval(e Exp, sp int) {
	switch e.(type) {
	case Lit:
		var literal Lit = e.(Lit)
		memory[sp] = literal.x
	// TODO(threadedstream): make a single huge Prim case
	case Plus:
		var plus = e.(Plus)
		eval(plus.x, sp)
		eval(plus.y, sp+1)
		memory[sp] += memory[sp+1]
	case Minus:
		var minus = e.(Minus)
		eval(minus.x, sp)
		eval(minus.y, sp+1)
		memory[sp] -= memory[sp+1]
	case Times:
		var times = e.(Times)
		eval(times.x, sp)
		eval(times.y, sp+1)
		memory[sp] *= memory[sp+1]
	case Div:
		var div = e.(Div)
		eval(div.x, sp)
		eval(div.y, sp+1)
		memory[sp] /= memory[sp+1]
	case Mod:
		var mod = e.(Mod)
		eval(mod.x, sp)
		eval(mod.y, sp+1)
		memory[sp] %= memory[sp+1]
	default:
		panic(fmt.Errorf("unknown node %v", e))
	}
}
