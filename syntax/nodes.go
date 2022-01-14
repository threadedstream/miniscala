package syntax

import (
	"text/scanner"
)

type Operator int

const (
	Plus               Operator = iota // +
	Minus                              // -
	Mul                                // *
	Div                                // /
	GreaterThan                        // >
	GreaterThanOrEqual                 // >=
	LessThan                           // <
	LessThanOrEqual                    // <=
	Equal                              // ==
	NotEqual                           // !=
	InvalidOperator    = -1
)

type Node interface {
	Pos() scanner.Position
}

type node struct {
	pos scanner.Position
}

func (n *node) Pos() scanner.Position {
	return n.pos
}

type Program struct {
	StmtList []Stmt
	EOF      scanner.Position
	node
}

type (
	Expr interface {
		Node
	}

	Operation struct {
		Op  Operator
		Lhs Expr
		Rhs Expr
		expr
	}

	BasicLit struct {
		Value string
		Kind  LitKind
		expr
	}

	// Name: Type
	Field struct {
		Name *Name
		Type Expr
	}

	Name struct {
		Value string
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

	// { Stmts }
	BlockStmt struct {
		Stmts []Stmt
		stmt
	}

	// while (Cond) { Body }
	WhileStmt struct {
		Cond Operation
		Body *BlockStmt
		stmt
	}

	// if (Cond) { Body } else ElseBody
	IfStmt struct {
		Cond     Operation
		Body     *BlockStmt
		ElseBody Stmt
		stmt
	}

	// Lhs = Rhs
	Assignment struct {
		Lhs Expr
		Rhs Expr
		stmt
	}

	// var Name = Rhs
	VarDeclStmt struct {
		Name Name
		Rhs  Expr
		stmt
	}

	// val Name = Rhs
	ValDeclStmt struct {
		Name Name
		Rhs  Expr
		stmt
	}

	// def Name ( ParamList ) :Type { Body }
	DefDeclStmt struct {
		Name      *Name
		ParamList []*Field
		Type      Expr
		Body      *BlockStmt
		stmt
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
