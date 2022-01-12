package main

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
		tokens = append(tokens, token)
	}

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
				}
			}
		}
		return &TokenUnknown{}
	case '=':
		cs.s.Next()
		if cs.s.Peek() == '=' {
			cs.s.Next()
			return &TokenEqual{}
		} else {
			return &TokenAssign{}
		}
	case '{':
		cs.s.Next()
		return &TokenOpenBrace{}
	case '}':
		cs.s.Next()
		return &TokenCloseBrace{}
	case '(':
		cs.s.Next()
		return &TokenOpenParen{}
	case ')':
		cs.s.Next()
		return &TokenCloseParen{}
	case '+':
		cs.s.Next()
		return &TokenPlus{}
	case '-':
		cs.s.Next()
		return &TokenMinus{}
	case '*':
		cs.s.Next()
		return &TokenMul{}
	case '>':
		cs.s.Next()
		if cs.s.Peek() == '=' {
			cs.s.Next()
			return &TokenGreaterThanOrEqual{}
		} else {
			return &TokenGreaterThan{}
		}
	case '<':
		cs.s.Next()
		if cs.s.Peek() == '=' {
			cs.s.Next()
			return &TokenLessThanOrEqual{}
		} else {
			return &TokenLessThan{}
		}
	case '!':
		cs.s.Next()
		if cs.s.Peek() == '=' {
			return &TokenNotEqual{}
		}
	case '"':
		cs.s.Next()
		return cs.tokenizeString()
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return cs.tokenizeNumber()
	case scanner.EOF:
		return &TokenEOF{}
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

func (cs *CharScanner) tokenizeKeyword(kwd string) Token {
	switch kwd {
	case "val":
		return &TokenVal{}
	case "var":
		return &TokenVar{}
	case "if":
		return &TokenIf{}
	case "while":
		return &TokenWhile{}
	case "def":
		return &TokenDef{}
	default:
		return &TokenUnknown{}
	}
}

func isKeyword(kwd string) bool {
	switch kwd {
	default:
		return false
	case "val", "var", "if", "while", "def":
		return true
	}
}

func isDelim(c rune) bool {
	return c == '{' || c == '}' || c == '(' || c == ')' || c == ';' || c == '='
}

func isOperator(op rune) bool {
	switch op {
	case '+', '-', '*', '/', '%':
		return true
	default:
		return false
	}
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
