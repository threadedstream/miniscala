package vm

import (
	"github.com/ThreadedStream/miniscala/backing"
	"github.com/ThreadedStream/miniscala/syntax"
	"strconv"
)

var (
	code []Instruction
)

type Compiler struct {
	Code              []Instruction
	hadCompilerErrors bool
	errors            []string
}

func litKindToValueType(kind syntax.LitKind) backing.ValueType {
	switch kind {
	default:
		return backing.Undefined
	case syntax.StringLit:
		return backing.String
	case syntax.FloatLit:
		return backing.Float
	case syntax.IntLit:
		return backing.Int
	case syntax.BoolLit:
		return backing.Bool
	}
}

func Compile(program *syntax.Program) []Instruction {
	for _, stmt := range program.StmtList {
		compileStmt(stmt)
	}

	return code
}

func compileStmt(stmt syntax.Stmt) {
	switch stmt.(type) {
	default:
		compileExpr(stmt)
	case *syntax.BlockStmt:
		compileBlockStmt(stmt)
	case *syntax.IfStmt:
		compileIfStmt(stmt)
	case *syntax.ReturnStmt:
		compileReturnStmt(stmt)
	case *syntax.DefDeclStmt:
		compileDefDeclStmt(stmt)
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

func compileDefDeclStmt(stmt syntax.Stmt) {
	defStmt := stmt.(*syntax.DefDeclStmt)
	chunk := Chunk{
		funcName: defStmt.Name.Value,
	}
	compileBlockStmt(defStmt.Body)
	chunk.instrStream = code
}

func compileExpr(expr syntax.Expr) {
	switch expr.(type) {
	default:
		// TODO(threadedstream): get out of habit panicking right away
		panic("unknown expression ")
	case *syntax.BasicLit:
		compileBasicLit(expr)
	case *syntax.Name:
		compileName(expr)
	case *syntax.Operation:
		compileOperation(expr)
	case *syntax.Call:
		compileCall(expr)
	}
}

func compileBasicLit(expr syntax.Expr) {
	basicLit := expr.(*syntax.BasicLit)
	loadInstr := &InstrLoadImm{}
	var value backing.Value
	switch basicLit.Kind {
	case syntax.StringLit:
		value.Value = basicLit.Value
		value.ValueType = backing.String
	case syntax.FloatLit:
		value.Value, _ = strconv.ParseFloat(basicLit.Value, 64)
		value.ValueType = backing.Float
	case syntax.IntLit:
		value.Value, _ = strconv.ParseInt(basicLit.Value, 10, 64)
		value.ValueType = backing.Int
	case syntax.BoolLit:
		value.Value, _ = strconv.ParseBool(basicLit.Value)
		value.ValueType = backing.Bool
	}

	loadInstr.Value = value
	code = append(code, loadInstr)
}

func compileName(expr syntax.Expr) {
	name := expr.(*syntax.Name)
	loadRefInstr := &InstrLoadRef{
		RefName: name.Value,
	}
	code = append(code, loadRefInstr)
}

func compileCall(expr syntax.Expr) {

}

func compileOperation(expr syntax.Expr) {
	operation := expr.(*syntax.Operation)
	compileExpr(operation.Lhs)
	compileExpr(operation.Rhs)
	switch operation.Op {
	default:
		// TODO(threadedstream): handle an error
		panic("")
	case syntax.Plus:
		code = append(code, &InstrAdd{})
	case syntax.Minus:
		code = append(code, &InstrSub{})
	case syntax.Mul:
		code = append(code, &InstrMul{})
	case syntax.Div:
		code = append(code, &InstrDiv{})
	case syntax.GreaterThan:
		code = append(code, &InstrGreaterThan{})
	case syntax.GreaterThanOrEqual:
		code = append(code, &InstrGreaterThanOrEqual{})
	case syntax.LessThan:
		code = append(code, &InstrLessThan{})
	case syntax.LessThanOrEqual:
		code = append(code, &InstrLessThanOrEqual{})
	case syntax.Equal:
		code = append(code, &InstrEqual{})
	}
}

func compileReturnStmt(stmt syntax.Stmt) {
	returnStmt := stmt.(*syntax.ReturnStmt)
	compileExpr(returnStmt.Value)
	code = append(code, &InstrReturn{})
}
