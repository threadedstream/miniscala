package interpreter

import (
	"errors"
	"fmt"
	"github.com/ThreadedStream/miniscala/assert"
	"github.com/ThreadedStream/miniscala/syntax"
	"strconv"
)

var (
	returned bool
)

func checkOpValues(op syntax.Operator, v1, v2 Value, localEnv Environment) {

	if v1.ValueType == Ref {
		v1, _ = lookup(v1.asString(), localEnv, true)
	}

	if v2.ValueType == Ref {
		v2, _ = lookup(v2.asString(), localEnv, true)
	}

	switch op {
	default:
		panic("unknown operation")
	case syntax.Plus:
		if (v1.isString() && v2.isString()) || (v1.isFloat() && v2.isFloat()) {
			return
		}
		panic("v1 and v2 must both be of type string or float")
	case syntax.Minus:
		if v1.isFloat() && v2.isFloat() {
			return
		}
		panic("v1 and v2 must both be of type float")
	case syntax.Mul:
		if v1.isFloat() && v2.isFloat() {
			return
		}
		panic("v1 and v2 must both be of type float")
	case syntax.Div:
		if v1.isFloat() && v2.isFloat() {
			if !v2.isZero() {
				return
			} else {
				panic("division by zero")
			}
		}
	}
}

func checkDefReturnType(v Value, defReturnType, funcName string) {
	if v.ValueType == miniscalaTypeToValueType(defReturnType) {
		return
	}

	panic(fmt.Errorf("(in %s()) expected to return a value of type %s, but it returned %s", funcName, defReturnType, v.valueTypeToStr()))
}

func isReservedFuncCall(funcName string) bool {
	switch funcName {
	default:
		return false
	case "print":
		return true
	}
}

func unwrapValue(value Value) interface{} {
	switch {
	default:
		return nil
	case value.isString():
		return value.asString()
	case value.isFloat():
		return value.asFloat()
	}
}

func callPrintFunc(call *syntax.Call, localEnv Environment) {
	// TODO(threadedstream): allow more arguments to print
	assert.Assert(len(call.ArgList) == 1, "expected a single argument for print got %d", len(call.ArgList))
	value := unwrapValue(visitExpr(call.ArgList[0], localEnv))
	fmt.Printf("%v", value)
}

func dispatchReservedCall(call *syntax.Call, localEnv Environment) {
	switch call.CalleeName.Value {
	case "print":
		callPrintFunc(call, localEnv)
	}
}

func checkAssignmentValidity(name string, localEnv Environment) error {
	// first, check against the presence of value associated with name
	value, _ := lookup(name, localEnv, true)

	// second, check against the possibility to change this value
	if value.Immutable {
		return errors.New("attempt to override val value")
	}

	return nil
}

func resolveRef(v1 Value, localEnv Environment) Value {
	switch {
	default:
		return v1
	case v1.ValueType == Ref:
		v, _ := lookup(v1.asString(), localEnv, true)
		return v
	}
}

func isCondTrue(cond syntax.Operation, localEnv Environment) bool {
	syntax.IsComparisonOp(cond.Op)
	var (
		lhs = visitExpr(cond.Lhs, localEnv)
		rhs = visitExpr(cond.Rhs, localEnv)
	)

	lhs = resolveRef(lhs, localEnv)
	rhs = resolveRef(rhs, localEnv)

	switch cond.Op {
	default:
		return false
	case syntax.GreaterThan:
		if lhs.isString() && rhs.isString() {
			return lhs.asString() > rhs.asString()
		} else if lhs.isFloat() && rhs.isFloat() {
			return lhs.asFloat() > rhs.asFloat()
		} else {
			panic("cast to string and float was unsuccessful")
		}
	case syntax.GreaterThanOrEqual:
		if lhs.isString() && rhs.isString() {
			return lhs.asString() >= rhs.asString()
		} else if lhs.isFloat() && rhs.isFloat() {
			return lhs.asFloat() >= rhs.asFloat()
		} else {
			panic("cast to string and float was unsuccessful")
		}
	case syntax.LessThan:
		if lhs.isString() && rhs.isString() {
			return lhs.asString() < rhs.asString()
		} else if lhs.isFloat() && rhs.isFloat() {
			return lhs.asFloat() < rhs.asFloat()
		} else {
			panic("cast to string and float was unsuccessful")
		}
	case syntax.LessThanOrEqual:
		if lhs.isString() && rhs.isString() {
			return lhs.asString() <= rhs.asString()
		} else if lhs.isFloat() && rhs.isFloat() {
			return lhs.asFloat() <= rhs.asFloat()
		} else {
			panic("cast to string and float was unsuccessful")
		}
	case syntax.Equal:
		return lhs.Value == rhs.Value
	case syntax.NotEqual:
		return lhs.Value != rhs.Value
	}
}

