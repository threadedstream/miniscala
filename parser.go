package main

import (
	"fmt"
	"strings"
	"text/scanner"
)

type Parser struct {
	in *TokenScanner
}

func (p *Parser) followsChar(c rune) bool {
	var delim = Delim{x: c}
	if p.in.peekToken == delim {
		p.in.next()
		return true
	}

	return false
}

func (p *Parser) followsString(s string) bool {
	var kwd = Keyword{x: s}
	if p.in.peekToken == kwd {
		p.in.next()
		return true
	}

	return false
}

func (p *Parser) requireChar(c rune) {
	var delim = Delim{x: c}
	if p.in.peekToken == delim {
		p.in.next()
		return
	}

	panic(fmt.Errorf("expected %c", c))
}

func (p *Parser) requireString(s string) {
	var kwd = Keyword{x: s}
	if p.in.peekToken == kwd {
		p.in.next()
		return
	}
	panic(fmt.Errorf("expected %s", s))
}

func (p *Parser) isName(x Token) bool {
	switch x.(type) {
	case Ident:
		return true
	default:
		return false
	}
}

func (p *Parser) isNum(x Token) bool {
	switch x.(type) {
	case Number:
		return true
	default:
		return false
	}
}

func (p *Parser) isInfixOp(min int, x Token) bool {
	switch x.(type) {
	case Ident:
		var ident = x.(Ident)
		return prec(ident.x) >= min
	default:
		return false
	}
}

func (p *Parser) name() string {
	if !p.in.hasNextP(p.isName) {
		panic("expected name")
	}
	var ident = p.in.next().(Ident)
	return ident.x
}

func (p *Parser) atom() Exp {
	switch p.in.peekToken.(type) {
	case Number:
		var num = p.in.peekToken.(Number)
		p.in.next()
		return Lit{x: num.x}
	case Ident:
		var ident = p.in.peekToken.(Ident)
		p.in.next()
		return Var{name: ident.x}
	case Delim:
		var delim = p.in.peekToken.(Delim)
		switch delim.x {
		case '(':
			p.in.next()
			var x = p.simpl()
			p.requireChar(')')
			return x
		case '{':
			p.in.next()
			var x = p.expr()
			p.requireChar('}')
			return x
		}
	}
}

func (p *Parser) simpl() Exp {
	return EOF{}
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
