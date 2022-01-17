package syntax

import (
	"io"
	"text/scanner"
	"unicode"
)

type (
	Reader interface {
		hasNext() bool
	}

	TokenReader interface {
		Reader
		peek() Token
		hasNextP(func(Token) bool)
		next() Token
	}

	CharReader interface {
		Reader
		peek() rune
		hasNextP(func(rune) bool)
		next() rune
	}

	CharScanner struct {
		CharReader
		s *scanner.Scanner
	}
)

func newCharScanner(reader io.Reader) *CharScanner {
	var charScanner = &CharScanner{
		s: new(scanner.Scanner),
	}
	charScanner.s = charScanner.s.Init(reader)
	return charScanner
}

func (cs *CharScanner) Tokenize() []Token {
	var tokens []Token
	for cs.s.Peek() != scanner.EOF {
		token := cs.tokenize()
		switch token.(type) {
		default:
			tokens = append(tokens, token)
		case *TokenComment:
			break
		}
	}

	// append EOF token
	tokens = append(tokens, &TokenEOF{})

	return tokens
}

func (cs *CharScanner) tokenize() Token {
	for unicode.IsSpace(cs.s.Peek()) {
		cs.s.Next()
	}
	switch cs.s.Peek() {
	default:
		peekChar := cs.s.Peek()
		if isAlphaNum(peekChar) {
			pos := cs.s.Pos()
			var tokenValue []rune
			for isAlphaNum(cs.s.Peek()) {
				tokenValue = append(tokenValue, cs.s.Peek())
				cs.s.Next()
			}
			if isKeyword(string(tokenValue)) {
				return cs.tokenizeKeyword(string(tokenValue))
			} else {
				return &TokenIdent{
					value: string(tokenValue),
					tok: tok{
						pos: pos,
					},
				}
			}
		}
		return &TokenUnknown{}
	case '=':
		pos := cs.s.Pos()
		cs.s.Next()
		if cs.s.Peek() == '=' {
			cs.s.Next()
			return &TokenEqual{
				tok: tok{
					pos: pos,
				},
			}
		} else {
			return &TokenAssign{
				tok: tok{
					pos: pos,
				},
			}
		}
	case '{':
		pos := cs.s.Pos()
		cs.s.Next()
		return &TokenOpenBrace{
			tok: tok{
				pos: pos,
			},
		}
	case '}':
		pos := cs.s.Pos()
		cs.s.Next()
		return &TokenCloseBrace{
			tok: tok{
				pos: pos,
			},
		}
	case '(':
		pos := cs.s.Pos()
		cs.s.Next()
		return &TokenOpenParen{
			tok: tok{
				pos: pos,
			},
		}
	case ')':
		pos := cs.s.Pos()
		cs.s.Next()
		return &TokenCloseParen{
			tok: tok{
				pos: pos,
			},
		}
	case '+':
		pos := cs.s.Pos()
		cs.s.Next()
		return &TokenPlus{
			tok: tok{
				pos: pos,
			},
		}
	case '-':
		pos := cs.s.Pos()
		cs.s.Next()
		return &TokenMinus{
			tok: tok{
				pos: pos,
			},
		}
	case '*':
		pos := cs.s.Pos()
		cs.s.Next()
		return &TokenMul{
			tok: tok{
				pos: pos,
			},
		}
	case '/':
		pos := cs.s.Pos()
		cs.s.Next()
		if cs.s.Peek() == '/' {
			cs.handleComment()
			return &TokenComment{}
		}
		return &TokenDiv{
			tok: tok{
				pos: pos,
			},
		}
	case '>':
		pos := cs.s.Pos()
		cs.s.Next()
		if cs.s.Peek() == '=' {
			cs.s.Next()
			return &TokenGreaterThanOrEqual{
				tok: tok{
					pos: pos,
				},
			}
		} else {
			return &TokenGreaterThan{
				tok: tok{
					pos: pos,
				},
			}
		}
	case '<':
		pos := cs.s.Pos()
		cs.s.Next()
		if cs.s.Peek() == '=' {
			cs.s.Next()
			return &TokenLessThanOrEqual{
				tok: tok{
					pos: pos,
				},
			}
		} else {
			return &TokenLessThan{
				tok: tok{
					pos: pos,
				},
			}
		}
	case '!':
		pos := cs.s.Pos()
		cs.s.Next()
		if cs.s.Peek() == '=' {
			return &TokenNotEqual{
				tok: tok{
					pos: pos,
				},
			}
		}
	case ':':
		pos := cs.s.Pos()
		cs.s.Next()
		return &TokenColon{
			tok: tok{
				pos: pos,
			},
		}
	case ',':
		pos := cs.s.Pos()
		cs.s.Next()
		return &TokenComma{
			tok: tok{
				pos: pos,
			},
		}
	case '"':
		pos := cs.s.Pos()
		cs.s.Next()
		tokenString := cs.tokenizeString()
		tokenString.pos = pos
		return tokenString
	case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
		pos := cs.s.Pos()
		tokenNumber := cs.tokenizeNumber()
		tokenNumber.pos = pos
		return tokenNumber
	case scanner.EOF:
		pos := cs.s.Pos()
		return &TokenEOF{
			tok: tok{
				pos: pos,
			},
		}
	}

	return &TokenUnknown{}
}

func (cs *CharScanner) peek() rune {
	return cs.s.Peek()
}

func (cs *CharScanner) hasNext() bool {
	return cs.s.Peek() != scanner.EOF
}

func (cs *CharScanner) hasNextP(f func(rune) bool) bool {
	return f(cs.peek())
}

func isCommentStart(c1, c2 rune) bool {
	return c1 == '/' && c2 == '/'
}

func (cs *CharScanner) tokenizeNumber() *TokenNumber {
	var digits []rune
	for unicode.IsDigit(cs.s.Peek()) {
		digits = append(digits, cs.s.Peek())
		cs.s.Next()
	}
	return &TokenNumber{value: string(digits)}
}

func (cs *CharScanner) tokenizeString() *TokenString {
	var tokenValue []rune
	for cs.s.Peek() != '"' {
		tokenValue = append(tokenValue, cs.s.Peek())
		cs.s.Next()
	}
	return &TokenString{
		value: string(tokenValue),
	}
}

func (cs *CharScanner) handleComment() {
	for cs.s.Peek() != '\n' && cs.s.Peek() != '\r' {
		cs.s.Next()
	}
}

func (cs *CharScanner) tokenizeKeyword(kwd string) Token {
	switch kwd {
	case "val":
		return &TokenVal{}
	case "var":
		return &TokenVar{}
	case "if":
		return &TokenIf{}
	case "else":
		return &TokenElse{}
	case "while":
		return &TokenWhile{}
	case "def":
		return &TokenDef{}
	case "return":
		return &TokenReturn{}
	default:
		return &TokenUnknown{}
	}
}

func isKeyword(kwd string) bool {
	switch kwd {
	default:
		return false
	case "val", "var", "if", "else", "while", "def", "return":
		return true
	}
}

func isDelim(c rune) bool {
	return c == '{' || c == '}' || c == '(' || c == ')' || c == ';' || c == '='
}

func isCondOp(op string) bool {
	switch op {
	case "==", "!=", "<", ">", "<=", ">=":
		return true
	default:
		return false
	}
}

func isAlphaNum(c rune) bool {
	return unicode.IsLetter(c) || unicode.IsDigit(c)
}