func Execute(program *syntax.Program) {
	for _, stmt := range program.StmtList {
		visitStmt(stmt, nil)
	}
}

func DumpEnvState() {
	state()
}

func visitStmt(stmt syntax.Stmt, localEnv Environment) Value {
	switch stmt.(type) {
	default:
		// we don't allow it in perspective. Right now, we're totally good with that
		return visitExpr(stmt, localEnv)
		//panic(fmt.Errorf("unknown node type %v", reflect.TypeOf(stmt)))
	case *syntax.VarDeclStmt:
		return visitVarDeclStmt(stmt, localEnv)
	case *syntax.ValDeclStmt:
		return visitValDeclStmt(stmt, localEnv)
	case *syntax.BlockStmt:
		return visitBlockStmt(stmt, localEnv)
	case *syntax.IfStmt:
		return visitIfStmt(stmt, localEnv)
	case *syntax.WhileStmt:
		return visitWhileStmt(stmt, localEnv)
	case *syntax.Assignment:
		return visitAssignment(stmt, localEnv)
	case *syntax.DefDeclStmt:
		return visitDefDeclStmt(stmt, localEnv)
	case *syntax.ReturnStmt:
		return visitReturnStmt(stmt, localEnv)
	}
}

func visitExpr(expr syntax.Expr, localEnv Environment) Value {
	switch expr.(type) {
	default:
		return Value{Value: nil}
	case *syntax.Name:
		return visitName(expr)
	case *syntax.BasicLit:
		return visitBasicLit(expr)
	case *syntax.Operation:
		return visitOperation(expr, localEnv)
	case *syntax.Call:
		return visitCall(expr, localEnv)
	}
}

// expressions
func visitName(expr syntax.Expr) Value {
	name := expr.(*syntax.Name)
	v := Value{
		Value:     name.Value,
		ValueType: Ref,
	}
	return v
}

func visitBasicLit(expr syntax.Expr) Value {
	basicLit := expr.(*syntax.BasicLit)
	v := Value{}
	switch basicLit.Kind {
	case syntax.FloatLit:
		v.Value, _ = strconv.ParseFloat(basicLit.Value, 64)
		v.ValueType = Float
	case syntax.StringLit:
		v.Value = basicLit.Value
		v.ValueType = String
	}
	return v
}

func visitOperation(expr syntax.Expr, localEnv Environment) Value {
	operation := expr.(*syntax.Operation)
	switch operation.Op {
	default:
		panic("unknown operation")
	case syntax.Plus:
		lhsValue := visitExpr(operation.Lhs, localEnv)
		rhsValue := visitExpr(operation.Rhs, localEnv)
		checkOpValues(syntax.Plus, lhsValue, rhsValue, localEnv)
		value := add(lhsValue, rhsValue, localEnv)
		return value
	case syntax.Minus:
		lhsValue := visitExpr(operation.Lhs, localEnv)
		rhsValue := visitExpr(operation.Rhs, localEnv)
		checkOpValues(syntax.Minus, lhsValue, rhsValue, localEnv)
		value := sub(lhsValue, rhsValue, localEnv)
		return value
	case syntax.Mul:
		lhsValue := visitExpr(operation.Lhs, localEnv)
		rhsValue := visitExpr(operation.Rhs, localEnv)
		checkOpValues(syntax.Mul, lhsValue, rhsValue, localEnv)
		value := mul(lhsValue, rhsValue, localEnv)
		return value
	case syntax.Div:
		lhsValue := visitExpr(operation.Lhs, localEnv)
		rhsValue := visitExpr(operation.Rhs, localEnv)
		checkOpValues(syntax.Div, lhsValue, rhsValue, localEnv)
		value := div(lhsValue, rhsValue, localEnv)
		return value
	}
}

