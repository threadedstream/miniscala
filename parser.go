package main

import (
	"fmt"
	"strings"
	"text/scanner"
)

type Parser struct {
	reader *Reader
}

func (p *Parser) follows(c rune) bool {
	//if p.reader.hasNextP()
	return false
}

func expected(s string) {
	panic(fmt.Errorf("expected %s", s))
}

func (r *Reader) accept(c rune) {
	if r.scanner.Peek() == c {
		r.scanner.Next()
	} else {
		expected(string(c))
	}
}

func (r *Reader) atom() Exp {
	if r.scanner.Peek() == '(' {
		r.accept('(')
		res := r.expr()
		r.accept(')')
		return res
	} else {
		return Lit{x: r.getNum()}
	}
}

func prec(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	}

	return 0
}

// 0 - right, 1 - left
func assoc(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/", "%":
		return 1
	}
	return -1
}

func (r *Reader) binOp(min int) Exp {
	res := r.atom()
	for r.hasNextP(isOperator) && prec(string(r.scanner.Peek())) >= min {
		op := r.scanner.Peek()
		r.scanner.Next()
		nextMin := prec(string(op)) + assoc(string(op))
		res = Prim{op: string(op), xs: []Exp{res, r.binOp(nextMin)}}
	}

	return res
}

func (r *Reader) expr() Exp {
	return r.binOp(0)
}

func parse(code string) Exp {
	reader := &Reader{
		scanner: new(scanner.Scanner),
	}
	reader.scanner = reader.scanner.Init(strings.NewReader(code))
	res := reader.expr()
	if reader.hasNext() {
		expected("EOF")
	}
	return res
}
