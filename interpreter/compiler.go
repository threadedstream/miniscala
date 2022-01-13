package interpreter

//var (
//	regs   = []string{"%rax", "%rcx", "%rdx"}
//	sp     = 0
//	offset = 0
//)
//
//func mem(i int) string {
//	if i >= offset {
//		return regs[i%len(regs)]
//	} else {
//		return fmt.Sprintf("stack(%d)", i)
//	}
//}
//
//func grow() {
//	if sp == offset+len(regs) {
//		fmt.Printf("push %s # evict %d\n", mem(offset), offset)
//		offset += 1
//	}
//
//	sp += 1
//}
//
//func shrink(dst int) {
//	if offset != dst {
//		fmt.Printf("addq $%d, %%rsp # reset stack %d to %d", (offset-dst)*8, offset, dst)
//	}
//	sp = dst
//	offset = dst
//}
//
//func transPrim(op string) {
//	switch op {
//	case "+":
//		fmt.Printf("addq %s, %s\n", mem(sp), mem(sp-1))
//	case "-":
//		fmt.Printf("subq %s, %s\n", mem(sp), mem(sp-1))
//	case "*":
//		fmt.Printf("imulq %s, %s\n", mem(sp), mem(sp-1))
//	case "/":
//		fmt.Printf("movq %s, %%rax\ncltd\nidivq %s\nmovq %%rax, %s\n", mem(sp), mem(sp-1), mem(sp))
//	case "%":
//		fmt.Printf("movq %s, %%rax\ncltd\nidivq %s\nmovq %%rdx, %s\n", mem(sp), mem(sp-1), mem(sp))
//	}
//}
//
//func trans1(e Exp, dst int, env map[string]int) {
//	switch e.(type) {
//	case Let:
//		var let = e.(Let)
//		trans(let.rhs, env)
//		env[let.name] = sp - 1
//		trans1(let.body, dst, env)
//	case Lit:
//		var lit = e.(Lit)
//		shrink(dst)
//		fmt.Printf("movq $%d, %s", lit.x, mem(dst))
//	case Var:
//		var variable = e.(Var)
//		shrink(dst)
//		if mem(env[variable.name]) != mem(dst) {
//			fmt.Printf("movq %s, %s", mem(env[variable.name]), mem(dst))
//		}
//	default:
//		trans(e, env)
//		var src = sp - 1
//		shrink(dst)
//		if mem(src) != mem(dst) {
//			fmt.Printf("movq %s, %s # drop %d to %d", mem(src), mem(dst), src, dst)
//		}
//	}
//
//}
//
//func trans(e Exp, env map[string]int) {
//	switch e.(type) {
//	case Lit:
//		var literal = e.(Lit)
//		grow()
//		fmt.Printf("movq $%d, %s # %d\n", literal.x, mem(sp-1), sp-1)
//	case Var:
//		var variable = e.(Var)
//		grow()
//		fmt.Printf("movq %s, %s # %d", mem(env[variable.name]), mem(sp-1), sp-1)
//	case Prim:
//		var prim = e.(Prim)
//		trans(prim.xs[0], env)
//		trans(prim.xs[1], env)
//		shrink(sp - 1)
//		transPrim(prim.op)
//	case Let:
//		var let = e.(Let)
//		trans(let.rhs, env)
//		env[let.name] = sp - 1
//		trans1(let.body, sp-1, env)
//	}
//}
