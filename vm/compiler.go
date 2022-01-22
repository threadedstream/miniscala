package vm

import (
	"github.com/ThreadedStream/miniscala/backing"
	"github.com/ThreadedStream/miniscala/syntax"
	"strconv"
)

var (
	code []Instruction
)

type compiler struct {
	code              []Instruction
	hadCompilerErrors bool
	errors            []string
}

func newCompiler() *compiler {
	return new(compiler)
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

func (c *compiler) compile(program *syntax.Program) {
	for _, stmt := range program.StmtList {
		c.compileStmt(stmt)
	}
}

func (c *compiler) compileStmt(stmt syntax.Stmt) {
	switch stmt.(type) {
	default:
		c.compileExpr(stmt)
	case *syntax.BlockStmt:
		c.compileBlockStmt(stmt)
	case *syntax.IfStmt:
		c.compileIfStmt(stmt)
	case *syntax.ReturnStmt:
		c.compileReturnStmt(stmt)
	case *syntax.DefDeclStmt:
		c.compileDefDeclStmt(stmt)
	}
}

func (c *compiler) compileBlockStmt(stmt syntax.Stmt) {
	block := stmt.(*syntax.BlockStmt)
	for _, currStmt := range block.Stmts {
		c.compileStmt(currStmt)
	}
}

func (c *compiler) compileIfStmt(stmt syntax.Stmt) {
	ifStmt := stmt.(*syntax.IfStmt)
	c.compileExpr(&ifStmt.Cond)
	jmpIfFalseInstr := &InstrJmpIfFalse{}
	code = append(code, jmpIfFalseInstr)
	priorCodeLen := len(code)
	c.compileBlockStmt(ifStmt.Body)
	posteriorCodeLen := len(code)
	jmpIfFalseInstr.Offset = posteriorCodeLen - priorCodeLen
}

func (c *compiler) compileDefDeclStmt(stmt syntax.Stmt) {
	defStmt := stmt.(*syntax.DefDeclStmt)
	chunk := newChunk(nil, defStmt.Name.Value)
	for _, param := range defStmt.ParamList {
		chunk.localValues[param.Name.Value] = backing.NullValue()
	}

	chunkStore[defStmt.Name.Value] = chunk
	c.compileBlockStmt(defStmt.Body)
	chunk.instrStream = code
}

func (c *compiler) compileExpr(expr syntax.Expr) {
	switch expr.(type) {
	default:
		// TODO(threadedstream): get out of habit panicking right away
		panic("unknown expression ")
	case *syntax.BasicLit:
		c.compileBasicLit(expr)
	case *syntax.Name:
		c.compileName(expr)
	case *syntax.Operation:
		c.compileOperation(expr)
	case *syntax.Call:
		c.compileCall(expr)
	}
}

func (c *compiler) compileBasicLit(expr syntax.Expr) {
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

func (c *compiler) compileName(expr syntax.Expr) {
	name := expr.(*syntax.Name)
	loadRefInstr := &InstrLoadRef{
		RefName: name.Value,
	}
	code = append(code, loadRefInstr)
}

func (c *compiler) compileCall(expr syntax.Expr) {
	call := expr.(*syntax.Call)
	chunk := lookupChunk(call.CalleeName.Value, true, nil)
	callInstr := &InstrCall{
		FuncName: call.CalleeName.Value,
	}

	for _, arg := range call.ArgList {
		c.compileExpr(arg)
	}

	i := 0
	for k, _ := range chunk.localValues {
		callInstr.ArgNames[i] = k
	}

	code = append(code, callInstr)
}

func (c *compiler) compileOperation(expr syntax.Expr) {
	operation := expr.(*syntax.Operation)
	c.compileExpr(operation.Lhs)
	c.compileExpr(operation.Rhs)
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

func (c *compiler) compileReturnStmt(stmt syntax.Stmt) {
	returnStmt := stmt.(*syntax.ReturnStmt)
	c.compileExpr(returnStmt.Value)
	code = append(code, &InstrReturn{})
}
