package syntax

import (
	"fmt"
	"github.com/ThreadedStream/miniscala/assert"
	"os"
	"reflect"
)

type AssociativityType int

const (
	RightAssociative AssociativityType = iota
	LeftAssociative
)

type Parser struct {
	tokenStream []Token
	currIdx     int
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
		panic(fmt.Errorf("expected %s but got %s", tokToString(p.curr()), tokToString(token)))
	}
}

func (p *Parser) match(token Token) bool {
	if reflect.TypeOf(p.curr()) != reflect.TypeOf(token) {
		return false
	}

	return true
}

func (p *Parser) binOp(min int) Node {
	res := p.atom()
	for IsOperator(p.curr()) && prec(p.curr()) >= min {
		op := tokenToOperator(p.curr())
		nextMin := prec(p.curr()) + int(assoc(p.curr()))
		p.next()
		res = &Operation{
			Op:  op,
			Lhs: res,
			Rhs: p.binOp(nextMin),
		}
	}

	return res
}

func (p *Parser) stmt() Stmt {
	switch p.curr().(type) {
	case *TokenVar:
		return p.varDeclStmt()
	case *TokenVal:
		return p.valDeclStmt()
	case *TokenIf:
		return p.ifStmt()
	case *TokenWhile:
		return p.whileStmt()
	case *TokenOpenBrace:
		return p.blockStmt()
	default:
		if (reflect.TypeOf(p.curr()) == reflect.TypeOf(&TokenIdent{})) &&
			(reflect.TypeOf(p.peek()) == reflect.TypeOf(&TokenAssign{})) {
			return p.assignment()
		} else {
			panic(fmt.Errorf("expected var, val, if, or ident, but got %s", tokToString(p.curr())))
		}
	}
}

func (p *Parser) atom() Node {
	switch p.curr().(type) {
	case *TokenNumber:
		tokenNum := p.curr().(*TokenNumber)
		p.consume(&TokenNumber{})
		return &BasicLit{
			Value: tokenNum.value,
			Kind:  FloatLit,
		}
	case *TokenOpenParen:
		p.consume(&TokenOpenParen{})
		simpNode := p.expr()
		p.consume(&TokenCloseParen{})
		return simpNode
	case *TokenOpenBrace:
		p.consume(&TokenOpenBrace{})
		expr := p.expr()
		p.consume(&TokenCloseBrace{})
		return expr
	case *TokenIdent:
		ident := p.curr().(*TokenIdent)
		p.next()
		return &Name{
			Value: ident.value,
		}
	default:
		panic("unknown node in atom()")
	}
}

func (p *Parser) program() *Program {
	program := new(Program)

	for reflect.TypeOf(p.curr()) != reflect.TypeOf(&TokenEOF{}) {
		program.StmtList = append(program.StmtList, p.stmt())
	}

	program.EOF = p.curr().Pos()

	return program
}

func (p *Parser) valDeclStmt() *ValDeclStmt {
	var valDeclStmt = &ValDeclStmt{}
	p.consume(&TokenVal{})
	if !p.match(&TokenIdent{}) {
		panic(fmt.Errorf("expected TokenIdent, but got %s", tokToString(p.curr())))
	}
	tokenIdent := p.curr().(*TokenIdent)
	valDeclStmt.Name = Name{Value: tokenIdent.value}
	p.next()
	p.consume(&TokenAssign{})
	valDeclStmt.Rhs = p.expr()

	return valDeclStmt
}

func (p *Parser) varDeclStmt() *VarDeclStmt {
	var varDeclStmt = &VarDeclStmt{}
	p.consume(&TokenVar{})
	if !p.match(&TokenIdent{}) {
		panic(fmt.Errorf("expected TokenIdent, but got %s", tokToString(p.curr())))
	}
	tokenIdent := p.curr().(*TokenIdent)
	varDeclStmt.Name = Name{Value: tokenIdent.value}
	p.next()
	p.consume(&TokenAssign{})
	varDeclStmt.Rhs = p.expr()

	return varDeclStmt
}

func (p *Parser) whileStmt() *WhileStmt {
	var whileStmt = &WhileStmt{}
	p.consume(&TokenWhile{})
	p.consume(&TokenOpenParen{})
	whileStmt.Cond = p.cond()
	p.consume(&TokenCloseParen{})
	whileStmt.Body = p.blockStmt()

	return whileStmt
}

