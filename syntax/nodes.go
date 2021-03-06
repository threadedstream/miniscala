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
	Mod                                // %
	GreaterThan                        // >
	GreaterThanOrEqual                 // >=
	LessThan                           // <
	LessThanOrEqual                    // <=
	Equal                              // ==
	NotEqual                           // !=
	LogicalAnd                         // &&
	LogicalOr                          // ||
	LogicalNot                         // !
	InvalidOperator    = -1
)

func OperatorToString(op Operator) string {
	switch op {
	default:
		return "?"
	case Plus:
		return "+"
	case Minus:
		return "-"
	case Mul:
		return "*"
	case Div:
		return "/"
	case Mod:
		return "%"
	case GreaterThan:
		return ">"
	case GreaterThanOrEqual:
		return ">="
	case LessThan:
		return "<"
	case LessThanOrEqual:
		return "<="
	case Equal:
		return "=="
	case NotEqual:
		return "!="
	case LogicalAnd:
		return "&&"
	case LogicalOr:
		return "||"
	case LogicalNot:
		return "!"
	}

}

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

	ErrExpr struct {
		expr
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
		expr
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

	ErrStmt struct {
		stmt
	}

	// { Stmts }
	BlockStmt struct {
		Stmts []Stmt
		stmt
	}

	// while (Cond) { Body }
	WhileStmt struct {
		Cond *Operation
		Body *BlockStmt
		stmt
	}

	Call struct {
		CalleeName *Name
		ArgList    []Expr
		stmt
	}

	// if (Cond) { Body } else ElseBody
	IfStmt struct {
		Cond     *Operation
		Body     *BlockStmt
		ElseBody Stmt
		stmt
	}

	ReturnStmt struct {
		Value Expr
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
		Name       *Name
		ParamList  []*Field
		ReturnType Expr
		Body       *BlockStmt
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
	IntLit
	BoolLit
)
