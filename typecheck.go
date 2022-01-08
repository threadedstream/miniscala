package main

//func check(e Exp, env map[string]bool) {
//	switch e.(type) {
//	case Lit:
//		break
//	case Var:
//		v := e.(Var)
//		if !env[v.name] {
//			panic(fmt.Errorf("unbound variable %s", v.name))
//		}
//	case Prim:
//		prim := e.(Prim)
//		if !isOperator(rune(prim.op[0])) {
//			panic(fmt.Errorf("undefined operator %c", rune(prim.op[0])))
//		}
//		for _, x := range prim.xs {
//			check(x, env)
//		}
//	}
//}
