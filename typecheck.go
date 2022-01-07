package main

import "fmt"

func check(e Exp, env map[string]bool) {
	switch e.(type) {
	case Lit:
		break
	case Var:
		v := e.(Var)
		if !env[v.name] {
			panic(fmt.Errorf("unbound variable %s", v.name))
		}
	case Prim:

	}
}
