package vm

import (
	"github.com/ThreadedStream/miniscala/backing"
	"github.com/ThreadedStream/miniscala/syntax"
)

var (
	code []Instruction
)

func Compile(program *syntax.Program) {
	for _, stmt := range program.StmtList {
		compileStmt(stmt)
	}
}

func compileStmt(stmt syntax.Stmt) {
	switch stmt.(type) {
	default:
		compileExpr(stmt)
	case *syntax.BlockStmt:
		compileBlockStmt(stmt)
	case *syntax.IfStmt:
		compileIfStmt(stmt)
	}
}

func compileBlockStmt(stmt syntax.Stmt) {
	block := stmt.(*syntax.BlockStmt)

	for _, currStmt := range block.Stmts {
		compileStmt(currStmt)
	}
}

func compileIfStmt(stmt syntax.Stmt) {
	ifStmt := stmt.(*syntax.IfStmt)

	compileExpr(&ifStmt.Cond)
	jmpIfFalseInstr := &InstrJmpIfFalse{}
	code = append(code, jmpIfFalseInstr)
	priorCodeLen := len(code)
	compileBlockStmt(ifStmt.Body)
	posteriorCodeLen := len(code)
	jmpIfFalseInstr.Offset = posteriorCodeLen - priorCodeLen
}

func compileExpr(expr syntax.Expr) {
	switch expr.(type) {
	case *syntax.BasicLit:
		compileBasicLit(expr)
	}
}

func compileBasicLit(expr syntax.Expr) {
	basicLit := expr.(*syntax.BasicLit)
	// TODO(threadedstream): function for mapping kinds to value types is needed
	_ = &InstrLoadImm{
		Value: backing.Value{
			Value: basicLit.Value,
		},
	}
	//code = append(code, )
}
