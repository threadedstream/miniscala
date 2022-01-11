package main

import "text/scanner"

type Operator int

const (
	PlusOp             Operator = iota // +
	MinusOp                            // -
	MulOp                              // *
	GreaterThan                        // >
	GreaterThanOrEqual                 // >=
	LessThan                           // <
	LessThanOrEqual                    // <=
	Equal                              // ==
	NotEqual                           // !=
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
	stmtList []Stmt
	EOF      scanner.Position
	node
}

type (
	Expr interface {
		Node
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
	Stmt interface {
		Node
	}

	BlockStmt struct {
		stmts []Stmt
		stmt
	}

	WhileStmt struct {
		cond Operation
		body *BlockStmt
		stmt
	}

	IfStmt struct {
		cond     Operation
		body     *BlockStmt
		elseBody Stmt
		stmt
	}

	Assignment struct {
		lhs Expr
		rhs Expr
		stmt
	}

	VarDeclStmt struct {
		name Name
		rhs  Expr
		stmt
	}

	ValDeclStmt struct {
		name Name
		rhs  Expr
		stmt
	}

	DefDeclStmt struct {
		name *Name
		body *BlockStmt
	}
)

type stmt struct {
	node
}

type LitKind uint8

const (
	StringLit LitKind = iota
	FloatLit
)
