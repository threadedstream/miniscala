package typecheck

import (
	"fmt"
	"github.com/ThreadedStream/miniscala/assert"
	"github.com/ThreadedStream/miniscala/backing"
	"github.com/ThreadedStream/miniscala/syntax"
	"os"
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
	typifyReservedFunctions()
	typecheckProgram(program)
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

func inferValueType(expr syntax.Expr, typeEnv backing.TypeEnv) backing.ValueType {
	switch expr.(type) {
	default:
		typecheckError("unknown node in inferValueType()\n")
		return backing.Undefined
	case *syntax.BasicLit:
		basicLit := expr.(*syntax.BasicLit)
		return backing.LitKindToValueType(basicLit.Kind)
	case *syntax.Name:
		name := expr.(*syntax.Name)
		typeInfo, ok := backing.LookupType(name.Value, typeEnv, false)
		if !ok {
			// probably, this is just a type name
			valueType := backing.MiniscalaTypeToValueType(name.Value)
			if valueType == backing.Undefined {
				typecheckError("name %s is neither a type name nor var, nor val\n", name.Value)
			}
			return valueType
		}
		return typeInfo.ValueType
	case *syntax.Field:
		field := expr.(*syntax.Field)
		return inferValueType(field.Type, typeEnv)
	case *syntax.Operation:
		operation := expr.(*syntax.Operation)
		lhsType := inferValueType(operation.Lhs, typeEnv)
		rhsType := inferValueType(operation.Rhs, typeEnv)
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

func typecheckProgram(program *syntax.Program) {
	for _, stmt := range program.StmtList {
		typecheckStmt(stmt, nil)
	}
}

func typecheckStmt(stmt syntax.Stmt, typeEnv backing.TypeEnv) {
	switch stmt.(type) {
	// add a block statement
	case *syntax.VarDeclStmt:
		typecheckVarDeclStmt(stmt, typeEnv)
	case *syntax.ValDeclStmt:
		typecheckValDeclStmt(stmt, typeEnv)
	case *syntax.DefDeclStmt:
		typecheckDefDeclStmt(stmt)
	case *syntax.IfStmt:
		typecheckIfStmt(stmt, typeEnv)
	case *syntax.WhileStmt:
		typecheckWhileStmt(stmt, typeEnv)
	case *syntax.Call:
		typecheckCall(stmt, typeEnv)
	case *syntax.BlockStmt:
		typecheckBlockStmt(stmt, backing.Undefined, typeEnv)
	case *syntax.Assignment:
		typecheckAssignment(stmt, typeEnv)
	}
}

func typecheckAssignment(stmt syntax.Stmt, typeEnv backing.TypeEnv) {
	// TODO(threadedstream):
}

func typecheckCall(stmt syntax.Stmt, typeEnv backing.TypeEnv) {
	callStmt := stmt.(*syntax.Call)
	calleeInfo, ok := backing.LookupType(callStmt.CalleeName.Value, nil, false)
	if !ok {
		errorPos := callStmt.Pos()
		typecheckError("[%d:%d] no function with name %s was found\n", errorPos.Line, errorPos.Column, callStmt.CalleeName.Value)
	}
	// first, check number of passed parameters
	if len(calleeInfo.ParamTypes) != len(callStmt.ArgList) {
		errorPos := callStmt.Pos()
		typecheckError("[%d:%d] function %s expects %d parameters, but %d were provided\n", errorPos.Line, errorPos.Column,
			callStmt.CalleeName.Value, len(calleeInfo.ParamTypes), len(callStmt.ArgList))
	}

	var valueTypes []backing.ValueType
	for _, arg := range callStmt.ArgList {
		valueTypes = append(valueTypes, inferValueType(arg, typeEnv))
	}

	for idx, paramType := range calleeInfo.ParamTypes {
		if !backing.TypesEqual(paramType, valueTypes[idx]) {
			errorPos := callStmt.Pos()
			typecheckError("[%d:%d] arg %d expected type %s, but %s was provided", errorPos.Line, errorPos.Column,
				idx, backing.ValueTypeToStr(paramType), backing.ValueTypeToStr(valueTypes[idx]))
		}
	}
}

func typecheckBlockStmt(stmt syntax.Stmt, expectedReturnType backing.ValueType, typeEnv backing.TypeEnv) {
	blockStmt := stmt.(*syntax.BlockStmt)
	for _, decStmt := range blockStmt.Stmts {
		switch decStmt.(type) {
		default:
			typecheckStmt(decStmt, typeEnv)
		case *syntax.ReturnStmt:
			// handling the special case
			typecheckReturnStmt(stmt, expectedReturnType, typeEnv)
		}
	}
}

func typecheckVarDeclStmt(stmt syntax.Stmt, typeEnv backing.TypeEnv) {
	varDeclStmt := stmt.(*syntax.VarDeclStmt)
	if syntax.IsKeyword(varDeclStmt.Name.Value) {
		errorPos := varDeclStmt.Pos()
		typecheckError("[%d:%d] name %s is reserved\n", errorPos.Line, errorPos.Column, varDeclStmt.Name.Value)
	}
	inferredType := inferValueType(varDeclStmt.Rhs, typeEnv)
	backing.StoreType(varDeclStmt.Name.Value, inferredType, false, nil, typeEnv)
}

func typecheckValDeclStmt(stmt syntax.Stmt, typeEnv backing.TypeEnv) {
	valDeclStmt := stmt.(*syntax.ValDeclStmt)
	// first, check if declared name is a keyword or not
	if syntax.IsKeyword(valDeclStmt.Name.Value) {
		errorPos := valDeclStmt.Pos()
		typecheckError("[%d:%d] name %s is reserved\n", errorPos.Line, errorPos.Column, valDeclStmt.Name.Value)
	}
	// TODO(threadedstream): there's a room for a constant folding optimization
	valueType := inferValueType(valDeclStmt.Rhs, typeEnv)
	backing.StoreType(valDeclStmt.Name.Value, valueType, true, nil, typeEnv)
}

func typecheckIfStmt(stmt syntax.Stmt, typeEnv backing.TypeEnv) {
	ifStmt := stmt.(*syntax.IfStmt)
	condValueType := inferValueType(&ifStmt.Cond, typeEnv)
	if condValueType != backing.Bool {
		errorPos := ifStmt.Pos()
		typecheckError("[%d:%d] condition is not of bool type", errorPos.Line, errorPos.Column)
	}
	typecheckBlockStmt(ifStmt.Body, backing.Undefined, typeEnv)
	typecheckStmt(ifStmt.ElseBody, typeEnv)
}

func typecheckWhileStmt(stmt syntax.Stmt, typeEnv backing.TypeEnv) {
	whileStmt := stmt.(*syntax.WhileStmt)
	condValueType := inferValueType(&whileStmt.Cond, typeEnv)
	if condValueType != backing.Bool {
		errorPos := whileStmt.Pos()
		typecheckError("[%d:%d] condition is not of bool type", errorPos.Line, errorPos.Column)
	}
	typecheckBlockStmt(whileStmt.Body, backing.Undefined, typeEnv)
}

func typecheckDefDeclStmt(stmt syntax.Stmt) {
	defDeclStmt := stmt.(*syntax.DefDeclStmt)
	typeEnv := make(backing.TypeEnv)
	var paramTypes []backing.ValueType
	for _, param := range defDeclStmt.ParamList {
		paramTypes = append(paramTypes, typecheckField(param, typeEnv))
	}
	backing.StoreType(defDeclStmt.Name.Value, backing.Function, true, paramTypes, nil)
	expectedReturnType := inferValueType(defDeclStmt.ReturnType, typeEnv)
	typecheckBlockStmt(defDeclStmt.Body, expectedReturnType, typeEnv)
}

func typecheckReturnStmt(stmt syntax.Stmt, expectedReturnType backing.ValueType, typeEnv backing.TypeEnv) {
	returnStmt := stmt.(*syntax.ReturnStmt)
	actualReturnType := inferValueType(returnStmt.Value, typeEnv)
	if actualReturnType != expectedReturnType {
		errorPos := returnStmt.Pos()
		typecheckError("[%d:%d] expected to return %s, but got %s\n", errorPos.Line, errorPos.Column,
			backing.ValueTypeToStr(expectedReturnType),
			backing.ValueTypeToStr(actualReturnType))
	}
}

func typecheckField(field *syntax.Field, typeEnv backing.TypeEnv) backing.ValueType {
	valueType := inferValueType(field, typeEnv)
	backing.StoreType(field.Name.Value, valueType, false, nil, typeEnv)
	return valueType
}
