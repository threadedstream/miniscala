package vm

import (
	"github.com/ThreadedStream/miniscala/backing"
	"github.com/ThreadedStream/miniscala/syntax"
	"strconv"
)

var (
	// TODO(threadedstream): add more reserved functions
	reservedFuncNames = []string{"print"}
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

func (c *compiler) prepareReservedFunctions() {
	for _, funcName := range reservedFuncNames {
		chunk := newChunk(nil, funcName)
		chunkStore[funcName] = chunk
	}
}

func (c *compiler) compile(program *syntax.Program) {
	c.prepareReservedFunctions()
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
	case *syntax.WhileStmt:
		c.compileWhileStmt(stmt)
	case *syntax.ReturnStmt:
		c.compileReturnStmt(stmt)
	case *syntax.DefDeclStmt:
		c.compileDefDeclStmt(stmt)
	case *syntax.Call:
		c.compileCall(stmt)
	case *syntax.Assignment:
		c.compileAssignment(stmt)
	case *syntax.VarDeclStmt:
		c.compileVarDeclStmt(stmt)
	case *syntax.ValDeclStmt:
		c.compileValDeclStmt(stmt)
	}
}

func (c *compiler) compileBlockStmt(stmt syntax.Stmt) {
	block := stmt.(*syntax.BlockStmt)
	for _, currStmt := range block.Stmts {
		c.compileStmt(currStmt)
	}
}

func (c *compiler) compileValDeclStmt(stmt syntax.Stmt) {
	valDeclStmt := stmt.(*syntax.ValDeclStmt)
	c.compileExpr(valDeclStmt.Rhs)
	c.code = append(c.code, &InstrSetLocal{
		Name:       valDeclStmt.Name.Value,
		StoringCtx: backing.Declare,
		Immutable:  true,
	})
}

func (c *compiler) compileVarDeclStmt(stmt syntax.Stmt) {
	varDeclStmt := stmt.(*syntax.VarDeclStmt)
	c.compileExpr(varDeclStmt.Rhs)
	c.code = append(c.code, &InstrSetLocal{
		Name:       varDeclStmt.Name.Value,
		StoringCtx: backing.Declare,
	})
}

func (c *compiler) compileIfStmt(stmt syntax.Stmt) {
	ifStmt := stmt.(*syntax.IfStmt)
	c.compileExpr(&ifStmt.Cond)
	jmpIfFalseInstr := &InstrJmpIfFalse{}
	c.code = append(c.code, jmpIfFalseInstr)
	priorCodeLen := len(c.code)
	c.compileBlockStmt(ifStmt.Body)
	posteriorCodeLen := len(c.code)
	jmpIfFalseInstr.Offset = posteriorCodeLen - priorCodeLen
}

func (c *compiler) compileAssignment(stmt syntax.Stmt) {
	assignment := stmt.(*syntax.Assignment)
	c.compileExpr(assignment.Rhs)
	// dirty little hack, not encouraged, by any means, in industry-strength compilers
	lhs := assignment.Lhs.(*syntax.Name)
	setLocalInstr := &InstrSetLocal{
		Name:       lhs.Value,
		StoringCtx: backing.Assign,
	}
	c.code = append(c.code, setLocalInstr)
}

func (c *compiler) compileWhileStmt(stmt syntax.Stmt) {
	whileStmt := stmt.(*syntax.WhileStmt)
	unCondJmpInit := len(c.code)
	c.compileExpr(&whileStmt.Cond)
	jmpIfFalseInstr := &InstrJmpIfFalse{}
	c.code = append(c.code, jmpIfFalseInstr)
	priorCodeLen := len(c.code)
	c.compileBlockStmt(whileStmt.Body)
	jmpInstr := &InstrJmp{}
	c.code = append(c.code, jmpInstr)
	posteriorCodeLen := len(c.code)
	jmpIfFalseInstr.Offset = posteriorCodeLen - priorCodeLen
	jmpInstr.Offset = unCondJmpInit - posteriorCodeLen
}

func (c *compiler) compileDefDeclStmt(stmt syntax.Stmt) {
	defStmt := stmt.(*syntax.DefDeclStmt)
	chunk := newChunk(nil, defStmt.Name.Value)
	for _, param := range defStmt.ParamList {
		chunk.localValues[param.Name.Value] = backing.NullValue()
	}

	chunk.doesReturn = defStmt.ReturnType.(*syntax.Name).Value != "Unit"
	chunkStore[defStmt.Name.Value] = chunk
	c.compileBlockStmt(defStmt.Body)
	switch c.code[len(c.code)-1].(type) {
	default:
		c.code = append(c.code, &InstrReturn{})
	case *InstrReturn:
		break
	}

	// dirty hack to mutate an element in a map
	chunk = chunkStore[defStmt.Name.Value]

	// allocate a space for an instruction buffer
	chunk.instrStream = make([]Instruction, len(c.code))
	copy(chunk.instrStream, c.code)
	c.code = nil
	c.code = make([]Instruction, 0)
	chunkStore[defStmt.Name.Value] = chunk
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
	c.code = append(c.code, loadInstr)
}

func (c *compiler) compileName(expr syntax.Expr) {
	name := expr.(*syntax.Name)
	loadRefInstr := &InstrLoadRef{
		RefName: name.Value,
	}
	c.code = append(c.code, loadRefInstr)
}

func (c *compiler) compileCall(expr syntax.Expr) {
	var (
		chunk Chunk
		call  = expr.(*syntax.Call)
	)

	chunk = lookupChunk(call.CalleeName.Value, true, nil)

	callInstr := &InstrCall{
		FuncName: call.CalleeName.Value,
	}

	for _, arg := range call.ArgList {
		c.compileExpr(arg)
	}

	for k, _ := range chunk.localValues {
		callInstr.ArgNames = append(callInstr.ArgNames, k)
	}

	c.code = append(c.code, callInstr)
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
		c.code = append(c.code, &InstrAdd{})
	case syntax.Minus:
		c.code = append(c.code, &InstrSub{})
	case syntax.Mul:
		c.code = append(c.code, &InstrMul{})
	case syntax.Div:
		c.code = append(c.code, &InstrDiv{})
	case syntax.GreaterThan:
		c.code = append(c.code, &InstrGreaterThan{})
	case syntax.GreaterThanOrEqual:
		c.code = append(c.code, &InstrGreaterThanOrEqual{})
	case syntax.LessThan:
		c.code = append(c.code, &InstrLessThan{})
	case syntax.LessThanOrEqual:
		c.code = append(c.code, &InstrLessThanOrEqual{})
	case syntax.Equal:
		c.code = append(c.code, &InstrEqual{})
	}
}

func (c *compiler) compileReturnStmt(stmt syntax.Stmt) {
	returnStmt := stmt.(*syntax.ReturnStmt)
	c.compileExpr(returnStmt.Value)
	c.code = append(c.code, &InstrReturn{})
}