func visitCall(expr syntax.Expr, localEnv Environment) Value {
	call := expr.(*syntax.Call)

	if isReservedFuncCall(call.CalleeName.Value) {
		// dispatch in case if call to a reserved function has been made
		dispatchReservedCall(call, localEnv)
		return Value{}
	}

	value, _ := lookup(call.CalleeName.Value, nil, true)
	defValue := value.asFunction()

	var funcFrame = make(Environment)
	// TODO(threadedstream): do some checks regarding the number of passed arguments
	for idx, arg := range call.ArgList {
		argValue := visitExpr(arg, localEnv)
		paramName := defValue.DefDeclStmt.ParamList[idx].Name.Value
		store(paramName, argValue, funcFrame, Assign)
	}

	returnValue := visitBlockStmt(defValue.DefDeclStmt.Body, funcFrame)

	checkDefReturnType(returnValue, defValue.DefDeclStmt.ReturnType.(*syntax.Name).Value, defValue.DefDeclStmt.Name.Value)

	// destroy local environment
	funcFrame = nil

	return returnValue
}

func visitReturnStmt(stmt syntax.Stmt, localEnv Environment) Value {
	returnStmt := stmt.(*syntax.ReturnStmt)
	returnValue := visitExpr(returnStmt.Value, localEnv)
	returnValue.Returned = true
	return returnValue
}

// statements
func visitVarDeclStmt(stmt syntax.Stmt, localEnv Environment) Value {
	varDecl := stmt.(*syntax.VarDeclStmt)
	value := visitExpr(varDecl.Rhs, localEnv)
	value.Immutable = false
	store(varDecl.Name.Value, value, localEnv, Declare)
	return Value{}
}

func visitValDeclStmt(stmt syntax.Stmt, localEnv Environment) Value {
	valDecl := stmt.(*syntax.ValDeclStmt)
	value := visitExpr(valDecl.Rhs, localEnv)
	value.Immutable = true
	store(valDecl.Name.Value, value, localEnv, Declare)
	return Value{}
}

func visitBlockStmt(stmt syntax.Stmt, localEnv Environment) Value {
	var (
		block = stmt.(*syntax.BlockStmt)
		value Value
	)

	for _, currStmt := range block.Stmts {
		value = visitStmt(currStmt, localEnv)
		if value.Returned {
			break
		}
	}

	return value
}

func visitIfStmt(stmt syntax.Stmt, localEnv Environment) Value {
	var (
		value  = Value{ValueType: Null}
		ifStmt = stmt.(*syntax.IfStmt)
	)

	if isCondTrue(ifStmt.Cond, localEnv) {
		value = visitStmt(ifStmt.Body, localEnv)
	} else {
		if ifStmt.ElseBody != nil {
			value = visitStmt(ifStmt.ElseBody, localEnv)
		}
	}

	return value
}

func visitWhileStmt(stmt syntax.Stmt, localEnv Environment) Value {
	var whileStmt = stmt.(*syntax.WhileStmt)
	for isCondTrue(whileStmt.Cond, localEnv) {
		visitStmt(whileStmt.Body, localEnv)
	}
	return Value{}
}

func visitAssignment(stmt syntax.Stmt, localEnv Environment) Value {
	assignment := stmt.(*syntax.Assignment)
	lhsValue := visitExpr(assignment.Lhs, localEnv)
	if lhsValue.ValueType != Ref {
		panic("lhs value in assignment should have a value type Ref")
	}
	rhsValue := visitExpr(assignment.Rhs, localEnv)
	if err := checkAssignmentValidity(lhsValue.asString(), localEnv); err != nil {
		panic(err)
	}
	store(lhsValue.asString(), rhsValue, localEnv, Assign)
	return Value{}
}

func visitDefDeclStmt(stmt syntax.Stmt, localEnv Environment) Value {
	var (
		defDeclStmt = stmt.(*syntax.DefDeclStmt)
		returnType  = visitExpr(defDeclStmt.ReturnType, localEnv)
		defValue    = &DefValue{
			DefDeclStmt: defDeclStmt,
			ReturnType:  miniscalaTypeToValueType(returnType.asString()),
		}
		value = Value{
			ValueType: Function,
		}
	)

	value.Value = defValue

	// functions reside in global environment exclusively
	store(defDeclStmt.Name.Value, value, nil, Declare)

	return Value{}
}