func (p *Parser) ifStmt() *IfStmt {
	var ifStmt = &IfStmt{}
	p.consume(&TokenIf{})
	p.consume(&TokenOpenParen{})
	ifStmt.Cond = p.cond()
	p.consume(&TokenCloseParen{})
	ifStmt.Body = p.blockStmt()
	if p.match(&TokenElse{}) {
		p.consume(&TokenElse{})
		ifStmt.ElseBody = p.stmt()
	}

	return ifStmt
}

func (p *Parser) cond() Operation {
	var condition = Operation{}
	condition.Lhs = p.expr()
	operator := tokenToOperator(p.curr())
	if !IsComparisonOp(operator) {
		panic(fmt.Errorf("expected operator, got %s", tokToString(p.curr())))
	}
	condition.Op = operator
	p.next()
	condition.Rhs = p.expr()

	return condition
}

func (p *Parser) assignment() *Assignment {
	var assignment = &Assignment{}
	if !p.match(&TokenIdent{}) {
		errorPos := p.curr().Pos()
		panic(fmt.Errorf("[%d:%d] expected TokenIdent, but got %s", errorPos.Line, errorPos.Column, tokToString(p.curr())))
	}
	ident := p.curr().(*TokenIdent)
	p.next()
	assignment.Lhs = &Name{Value: ident.value}
	p.consume(&TokenAssign{})
	assignment.Rhs = p.expr()

	return assignment
}

func (p *Parser) blockStmt() *BlockStmt {
	block := new(BlockStmt)
	block.pos = p.curr().Pos()

	// constrain block statement to actually contain statements, i.e
	// do not allow empty blocks
	p.consume(&TokenOpenBrace{})
	block.Stmts = p.stmts()
	p.consume(&TokenCloseBrace{})

	return block
}

func (p *Parser) stmts() []Stmt {
	var stmtList []Stmt
	for (reflect.TypeOf(p.curr()) != reflect.TypeOf(&TokenEOF{})) &&
		(reflect.TypeOf(p.curr()) != reflect.TypeOf(&TokenCloseBrace{})) {

		stmtList = append(stmtList, p.stmt())
	}

	return stmtList
}

func (p *Parser) expr() Node {
	switch p.curr().(type) {
	case *TokenNumber, *TokenOpenParen, *TokenOpenBrace, *TokenIdent:
		return p.binOp(0)
	default:
		panic(fmt.Errorf("expected TokenNumber, TokenOpenParen, TokenOpenBrace, TokenIdent, TokenVal, TokenVar, but got %s", tokToString(p.curr())))
	}
}

func Parse(path string) *Program {
	stream, err := os.Open(path)
	if err != nil {
		panic("no file with such path was found")
	}

	scanner := newCharScanner(stream)
	tokens := scanner.Tokenize()

	var parser = &Parser{
		tokenStream: tokens,
		currIdx:     0,
	}

	return parser.program()
}

func prec(token Token) int {
	assert.Assert(IsOperator(token), "should be an operator")
	switch token.(type) {
	default:
		return 2
	case *TokenPlus, *TokenMinus:
		return 1
	}
}

func assoc(token Token) AssociativityType {
	assert.Assert(IsOperator(token), "should be an operator")
	switch token.(type) {
	default:
		return RightAssociative
	case *TokenPlus, *TokenMinus, *TokenMul, *TokenDiv:
		return LeftAssociative
	}
}

func IsOperator(token Token) bool {
	switch token.(type) {
	default:
		return false
	// TODO(threadedstream): add comparison operators
	case *TokenPlus, *TokenMinus, *TokenMul, *TokenDiv:
		return true
	}
}

func tokenToOperator(token Token) Operator {
	switch token.(type) {
	case *TokenPlus:
		return Plus
	case *TokenMinus:
		return Minus
	case *TokenMul:
		return Mul
	case *TokenGreaterThan:
		return GreaterThan
	case *TokenGreaterThanOrEqual:
		return GreaterThanOrEqual
	case *TokenLessThan:
		return LessThan
	case *TokenLessThanOrEqual:
		return LessThanOrEqual
	case *TokenEqual:
		return Equal
	case *TokenNotEqual:
		return NotEqual
	default:
		return InvalidOperator
	}
}

func IsComparisonOp(op Operator) bool {
	switch op {
	default:
		return false
	case GreaterThan, GreaterThanOrEqual, LessThan, LessThanOrEqual, Equal, NotEqual:
		return true
	}
}