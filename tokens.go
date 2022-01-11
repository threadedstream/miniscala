package main

import "text/scanner"

type Token interface {
	Pos() scanner.Position
	Text() string
}

type tok struct {
	text string
	pos  scanner.Position
}

type (
	TokenVar struct {
		tok
	}

	TokenVal struct {
		tok
	}

	TokenSemicolon struct {
		tok
	}

	TokenAssign struct {
		tok
	}

	TokenPlus struct {
		tok
	}

	TokenMinus struct {
		tok
	}

	TokenMul struct {
		tok
	}

	TokenIdent struct {
		name string
		tok
	}

	TokenNumber struct {
		value float64
		tok
	}

	TokenString struct {
		value string
		tok
	}

	TokenWhile struct {
		tok
	}

	TokenIf struct {
		tok
	}
)

//const (
//	TokenVar int = iota
//	TokenVal
//	TokenSemicolon
//	TokenAssign
//	TokenPlus
//	TokenMinus
//	TokenMul
//	TokenIdent
//  TokenNumber
//  TokenString
//)

func (t *tok) Pos() scanner.Position {
	return t.pos
}

func (t *tok) Text() string {
	return t.text
}
