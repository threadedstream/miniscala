package main

import "text/scanner"

type Token interface {
	Pos() scanner.Position
}

type tok struct {
	pos scanner.Position
}

type (
	TokenVar struct {
		tok
	}

	TokenVal struct {
		tok
	}

	TokenDef struct {
		tok
	}

	TokenSemicolon struct {
		tok
	}

	TokenAssign struct {
		tok
	}

	TokenEqual struct {
		tok
	}

	TokenNotEqual struct {
		tok
	}

	TokenGreaterThan struct {
		tok
	}

	TokenGreaterThanOrEqual struct {
		tok
	}

	TokenLessThan struct {
		tok
	}

	TokenLessThanOrEqual struct {
		tok
	}

	TokenOpenBrace struct {
		tok
	}

	TokenCloseBrace struct {
		tok
	}

	TokenOpenParen struct {
		tok
	}

	TokenCloseParen struct {
		tok
	}

	TokenEOF struct {
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
		value string
		tok
	}

	TokenNumber struct {
		value string
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

	TokenUnknown struct {
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

func tokToString(token Token) string {
	switch token.(type) {
	case *TokenVar:
		return "TokenVar"
	case *TokenVal:
		return "TokenVal"
	case *TokenDef:
		return "TokenDef"
	case *TokenSemicolon:
		return "TokenSemicolon"
	case *TokenAssign:
		return "TokenAssign"
	case *TokenEqual:
		return "TokenEqual"
	case *TokenNotEqual:
		return "TokenNotEqual"
	case *TokenGreaterThan:
		return "TokenGreaterThan"
	case *TokenGreaterThanOrEqual:
		return "TokenGreaterThanOrEqual"
	case *TokenLessThan:
		return "TokenLessThan"
	case *TokenLessThanOrEqual:
		return "TokenLessThanOrEqual"
	case *TokenPlus:
		return "TokenPlus"
	case *TokenMinus:
		return "TokenMinus"
	case *TokenMul:
		return "TokenMul"
	case *TokenString:
		return "TokenString"
	case *TokenIf:
		return "TokenIf"
	case *TokenWhile:
		return "TokenWhile"
	case *TokenIdent:
		return "TokenIdent"
	case *TokenOpenBrace:
		return "TokenOpenBrace"
	case *TokenCloseBrace:
		return "TokenCloseBrace"
	case *TokenOpenParen:
		return "TokenOpenParen"
	case *TokenCloseParen:
		return "TokenCloseParen"
	case *TokenEOF:
		return "TokenEOF"
	default:
		return "TokenUnknown"
	}
}
