package typecheck

import (
	"fmt"
	"github.com/ThreadedStream/miniscala/assert"
	"github.com/ThreadedStream/miniscala/backing"
	"github.com/ThreadedStream/miniscala/syntax"
	"os"
	"text/scanner"
)

type typecheckerror struct {
	fmt  string
	args []interface{}
}

var (
	errors    []typecheckerror
	hadErrors = false
	// map from reserved functions' names to the type of their parameters
	reservedFunctions = map[string][]backing.ValueType{
		"print": {backing.Any},
	}

	venv backing.SymbolTable
	tenv backing.SymbolTable
)

func typecheckError(format string, args ...interface{}) {
	errors = append(errors, typecheckerror{
		fmt:  format,
		args: args,
	})
	hadErrors = true
}

func Typecheck(program *syntax.Program) {
	assert.Assert(program != nil, "program is nil!!!")
	venv = backing.BaseValueEnv()
	tenv = backing.BaseTypeEnv()
	level := backing.OutermostLevel()
	//typifyReservedFunctions()
	typecheckProgram(program, level)
	if hadErrors {
		for _, err := range errors {
			fmt.Fprintf(os.Stderr, err.fmt, err.args...)
		}
		os.Exit(1)
	}
}

func typifyReservedFunctions() {
	for name, paramTypes := range reservedFunctions {
		backing.StoreType(name, backing.Unit, true, paramTypes, nil)
	}
}

func typecheckExpr(expr syntax.Expr, level *backing.Level) backing.ValueType {
	switch expr.(type) {
	default:
		typecheckError("unknown node in typecheckExpr()\n")
		return backing.Undefined
	case *syntax.BasicLit:
		basicLit := expr.(*syntax.BasicLit)
		return backing.LitKindToValueType(basicLit.Kind)
	case *syntax.Name:
		name := expr.(*syntax.Name)
		// name can be a type, as well as a value
		entry := backing.SLook(venv, backing.SSymbol(name.Value)).(*backing.EnvEntry)
		if entry != nil {
			// probably, this is just a type name
			valueType := backing.SLook(tenv, backing.SSymbol(name.Value))
			if valueType == nil {
				typecheckError("name %s is neither a type name nor var, nor val\n", name.Value)
			}
			return valueType.(backing.ValueType)
		}
		return entry.ResultType
	case *syntax.Field:
		field := expr.(*syntax.Field)
		return typecheckExpr(field.Type, level)
	case *syntax.Operation:
		operation := expr.(*syntax.Operation)
		lhsType := typecheckExpr(operation.Lhs, level)
		rhsType := typecheckExpr(operation.Rhs, level)
		resultingType, compatible := typesCompatible(lhsType, rhsType, operation.Op)
		if !compatible {
			errorPos := operation.Pos()
			typecheckError("[%d:%d] types %s and %s are not compatible\n",
				errorPos.Line, errorPos.Column,
				backing.ValueTypeToStr(lhsType),
				backing.ValueTypeToStr(rhsType))
		}
		return resultingType
	}
}

// typesCompatible checks against compatibility of passed types and returns
// a resulting type upon success
func typesCompatible(t1, t2 backing.ValueType, op syntax.Operator) (backing.ValueType, bool) {
	switch op {
	default:
		return backing.Undefined, false
	case syntax.Plus:
		switch {
		default:
			return backing.Undefined, false
		case t1 == backing.Float && t2 == backing.Int:
			return backing.Float, true
		case t1 == backing.Int && t2 == backing.Float:
			return backing.Float, true
		case t1 == backing.Float && t2 == backing.Float:
			return backing.Float, true
		case t1 == backing.Int && t2 == backing.Int:
			return backing.Int, true
		case t1 == backing.String && t2 == backing.String:
			return backing.String, true
		}
	case syntax.Minus, syntax.Mul, syntax.Div:
		switch {
		default:
			return backing.Undefined, false
		case t1 == backing.Float && t2 == backing.Int:
			return backing.Float, true
		case t1 == backing.Int && t2 == backing.Float:
			return backing.Float, true
		case t1 == backing.Float && t2 == backing.Float:
			return backing.Float, true
		case t1 == backing.Int && t2 == backing.Int:
			return backing.Int, true
		}
	case syntax.GreaterThan, syntax.GreaterThanOrEqual,
		syntax.LessThan, syntax.LessThanOrEqual,
		syntax.Equal, syntax.NotEqual:

		switch {
		default:
			return backing.Undefined, false
		case t1 == backing.Float && t2 == backing.Int:
			return backing.Bool, true
		case t1 == backing.Int && t2 == backing.Float:
			return backing.Bool, true
		case t1 == backing.Int && t2 == backing.Int:
			return backing.Bool, true
		case t1 == backing.String && t2 == backing.String:
			return backing.Bool, true
		case t1 == backing.Bool && t2 == backing.Bool:
			return backing.Bool, true
		}
	}
}

