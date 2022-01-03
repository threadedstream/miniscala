package main

import "fmt"

var (
	regs = []string{"%rbx", "%rcx", "%rdi", "%rsi", "%r8", "%r9"}
)

func transSsa(e Exp, sp int) {
	switch e.(type) {
	case Lit:
		var literal = e.(Lit)
		fmt.Printf("memory(%d) = %d\n", sp, literal.x)
	case Plus:
		plusExpr := e.(Plus)
		transSsa(plusExpr.x, sp)
		transSsa(plusExpr.y, sp+1)
		fmt.Printf("memory(%d) += memory(%d)\n", sp, sp+1)
	case Minus:
		minusExpr := e.(Minus)
		transSsa(minusExpr.x, sp)
		transSsa(minusExpr.y, sp+1)
		fmt.Printf("memory(%d) -= memory(%d)\n", sp, sp+1)
	case Times:
		timesExpr := e.(Times)
		transSsa(timesExpr.x, sp)
		transSsa(timesExpr.y, sp+1)
		fmt.Printf("memory(%d) *= memory(%d)\n", sp, sp+1)
	case Div:
		divExpr := e.(Div)
		transSsa(divExpr.x, sp)
		transSsa(divExpr.y, sp+1)
		fmt.Printf("memory(%d) /= memory(%d)\n", sp, sp+1)
	case Mod:
		modExpr := e.(Mod)
		transSsa(modExpr.x, sp)
		transSsa(modExpr.y, sp+1)
		fmt.Printf("memory(%d) %%= memory(%d)\n", sp, sp+1)
	default:
		panic(fmt.Errorf("encountered unknown node %v", e))
	}
}

func trans(e Exp, sp int) {
	switch e.(type) {
	case Lit:
		var literal = e.(Lit)
		fmt.Printf("movq $%d, %s\n", literal.x, regs[sp])
	case Plus:
		var plusExpr = e.(Plus)
		trans(plusExpr.x, sp)
		trans(plusExpr.y, sp+1)
		fmt.Printf("addq %s, %s\n", regs[sp+1], regs[sp])
	case Minus:
		minusExpr := e.(Minus)
		trans(minusExpr.x, sp)
		trans(minusExpr.y, sp+1)
		fmt.Printf("subq %s, %s", regs[sp+1], regs[sp])
	case Times:
		timesExpr := e.(Times)
		trans(timesExpr.x, sp)
		trans(timesExpr.y, sp+1)
		fmt.Printf("imulq %s, %s\n", regs[sp+1], regs[sp])
	case Div:
		divExpr := e.(Div)
		trans(divExpr.x, sp)
		trans(divExpr.y, sp+1)
		fmt.Printf("movq %s, %%rax\ncltd\nidivq %s\nmovq %%rax, %s", regs[sp], regs[sp+1], regs[sp])
	case Mod:
		modExpr := e.(Mod)
		trans(modExpr.x, sp)
		trans(modExpr.y, sp+1)
		fmt.Printf("movq %s, %%rax\ncltd\nidivq %s\nmovq %%rdx, %s", regs[sp], regs[sp+1], regs[sp])
	default:
		panic(fmt.Errorf("encountered unknown node %v", e))
	}
}
