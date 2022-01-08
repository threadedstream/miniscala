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

func checkOpValues(op Operator, v1, v2 Value) {
	if v1.valueType != Literal && v2.valueType != Literal {
		panic("v1 and v2 must be of literal type")
	}

	switch op {
	case PlusOp:
		if (v1.isString() && v2.isString()) || (v1.isFloat() && v2.isFloat()) {
			return
		}
		panic("v1 and v2 must both be of type string or float")
	case MinusOp:
		if v1.isFloat() && v2.isFloat() {
			return
		}
		panic("v1 and v2 must both be of type float")
	case MulOp:
		if v1.isFloat() && v2.isFloat() {
			return
		}
		panic("v1 and v2 must both be of type float")
	default:
		panic("unknown operation")
	}
}