func typecheckProgram(program *syntax.Program, level *backing.Level) {
	for _, stmt := range program.StmtList {
		typecheckStmt(stmt, nil)
	}
}

func typecheckStmt(stmt syntax.Stmt, level *backing.Level) {
	switch stmt.(type) {
	// add a block statement
	case *syntax.VarDeclStmt:
		typecheckVarDeclStmt(stmt, level)
	case *syntax.ValDeclStmt:
		typecheckValDeclStmt(stmt, level)
	case *syntax.DefDeclStmt:
		typecheckDefDeclStmt(stmt, level)
	case *syntax.IfStmt:
		typecheckIfStmt(stmt, level)
	case *syntax.WhileStmt:
		typecheckWhileStmt(stmt, level)
	case *syntax.Call:
		typecheckCall(stmt, level)
	case *syntax.BlockStmt:
		typecheckBlockStmt(stmt, level)
	case *syntax.Assignment:
		typecheckAssignment(stmt, level)
	}
}

func typecheckAssignment(stmt syntax.Stmt, level *backing.Level) {
	// TODO(threadedstream):
}

func typecheckCall(stmt syntax.Stmt, level *backing.Level) {
	callStmt := stmt.(*syntax.Call)
	calleeEntry := backing.SLook(venv, backing.SSymbol(callStmt.CalleeName.Value)).(*backing.EnvEntry)
	if calleeEntry == nil {
		errorPos := callStmt.Pos()
		typecheckError("[%d:%d] no function with name %s was found\n", errorPos.Line, errorPos.Column, callStmt.CalleeName.Value)
	}
	if calleeEntry.Kind != backing.EntryFun {
		errorPos := callStmt.Pos()
		typecheckError("[%d:%d] %s is not a function\n", errorPos.Line, errorPos.Column, callStmt.CalleeName.Value)
	}

	// first, check number of passed parameters
	if len(calleeEntry.ParamTypes) != len(callStmt.ArgList) {
		errorPos := callStmt.Pos()
		typecheckError("[%d:%d] function %s expects %d parameters, but %d were provided\n", errorPos.Line, errorPos.Column,
			callStmt.CalleeName.Value, len(calleeEntry.ParamTypes), len(callStmt.ArgList))
	}

	var valueTypes []backing.ValueType
	for _, arg := range callStmt.ArgList {
		valueTypes = append(valueTypes, typecheckExpr(arg, level))
	}

	for idx, paramType := range calleeEntry.ParamTypes {
		if !backing.TypesEqual(paramType, valueTypes[idx]) {
			errorPos := callStmt.Pos()
			typecheckError("[%d:%d] arg %d expected type %s, but %s was provided", errorPos.Line, errorPos.Column,
				idx, backing.ValueTypeToStr(paramType), backing.ValueTypeToStr(valueTypes[idx]))
		}
	}
}

func typecheckBlockStmt(stmt syntax.Stmt, level *backing.Level) []struct {
	exprType backing.ValueType
	pos      scanner.Position
} {
	blockStmt := stmt.(*syntax.BlockStmt)
	returnValueTypes := make([]struct {
		exprType backing.ValueType
		pos      scanner.Position
	}, 0)
	for _, decStmt := range blockStmt.Stmts {
		switch decStmt.(type) {
		default:
			typecheckStmt(decStmt, level)
		case *syntax.ReturnStmt:
			// handling the special case
			returnValueTypes = append(returnValueTypes, typecheckReturnStmt(stmt, level))
		}
	}
	return returnValueTypes
}

