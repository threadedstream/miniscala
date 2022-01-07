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
		peekChar  rune
		peekChar1 rune
		buffer    string
		s         *scanner.Scanner
	}

	TokenScanner struct {
		TokenReader
		peekToken  Token
		peekToken1 Token
		in         *CharScanner
		position   scanner.Position
	}
)

func newCharScanner(reader io.Reader) *CharScanner {
	var charScanner = &CharScanner{
		s: new(scanner.Scanner),
	}
	charScanner.s = charScanner.s.Init(reader)
	// heap allocated temporary buffer
	var temp = make([]byte, 128)
	_, err := reader.Read(temp)
	if err != nil {
		panic("newCharScanner() failed to read into buffer")
	}
	charScanner.buffer = string(temp)
	// trigger gc
	temp = nil
	charScanner.peekChar = charScanner.s.Next()
	charScanner.peekChar1 = charScanner.s.Next()
	return charScanner
}

func newTokenScanner(reader io.Reader) *TokenScanner {
	var tokenScanner = new(TokenScanner)
	tokenScanner.in = newCharScanner(reader)
	tokenScanner.peekToken = tokenScanner.getToken()
	tokenScanner.peekToken1 = tokenScanner.getToken()
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

func (cs *CharScanner) hasNextP2(f func(rune, rune) bool) bool {
	return f(cs.peekChar, cs.peekChar1)
}

func (cs *CharScanner) next() rune {
	cs.peekChar = cs.peekChar1
	cs.peekChar1 = cs.s.Next()
	return cs.peekChar
}

func (ts *TokenScanner) hasNext() bool {
	return ts.peekToken != EOF{}
}

func (ts *TokenScanner) hasNextP(f func(Token) bool) bool {
	return f(ts.peekToken)
}

func (ts *TokenScanner) hasNextP2(f func(Token, Token) bool) bool {
	return f(ts.peekToken, ts.peekToken1)
}

func (ts *TokenScanner) next() Token {
	temp := ts.getToken()
	ts.peekToken = temp
	return temp
}

func (ts *TokenScanner) pos() scanner.Position {
	return ts.in.s.Pos()
}

func isCommentStart(c1, c2 rune) bool {
	return c1 == '/' && c2 == '/'
}

func (ts *TokenScanner) skipWhitespaces() {
	for ts.in.hasNextP(unicode.IsSpace) || ts.in.hasNextP2(isCommentStart) {
		if ts.in.peekChar == '/' {
			ts.in.next()
			for ts.in.peekChar != '\n' && ts.in.peekChar != scanner.EOF {
				ts.in.next()
			}
		}
	}
}

func (ts *TokenScanner) getNum() Number {
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

//func pointToSource(si SourceInfo) {
//	var width = 0
//	if si.end.Line == si.start.Line {
//		width = si.end.Column - si.start.Column
//	} else {
//		width = 0
//	}
//	beginning := si.start.Line * si.start.Column
//	substr :=
//}

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
		errPos := ts.pos()
		panic(fmt.Errorf("unexpected character %c at %d:%d\n", ts.in.peek(), errPos.Line, errPos.Column))
	}
}

func (ts *TokenScanner) getToken() Token {
	var gap = ts.pos()
	ts.skipWhitespaces()
	var (
		start = ts.pos()
		token = ts.getRawToken()
		end   = ts.pos()
	)

	// TODO(threadedstream): Dear stranger, i sincerely apologize
	// for the shocking feeling you're having right now just by looking
	// at the horrendous piece of code taking place right below this comment.
	// Don't get me wrong, i'm just trying to translate code written in Scala to Go.
	// Perhaps, there's a much better alternative to this, but i'm no expert and
	// , left with no choice, doing it that way

	switch token.(type) {
	case Ident:
		ident := token.(Ident)
		ident.position = SourceInfo{gap: gap, start: start, end: end}
	case Keyword:
		keyword := token.(Keyword)
		keyword.position = SourceInfo{gap: gap, start: start, end: end}
	case Number:
		number := token.(Number)
		number.position = SourceInfo{gap: gap, start: start, end: end}
	case Delim:
		delim := token.(Delim)
		delim.position = SourceInfo{gap: gap, start: start, end: end}
	}

	return token
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

func (ts *TokenScanner) getOperator() Ident {
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
