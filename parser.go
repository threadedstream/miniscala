package main

import (
	"fmt"
	"io"
	"reflect"
)

type Parser struct {
	tokenStream []Token
	currIdx     int
}

//func (p *Parser) stmt() Stmt {
//	switch p.in.peekToken.(type) {
//	case *TokenVal:
//		return p.valDecl()
//	case *TokenVar:
//		return p.varDecl()
//	}
//	return nil
//}

func newParser(reader io.Reader) *Parser {
	scanner := newCharScanner(reader)
	return &Parser{
		tokenStream: scanner.Tokenize(),
		currIdx:     0,
	}
}

func (p *Parser) peek() Token {
	if p.currIdx+1 < len(p.tokenStream) {
		return p.tokenStream[p.currIdx+1]
	}
	return &TokenEOF{}
}

func (p *Parser) next() Token {
	token := p.peek()
	p.currIdx += 1
	return token
}

func (p *Parser) curr() Token {
	return p.tokenStream[p.currIdx]
}

func (p *Parser) consume(token Token) {
	if reflect.TypeOf(p.curr()) == reflect.TypeOf(token) {
		p.next()
	} else {
		panic(fmt.Errorf("expected token of type %v but got %v", reflect.TypeOf(p.curr()), reflect.TypeOf(token)))
	}
}

func (p *Parser) binOp(min int) Node {

}

func (p *Parser) simp() Node {
	return p.binOp(0)
}

func (p *Parser) atom() Node {
	switch p.curr().(type) {
	case *TokenNumber:
		tokenNum := p.curr().(*TokenNumber)
		p.consume(&TokenNumber{})
		return &BasicLit{
			value: tokenNum.value,
			kind:  FloatLit,
		}
	case *TokenOpenParen:
		p.consume(&TokenOpenParen{})
		simpNode := p.simp()
		p.consume(&TokenCloseParen{})
		return simpNode
	case *TokenOpenBrace:
		p.consume(&TokenOpenBrace{})
		expr := p.expr()
		p.consume(&TokenCloseBrace{})
		return expr
	case *TokenIdent:
		ident := p.curr().(*TokenIdent)
		return &Name{
			value: ident.value,
		}
	default:
		panic("unknown node in atom()")
	}
}

func (p *Parser) expr() Node {

}

func (p *Parser) varDecl() *VarDeclStmt {
	return nil
}

func (p *Parser) valDecl() *ValDeclStmt {
	return nil
}

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

func isComparisonOp(op Operator) bool {
	switch op {
	default:
		return false
	case GreaterThan, GreaterThanOrEqual, LessThan, LessThanOrEqual, Equal, NotEqual:
		return true
	}
}
