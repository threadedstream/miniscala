package main

import "fmt"

var (
	regs   = []string{"%rax", "%rcx", "%rdx"}
	sp     = 0
	offset = 0
)

func mem(i int) string {
	if i >= offset {
		return regs[i%len(regs)]
	} else {
		return fmt.Sprintf("stack(%d)", i)
	}
}

func grow() {
	if sp == offset+len(regs) {
		fmt.Printf("push %s # evict %d\n", mem(offset), offset)
		offset += 1
	}

	sp += 1
}

func shrink() {
	sp -= 1
	if sp == offset {
		offset -= 1
		fmt.Printf("pop %s # reload %d\n", mem(offset), offset)
	}
}

func trans(e Exp) {
	switch e.(type) {
	case Lit:
		var literal = e.(Lit)
		grow()
		fmt.Printf("movq $%d, %s # %d\n", literal.x, mem(sp-1), sp-1)
	case Prim:
		var prim = e.(Prim)
		switch prim.op {
		case "+":
			trans(prim.xs[0])
			trans(prim.xs[1])
			shrink()
			fmt.Printf("addq %s, %s\n", mem(sp), mem(sp-1))
		case "-":
			trans(prim.xs[0])
			trans(prim.xs[1])
			shrink()
			fmt.Printf("subq %s, %s\n", mem(sp), mem(sp-1))
		case "*":
			trans(prim.xs[0])
			trans(prim.xs[1])
			shrink()
			fmt.Printf("imulq %s, %s\n", mem(sp), mem(sp-1))
		case "/":
			trans(prim.xs[0])
			trans(prim.xs[1])
			shrink()
			fmt.Printf("movq %s, %%rax\ncltd\nidivq %s\nmovq %%rax, %s\n", mem(sp), mem(sp-1), mem(sp))
		case "%":
			trans(prim.xs[0])
			trans(prim.xs[1])
			fmt.Printf("movq %s, %%rax\ncltd\nidivq %s\nmovq %%rdx, %s\n", mem(sp), mem(sp-1), mem(sp))
		default:
			expected(fmt.Sprintf("+, -, *, /, %%, but found %s", prim.op))
		}
	}
}
