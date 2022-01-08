package main

//type Parser struct {
//	in *TokenScanner
//}

//func (p *Parser) skippedNewLine() bool {
//	return p.in.peekToken.pos().start.Line != p.in.peekToken.pos().gap.Line
//}
//
//func (p *Parser) followsChar(c rune) bool {
//	var delim = Delim{x: c}
//	if p.in.peekToken == delim {
//		p.in.next()
//		return true
//	} else if c == ';' && p.skippedNewLine() {
//		return true
//	}
//
//	return false
//}
//
//func (p *Parser) followsString(s string) bool {
//	var kwd = Keyword{x: s}
//	if p.in.peekToken == kwd {
//		p.in.next()
//		return true
//	}
//
//	return false
//}
//
//func (p *Parser) requireChar(c rune) {
//	var delim = Delim{x: c}
//	if p.in.peekToken == delim {
//		p.in.next()
//		return
//	}
//
//	panic(fmt.Errorf("expected %c", c))
//}
//
//func (p *Parser) requireString(s string) {
//	var kwd = Keyword{x: s}
//	if p.in.peekToken == kwd {
//		p.in.next()
//		return
//	}
//	panic(fmt.Errorf("expected %s", s))
//}
//
//func (p *Parser) isName(x Token) bool {
//	switch x.(type) {
//	case Ident:
//		return true
//	default:
//		return false
//	}
//}
//
//func (p *Parser) isNum(x Token) bool {
//	switch x.(type) {
//	case Number:
//		return true
//	default:
//		return false
//	}
//}
//
//func (p *Parser) isInfixOp(min int, x Token) bool {
//	switch x.(type) {
//	case Ident:
//		var ident = x.(Ident)
//		return prec(ident.x) >= min
//	default:
//		return false
//	}
//}
//
//// val x = 4; x + 5
//
//func (p *Parser) name() string {
//	if !p.in.hasNextP(p.isName) {
//		panic("expected name")
//	}
//	var ident = p.in.peekToken.(Ident)
//	p.in.next()
//	return ident.x
//}
//
//func (p *Parser) atom() Exp {
//	switch p.in.peekToken.(type) {
//	case Number:
//		var num = p.in.peekToken.(Number)
//		p.in.next()
//		return Lit{x: num.x}
//	case Ident:
//		var ident = p.in.peekToken.(Ident)
//		p.in.next()
//		return Var{name: ident.x}
//	case Delim:
//		var delim = p.in.peekToken.(Delim)
//		switch delim.x {
//		case '(':
//			p.in.next()
//			var x = p.simpl()
//			p.requireChar(')')
//			return x
//		case '{':
//			p.in.next()
//			var x = p.expr()
//			p.requireChar('}')
//			return x
//		default:
//			errorPos := p.in.pos()
//			panic(fmt.Errorf("[%d:%d] expected '(' or '{', but got %c", errorPos.Line, errorPos.Column, delim.x))
//		}
//	default:
//		errorPos := p.in.pos()
//		panic(fmt.Errorf("[%d:%d] expected number, identifier, or delimiter, but got %s", errorPos.Line, errorPos.Column, p.in.peekToken.str()))
//	}
//}
//
//func (p *Parser) simpl() Exp {
//	return p.binOp(0)
//}
//
//func prec(op string) int {
//	switch op {
//	case "+", "-":
//		return 1
//	case "*", "/":
//		return 2
//	}
//
//	return 0
//}
//
//// 0 - right, 1 - left
//func assoc(op string) int {
//	switch op {
//	case "+", "-":
//		return 1
//	case "*", "/", "%":
//		return 1
//	}
//	return -1
//}
//
//func (p *Parser) binOp(min int) Exp {
//	res := p.atom()
//	for p.in.hasNextP(func(token Token) bool {
//		return p.isInfixOp(min, token)
//	}) && !p.skippedNewLine() {
//		op := p.name()
//		nextMin := prec(op) + assoc(op)
//		rhs := p.binOp(nextMin)
//		res = Prim{op: op, xs: []Exp{res, rhs}}
//	}
//
//	return res
//}
//
//func (p *Parser) expr() Exp {
//	switch p.in.peekToken.(type) {
//	case Keyword:
//		var kwd = p.in.peekToken.(Keyword)
//		if kwd.x == "val" {
//			p.requireString("val")
//			var name = p.name()
//			p.requireChar('=')
//			var rhs = p.simpl()
//			p.requireChar(';')
//			var body = p.expr()
//			return Let{name: name, rhs: rhs, body: body}
//		}
//	default:
//		return p.simpl()
//	}
//	return EmptyExp{}
//}
//
//func parse(path string) Exp {
//	stream, err := os.Open(path)
//	if err != nil {
//		panic("no file with such path was found")
//	}
//
//	var parser = &Parser{
//		in: newTokenScanner(stream),
//	}
//	res := parser.expr()
//	if parser.in.hasNext() {
//		panic("expected EOF")
//	}
//	return res
//}
