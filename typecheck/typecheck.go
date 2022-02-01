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
	//reservedFunctions = map[string][]backing.ValueType{
	//	"print": {backing.Any},
	//}

	venv = &backing.Venv
	tenv = &backing.Tenv
)

func typecheckError(format string, args ...interface{}) {
	errors = append(errors, typecheckerror{
		fmt:  format,
		args: args,
	})
	hadErrors = true
}

func Typecheck(program *syntax.Program) bool {
	assert.Assert(program != nil, "program is nil!!!")
	*venv = backing.BaseValueEnv()
	*tenv = backing.BaseTypeEnv()
	level := backing.OutermostLevel()
	//typifyReservedFunctions()
	typecheckProgram(program, level)
	if hadErrors {
		for _, err := range errors {
			fmt.Fprintf(os.Stderr, err.fmt, err.args...)
		}
	}
	return hadErrors
}

//func typifyReservedFunctions() {
//	for name, paramTypes := range reservedFunctions {
//		backing.StoreType(name, backing.Unit, true, paramTypes, nil)
//	}
//}

func typecheckUnary(targetType backing.ValueType, op syntax.Operator) (backing.ValueType, bool) {
	switch op {
	default:
		return backing.Undefined, false
	case syntax.Minus:
		switch targetType {
		default:
			return backing.Undefined, false
		case backing.Int:
			return backing.Int, true
		case backing.Float:
			return backing.Float, true
		}
	case syntax.LogicalNot:
		switch targetType {
		default:
			return backing.Undefined, false
		case backing.Bool:
			return backing.Bool, true
		}
	}
}

func typecheckExpr(expr syntax.Expr, level *backing.Level) backing.ValueType {
	switch expr.(type) {
	default:
		errorPos := expr.Pos()
		typecheckError("[%d:%d] undefined type\n", errorPos.Line, errorPos.Column)
		return backing.Undefined
	case *syntax.BasicLit:
		basicLit := expr.(*syntax.BasicLit)
		return backing.LitKindToValueType(basicLit.Kind)
	case *syntax.Name:
		name := expr.(*syntax.Name)
		// name can be a type, as well as a value
		entry := backing.SLook(*venv, backing.SSymbol(name.Value))
		if entry == nil {
			// probably, this is just a type name
			valueType := backing.SLook(*tenv, backing.SSymbol(name.Value))
			if valueType == nil {
				errorPos := name.Pos()
				typecheckError("[%d:%d] name %s is neither a type name nor var, nor val\n", errorPos.Line, errorPos.Column, name.Value)
				return backing.Undefined
			}
			return valueType.(backing.ValueType)
		}
		return entry.(*backing.EnvEntry).ResultType
	case *syntax.Field:
		field := expr.(*syntax.Field)
		return typecheckExpr(field.Type, level)
	case *syntax.Operation:
		var lhsType, rhsType backing.ValueType
		operation := expr.(*syntax.Operation)
		lhsType = typecheckExpr(operation.Lhs, level)
		if operation.Rhs == nil {
			// handling unary
			resultingType, ok := typecheckUnary(lhsType, operation.Op)
			if !ok {
				errorPos := operation.Pos()
				typecheckError("[%d:%d] unary %d didn't expect expression of type %s", errorPos.Line, errorPos.Column,
					operation.Op, backing.ValueTypeToStr(resultingType))
				return backing.Undefined
			}
			return resultingType
		}
		rhsType = typecheckExpr(operation.Rhs, level)
		resultingType, compatible := typesCompatible(lhsType, rhsType, operation.Op)
		if !compatible {
			errorPos := operation.Pos()
			typecheckError("[%d:%d] types %s and %s are not compatible\n",
				errorPos.Line, errorPos.Column,
				backing.ValueTypeToStr(lhsType),
				backing.ValueTypeToStr(rhsType))
		}
		return resultingType

	// although call is the statement, we might implicitly treat it
	// as an expression in that particular case
	case *syntax.Call:
		return typecheckCall(expr, level)
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
		case t1 == backing.Float && t2 == backing.Float:
			return backing.Bool, true
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
	case syntax.Mod:
		switch {
		default:
			return backing.Undefined, false
		case t1 == backing.Int && t2 == backing.Int:
			return backing.Int, true
		}
	}
}

