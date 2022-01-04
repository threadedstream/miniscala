package main

import (
	"fmt"
	"strings"
	"text/scanner"
	"unicode"
)

func expected(s string) {
	panic(fmt.Errorf("expected %s", s))
}

type Reader struct {
	scanner *scanner.Scanner
}

func (r *Reader) hasNext() bool {
	return r.scanner.Peek() != scanner.EOF
}

func (r *Reader) hasNextP(predicate func(rune) bool) bool {
	return predicate(r.scanner.Peek())
}

func (r *Reader) accept(c rune) {
	if r.scanner.Peek() == c {
		r.scanner.Next()
	} else {
		expected(string(c))
	}
}

func isOperator(c rune) bool {
	return c == '+' || c == '-' || c == '*' || c == '/' || c == '%'
}

func (r *Reader) getNum() int {
	if r.hasNextP(unicode.IsDigit) {
		n := 0
		for r.hasNextP(unicode.IsDigit) {
			n = 10*n + (int)(r.scanner.Next()-'0')
		}
		return n
	} else {
		expected("number")
	}
	// shouldn't reach this point
	return 0
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
