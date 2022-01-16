package syntax

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

	TokenComma struct {
		tok
	}

	TokenAssign struct {
		tok
	}

	TokenReturn struct {
		tok
	}

	TokenEqual struct {
		tok
	}

	TokenNotEqual struct {
		tok
	}

	TokenColon struct {
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

	TokenPlus struct {
		tok
	}

	TokenMinus struct {
		tok
	}

	TokenMul struct {
		tok
	}

	TokenDiv struct {
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

	TokenElse struct {
		tok
	}

	TokenEOF struct {
		tok
	}

	TokenUnknown struct {
		tok
	}
)

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
	case *TokenColon:
		return "TokenColon"
	case *TokenPlus:
		return "TokenPlus"
	case *TokenMinus:
		return "TokenMinus"
	case *TokenMul:
		return "TokenMul"
	case *TokenDiv:
		return "TokenDiv"
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
	case *TokenReturn:
		return "TokenReturn"
	case *TokenEOF:
		return "TokenEOF"
	default:
		return "TokenUnknown"
	}
}