func typecheckProgram(program *syntax.Program, level *backing.Level) {
	for _, stmt := range program.StmtList {
		typecheckStmt(stmt, level)
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
	assignment := stmt.(*syntax.Assignment)
	assigneeName := assignment.Lhs.(*syntax.Name).Value
	// should make assignment's Lhs of type *Name
	lhs := backing.SLook(*venv, backing.SSymbol(assigneeName))
	if lhs == nil {
		errorPos := assignment.Pos()
		typecheckError("[%d:%d] assigning to the undefined variable %s\n",
			errorPos.Line, errorPos.Column,
			assigneeName,
		)
		return
	}

	lhsEntry := lhs.(*backing.EnvEntry)

	rhsType := typecheckExpr(assignment.Rhs, level)

	if lhsEntry.Immutable {
		errorPos := assignment.Pos()
		// reporting the type mismatch issue
		typecheckError("[%d:%d] %s is immutable, thus non-assignable", errorPos.Line, errorPos.Column, assigneeName)
		return
	}

	if lhsEntry.ResultType != rhsType {
		errorPos := assignment.Pos()
		typecheckError("[%d:%d] expected to have rhs type %s, but got %s\n",
			errorPos.Line, errorPos.Column,
			backing.ValueTypeToStr(lhsEntry.ResultType),
			backing.ValueTypeToStr(rhsType))
	}
}

func typecheckCall(stmt syntax.Stmt, level *backing.Level) backing.ValueType {
	callStmt := stmt.(*syntax.Call)
	entry := backing.SLook(*venv, backing.SSymbol(callStmt.CalleeName.Value))
	if entry == nil {
		errorPos := callStmt.Pos()
		typecheckError("[%d:%d] no function with name %s was found\n", errorPos.Line, errorPos.Column, callStmt.CalleeName.Value)
		// bravely return at that point, as it panics if entry is nil
		return backing.Undefined
	}

	calleeEntry := entry.(*backing.EnvEntry)
	if calleeEntry.Kind != backing.EntryFun {
		errorPos := callStmt.Pos()
		typecheckError("[%d:%d] %s is not a function\n", errorPos.Line, errorPos.Column, callStmt.CalleeName.Value)
		return backing.Undefined
	}

	// first, check number of passed parameters
	if len(calleeEntry.ParamTypes) != len(callStmt.ArgList) {
		errorPos := callStmt.Pos()
		typecheckError("[%d:%d] function %s expects %d parameters, but %d were provided\n", errorPos.Line, errorPos.Column,
			callStmt.CalleeName.Value, len(calleeEntry.ParamTypes), len(callStmt.ArgList))
		return backing.Undefined
	}

	var valueTypes []backing.ValueType
	for _, arg := range callStmt.ArgList {
		argType := typecheckExpr(arg, level)
		if argType == backing.Undefined {
			return backing.Undefined
		}
		valueTypes = append(valueTypes, argType)
	}

	for idx, paramType := range calleeEntry.ParamTypes {
		if !backing.TypesEqual(paramType, valueTypes[idx]) {
			errorPos := callStmt.Pos()
			typecheckError("[%d:%d] parameter %d expected type %s, but %s was provided\n", errorPos.Line, errorPos.Column,
				idx+1, backing.ValueTypeToStr(paramType), backing.ValueTypeToStr(valueTypes[idx]))
			return backing.Undefined
		}
	}

	return calleeEntry.ResultType
}

func typecheckBlockStmt(stmt syntax.Stmt, level *backing.Level) (backing.ValueType, scanner.Position) {
	blockStmt := stmt.(*syntax.BlockStmt)
	for _, decStmt := range blockStmt.Stmts {
		switch decStmt.(type) {
		default:
			typecheckStmt(decStmt, level)
		case *syntax.ReturnStmt:
			// handling the special case
			return typecheckReturnStmt(decStmt, level)
		}
	}
	return backing.Unit, scanner.Position{}
}

func typecheckVarDeclStmt(stmt syntax.Stmt, level *backing.Level) {
	varDeclStmt := stmt.(*syntax.VarDeclStmt)
	if syntax.IsKeyword(varDeclStmt.Name.Value) {
		errorPos := varDeclStmt.Pos()
		typecheckError("[%d:%d] name %s is reserved\n", errorPos.Line, errorPos.Column, varDeclStmt.Name.Value)
		return
	}
	inferredType := typecheckExpr(varDeclStmt.Rhs, level)
	backing.SEnter(
		*venv, backing.SSymbol(varDeclStmt.Name.Value), backing.MakeVarEntry(
			varDeclStmt.Name.Value,
			level,
			inferredType,
			false,
		),
	)
}

func typecheckValDeclStmt(stmt syntax.Stmt, level *backing.Level) {
	valDeclStmt := stmt.(*syntax.ValDeclStmt)
	// first, check if declared name is a keyword or not
	if syntax.IsKeyword(valDeclStmt.Name.Value) {
		errorPos := valDeclStmt.Pos()
		typecheckError("[%d:%d] name %s is reserved\n", errorPos.Line, errorPos.Column, valDeclStmt.Name.Value)
		return
	}
	valueType := typecheckExpr(valDeclStmt.Rhs, level)
	backing.SEnter(
		*venv, backing.SSymbol(valDeclStmt.Name.Value), backing.MakeVarEntry(
			valDeclStmt.Name.Value,
			level,
			valueType,
			true,
		),
	)
}

func typecheckIfStmt(stmt syntax.Stmt, level *backing.Level) {
	ifStmt := stmt.(*syntax.IfStmt)
	condValueType := typecheckExpr(ifStmt.Cond, level)
	if condValueType != backing.Bool {
		errorPos := ifStmt.Pos()
		typecheckError("[%d:%d] condition is not of bool type\n", errorPos.Line, errorPos.Column)
		return
	}
	typecheckBlockStmt(ifStmt.Body, level)
	if ifStmt.ElseBody != nil {
		typecheckStmt(ifStmt.ElseBody, level)
	}
}

func typecheckWhileStmt(stmt syntax.Stmt, level *backing.Level) {
	whileStmt := stmt.(*syntax.WhileStmt)
	condValueType := typecheckExpr(whileStmt.Cond, level)
	if condValueType != backing.Bool {
		errorPos := whileStmt.Pos()
		typecheckError("[%d:%d] condition is not of bool type\n", errorPos.Line, errorPos.Column)
		return
	}
	backing.SBeginScope(*venv)
	typecheckBlockStmt(whileStmt.Body, level)
	backing.SEndScope(*venv)
}

func typecheckDefDeclStmt(stmt syntax.Stmt, level *backing.Level) {
	defDeclStmt := stmt.(*syntax.DefDeclStmt)
	var paramTypes []backing.ValueType

	expectedReturnType := typecheckExpr(defDeclStmt.ReturnType, level)
	funLevel := backing.NewLevel(defDeclStmt.Name.Value, level)
	for _, param := range defDeclStmt.ParamList {
		paramTypes = append(paramTypes, typecheckField(param, funLevel))
	}

	backing.SEnter(
		*venv, backing.SSymbol(defDeclStmt.Name.Value), backing.MakeFunEntry(
			defDeclStmt.Name.Value,
			paramTypes,
			funLevel,
			expectedReturnType),
	)

	backing.SBeginScope(*venv)
	for idx, param := range defDeclStmt.ParamList {
		backing.SEnter(
			*venv, backing.SSymbol(param.Name.Value), backing.MakeVarEntry(
				param.Name.Value,
				funLevel,
				paramTypes[idx],
				false,
			),
		)
	}

	returnType, pos := typecheckBlockStmt(defDeclStmt.Body, level)
	if returnType != expectedReturnType {
		errorPos := pos
		typecheckError("[%d:%d] expected return type %s but got %s\n",
			errorPos.Line,
			errorPos.Column,
			backing.ValueTypeToStr(expectedReturnType),
			backing.ValueTypeToStr(returnType))
	}

	backing.SEndScope(*venv)
}

func typecheckReturnStmt(stmt syntax.Stmt, level *backing.Level) (backing.ValueType, scanner.Position) {
	returnStmt := stmt.(*syntax.ReturnStmt)
	returnType := typecheckExpr(returnStmt.Value, level)
	pos := returnStmt.Pos()
	return returnType, pos
}

func typecheckField(field *syntax.Field, level *backing.Level) backing.ValueType {
	valueType := typecheckExpr(field, level)
	//backing.StoreType(field.Name.Value, valueType, false, nil, )
	return valueType
}
