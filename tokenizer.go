package main

import (
	"fmt"
	"text/scanner"
	"unicode"
)

type Reader struct {
	scanner *scanner.Scanner
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

func (r *Reader) hasNext() bool {
	return r.scanner.Peek() != scanner.EOF
}

func (r *Reader) hasNextP(predicate func(rune) bool) bool {
	return predicate(r.scanner.Peek())
}

func (r *Reader) skipWhitespaces() {
	for r.hasNextP(unicode.IsSpace) {
		r.scanner.Next()
	}
}

func (r *Reader) getNum() Token {
	if r.hasNextP(unicode.IsDigit) {
		n := 0
		for r.hasNextP(unicode.IsDigit) {
			n = 10*n + (int)(r.scanner.Next()-'0')
		}
		return Number{x: n}
	} else {
		panic("expected number")
	}
}

func (r *Reader) getRawToken() Token {
	if r.hasNextP(unicode.IsLetter) {
		return r.getName()
	} else if r.hasNextP(isOperator) {
		return r.getOperator()
	} else if r.hasNextP(unicode.IsDigit) {
		return r.getNum()
	} else if r.hasNextP(isDelim) {
		return Delim{x: r.scanner.Next()}
	} else if !r.hasNext() {
		return EOF{}
	} else {
		panic(fmt.Errorf("unexpected character %c", r.scanner.Peek()))
	}
}

func (r *Reader) getToken() Token {
	r.skipWhitespaces()
	return r.getRawToken()
}

func (r *Reader) getName() Token {
	var buf []rune
	for r.hasNextP(isAlphaNum) {
		buf = append(buf, r.scanner.Next())
	}
	s := string(buf)

	if isKeyword(s) {
		return Keyword{x: s}
	} else {
		return Ident{x: s}
	}
}

func (r *Reader) getOperator() Token {
	if r.hasNextP(isOperator) {
		return Ident{x: string(r.scanner.Next())}
	} else {
		panic("expected operator")
	}
}
