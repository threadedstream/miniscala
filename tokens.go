package main

import (
	"fmt"
	"text/scanner"
)
import "C"

type (
	SourceInfo struct {
		file  string
		gap   scanner.Position
		start scanner.Position
		end   scanner.Position
	}
)

type (
	Token interface {
		pos() SourceInfo
		str() string
	}

	EOF struct {
		Token
	}

	Number struct {
		Token
		x        int
		position SourceInfo
	}

	Ident struct {
		Token
		x        string
		position SourceInfo
	}

	Keyword struct {
		Token
		x        string
		position SourceInfo
	}

	Delim struct {
		Token
		x        rune
		position SourceInfo
	}
)

// TODO(threadedstream): what should i replace it with?

func (n Number) str() string {
	return fmt.Sprintf("%d", n.x)
}

func (n Number) pos() SourceInfo {
	return n.position
}

func (i Ident) str() string {
	return i.x
}

func (i Ident) pos() SourceInfo {
	return i.position
}

func (k Keyword) str() string {
	return k.x
}

func (k Keyword) pos() SourceInfo {
	return k.position
}

func (d Delim) str() string {
	return string(d.x)
}

func (d Delim) pos() SourceInfo {
	return d.position
}