func typecheckVarDeclStmt(stmt syntax.Stmt, level *backing.Level) {
	varDeclStmt := stmt.(*syntax.VarDeclStmt)
	if syntax.IsKeyword(varDeclStmt.Name.Value) {
		errorPos := varDeclStmt.Pos()
		typecheckError("[%d:%d] name %s is reserved\n", errorPos.Line, errorPos.Column, varDeclStmt.Name.Value)
	}
	inferredType := typecheckExpr(varDeclStmt.Rhs, level)
	backing.SEnter(venv, backing.SSymbol(varDeclStmt.Name.Value), backing.MakeVarEntry(
		varDeclStmt.Name.Value,
		level,
		inferredType),
	)
	//backing.StoreType(varDeclStmt.Name.Value, inferredType, false, nil, level)
}

func typecheckValDeclStmt(stmt syntax.Stmt, level *backing.Level) {
	valDeclStmt := stmt.(*syntax.ValDeclStmt)
	// first, check if declared name is a keyword or not
	if syntax.IsKeyword(valDeclStmt.Name.Value) {
		errorPos := valDeclStmt.Pos()
		typecheckError("[%d:%d] name %s is reserved\n", errorPos.Line, errorPos.Column, valDeclStmt.Name.Value)
	}
	// TODO(threadedstream): there's a room for a constant folding optimization
	valueType := typecheckExpr(valDeclStmt.Rhs, level)
	backing.SEnter(
		venv, backing.SSymbol(valDeclStmt.Name.Value), backing.MakeVarEntry(
			valDeclStmt.Name.Value,
			level,
			valueType,
		),
	)
	//backing.StoreType(valDeclStmt.Name.Value, valueType, true, nil, level)
}

func typecheckIfStmt(stmt syntax.Stmt, level *backing.Level) {
	ifStmt := stmt.(*syntax.IfStmt)
	condValueType := typecheckExpr(&ifStmt.Cond, level)
	if condValueType != backing.Bool {
		errorPos := ifStmt.Pos()
		typecheckError("[%d:%d] condition is not of bool type", errorPos.Line, errorPos.Column)
	}
	typecheckBlockStmt(ifStmt.Body, level)
	typecheckStmt(ifStmt.ElseBody, level)
}

func typecheckWhileStmt(stmt syntax.Stmt, level *backing.Level) {
	whileStmt := stmt.(*syntax.WhileStmt)
	condValueType := typecheckExpr(&whileStmt.Cond, level)
	if condValueType != backing.Bool {
		errorPos := whileStmt.Pos()
		typecheckError("[%d:%d] condition is not of bool type", errorPos.Line, errorPos.Column)
	}
	typecheckBlockStmt(whileStmt.Body, level)
}

func typecheckDefDeclStmt(stmt syntax.Stmt, level *backing.Level) {
	defDeclStmt := stmt.(*syntax.DefDeclStmt)
	var paramTypes []backing.ValueType

	expectedReturnType := typecheckExpr(defDeclStmt.ReturnType, level)
	funLevel := backing.NewLevel(defDeclStmt.Name.Value, level)
	backing.SEnter(
		venv, backing.SSymbol(defDeclStmt.Name.Value), backing.MakeFunEntry(
			defDeclStmt.Name.Value,
			paramTypes,
			funLevel,
			expectedReturnType),
	)

	backing.SBeginScope(venv)
	for _, param := range defDeclStmt.ParamList {
		backing.SEnter(
			venv, backing.SSymbol(param.Name.Value), backing.MakeVarEntry(
				param.Name.Value,
				level,
				typecheckField(param, funLevel),
			),
		)
	}

	returnTypes := typecheckBlockStmt(defDeclStmt.Body, level)
	for _, returnType := range returnTypes {
		if returnType.exprType != expectedReturnType {
			errorPos := returnType.pos
			typecheckError("[%d:%d] expected return type %s but got %s",
				errorPos.Line,
				errorPos.Column,
				backing.ValueTypeToStr(expectedReturnType),
				backing.ValueTypeToStr(returnType.exprType))
		}
	}

	backing.SEndScope(venv)
}

func typecheckReturnStmt(stmt syntax.Stmt, level *backing.Level) struct {
	exprType backing.ValueType
	pos      scanner.Position
} {
	returnStmt := stmt.(*syntax.ReturnStmt)
	returnType := typecheckExpr(returnStmt.Value, level)
	pos := returnStmt.Pos()
	return struct {
		exprType backing.ValueType
		pos      scanner.Position
	}{
		exprType: returnType,
		pos:      pos,
	}
}

func typecheckField(field *syntax.Field, level *backing.Level) backing.ValueType {
	valueType := typecheckExpr(field, level)
	//backing.StoreType(field.Name.Value, valueType, false, nil, )
	return valueType
}
