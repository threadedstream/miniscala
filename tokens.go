package main

import (
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

//func (n *Number) setPos(si SourceInfo) {
//	n.position = si
//}
//
func (n Number) pos() SourceInfo {
	return n.position
}

//
//func (i *Ident) setPos(si SourceInfo) {
//	i.position = si
//}
//
func (i Ident) pos() SourceInfo {
	return i.position
}

//
//func (k *Keyword) setPos(si SourceInfo) {
//	k.position = si
//}
//
func (k Keyword) pos() SourceInfo {
	return k.position
}

//
//func (d *Delim) setPos(si SourceInfo) {
//	d.position = si
//}
//
func (d Delim) pos() SourceInfo {
	return d.position
}
