package main

import "text/scanner"

type Operator int

const (
	PlusOp Operator = iota
	MinusOp
	MulOp
)

type SourceInfo struct {
	start scanner.Position
	end   scanner.Position
}

type Node interface {
	Pos() SourceInfo
}

type node struct {
	pos SourceInfo
}

func (n *node) Pos() SourceInfo {
	return n.pos
}

type Program struct {
	nodeList []Node
	EOF      scanner.Position
	node
}

type (
	VarDecl struct {
		name Name
		rhs  Expr
		node
	}

	ValDecl struct {
		name Name
		rhs  Expr
		node
	}
)

type (
	Expr interface {
		Node
	}

	Assignment struct {
		lhs Expr
		rhs Expr
		expr
	}

	Operation struct {
		op  Operator
		lhs Expr
		rhs Expr
		expr
	}

	BasicLit struct {
		value string
		kind  LitKind
		expr
	}

	Name struct {
		value string
		expr
	}
)

type expr struct {
	node
}

type (
	Statement interface {
		Node
	}

	WhileStmt struct {
		cond Operation
		body Expr
		statement
	}

	IfStmt struct {
		cond     Operation
		body     Expr
		elseBody Expr
		statement
	}
)

type statement struct {
	node
}

type LitKind uint8

const (
	StringLit LitKind = iota
	FloatLit
)
