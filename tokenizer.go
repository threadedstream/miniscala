package main

import (
	"fmt"
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

	TokenScanner struct {
		TokenReader
		peekToken Token
		in        *CharScanner
	}
)

func newCharScanner(reader io.Reader) *CharScanner {
	var charScanner = new(CharScanner)
	charScanner.s = charScanner.s.Init(reader)
	return charScanner
}

func newTokenScanner(reader io.Reader) *TokenScanner {
	var tokenScanner = new(TokenScanner)
	tokenScanner.in = newCharScanner(reader)
	tokenScanner.next()
	return tokenScanner
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

func (cs *CharScanner) next() rune {
	return cs.s.Next()
}

func (ts *TokenScanner) hasNext() bool {
	return ts.peek() != EOF{}
}

func (ts *TokenScanner) hasNextP(f func(Token) bool) bool {
	return f(ts.peek())
}

func (ts *TokenScanner) next() Token {
	temp := ts.getToken()

	if _, ok := temp.(EOF); !ok {
		ts.peekToken = temp
	}

	return temp
}

func (ts *TokenScanner) skipWhitespaces() {
	for ts.in.hasNextP(unicode.IsSpace) {
		ts.in.next()
	}
}

func (ts *TokenScanner) getNum() Token {
	if ts.in.hasNextP(unicode.IsDigit) {
		n := 0
		for ts.in.hasNextP(unicode.IsDigit) {
			n = 10*n + (int)(ts.in.next()-'0')
		}
		return Number{x: n}
	} else {
		panic("expected number")
	}
}

func (ts *TokenScanner) getRawToken() Token {
	if ts.in.hasNextP(unicode.IsLetter) {
		return ts.getName()
	} else if ts.in.hasNextP(isOperator) {
		return ts.getOperator()
	} else if ts.in.hasNextP(unicode.IsDigit) {
		return ts.getNum()
	} else if ts.in.hasNextP(isDelim) {
		return Delim{x: ts.in.next()}
	} else if !ts.in.hasNext() {
		return EOF{}
	} else {
		panic(fmt.Errorf("unexpected character %c", ts.in.peek()))
	}
}

func (ts *TokenScanner) getToken() Token {
	ts.skipWhitespaces()
	return ts.getRawToken()
}

func (ts *TokenScanner) getName() Token {
	var buf []rune
	for ts.in.hasNextP(isAlphaNum) {
		buf = append(buf, ts.in.next())
	}
	s := string(buf)

	if isKeyword(s) {
		return Keyword{x: s}
	} else {
		return Ident{x: s}
	}
}

func (ts *TokenScanner) getOperator() Token {
	if ts.in.hasNextP(isOperator) {
		return Ident{x: string(ts.in.next())}
	} else {
		panic("expected operator")
	}
}

func isKeyword(kwd string) bool {
	return kwd == "val"
}

func isDelim(c rune) bool {
	return c == '{' || c == '}' || c == '(' || c == ')' || c == ';' || c == '='
}

func isOperator(c rune) bool {
	return c == '+' || c == '-' || c == '*' || c == '/' || c == '%'
}

func isAlphaNum(c rune) bool {
	return unicode.IsLetter(c) || unicode.IsDigit(c)
}
