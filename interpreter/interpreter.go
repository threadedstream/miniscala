package interpreter

import (
	"errors"
	"fmt"
	"github.com/ThreadedStream/miniscala/assert"
	"github.com/ThreadedStream/miniscala/syntax"
	"github.com/ThreadedStream/miniscala/vm"
	"strconv"
)

var (
	returned bool
)

func checkOpValues(op syntax.Operator, v1, v2 vm.Value, localEnv Environment) {

	if v1.ValueType == vm.Ref {
		v1, _ = lookup(v1.AsString(), localEnv, true)
	}

	if v2.ValueType == vm.Ref {
		v2, _ = lookup(v2.AsString(), localEnv, true)
	}

	switch op {
	default:
		panic("unknown operation")
	case syntax.Plus:
		if (v1.IsString() && v2.IsString()) || (v1.IsFloat() && v2.IsFloat()) {
			return
		}
		panic("v1 and v2 must both be of type string or float")
	case syntax.Minus:
		if v1.IsFloat() && v2.IsFloat() {
			return
		}
		panic("v1 and v2 must both be of type float")
	case syntax.Mul:
		if v1.IsFloat() && v2.IsFloat() {
			return
		}
		panic("v1 and v2 must both be of type float")
	case syntax.Div:
		if v1.IsFloat() && v2.IsFloat() {
			if !v2.IsZero() {
				return
			} else {
				panic("division by zero")
			}
		}
	}
}

func checkDefReturnType(v vm.Value, defReturnType, funcName string) {
	if v.ValueType == vm.MiniscalaTypeToValueType(defReturnType) {
		return
	}

	panic(fmt.Errorf("(in %s()) expected to return a value of type %s, but it returned %s", funcName, defReturnType, v.ValueTypeToStr()))
}

func isReservedFuncCall(funcName string) bool {
	switch funcName {
	default:
		return false
	case "print":
		return true
	}
}

func unwrapValue(value vm.Value) interface{} {
	switch {
	default:
		return nil
	case value.IsString():
		return value.AsString()
	case value.IsFloat():
		return value.AsFloat()
	}
}

func add(v1, v2 vm.Value, localEnv Environment) vm.Value {
	if v1.ValueType == vm.Ref {
		v1, _ = lookup(v1.AsString(), localEnv, true)
	}
	if v2.ValueType == vm.Ref {
		v2, _ = lookup(v2.AsString(), localEnv, true)
	}

	switch v1.Value.(type) {
	default:
		return vm.Value{
			Value:     nil,
			ValueType: vm.Undefined,
		}
	case float64:
		return vm.Value{
			Value:     v1.AsFloat() + v2.AsFloat(),
			ValueType: vm.Float,
		}
	case string:
		return vm.Value{
			Value:     v1.AsString() + v2.AsString(),
			ValueType: vm.String,
		}
	}
}

func sub(v1, v2 vm.Value, localEnv Environment) vm.Value {
	if v1.ValueType == vm.Ref {
		v1, _ = lookup(v1.AsString(), localEnv, true)
	}
	if v2.ValueType == vm.Ref {
		v2, _ = lookup(v2.AsString(), localEnv, true)
	}

	switch v1.Value.(type) {
	case float64:
		return vm.Value{
			Value:     v1.AsFloat() - v2.AsFloat(),
			ValueType: vm.Float,
		}
	default:
		return vm.Value{
			Value:     nil,
			ValueType: vm.Undefined,
		}
	}
}

func mul(v1, v2 vm.Value, localEnv Environment) vm.Value {
	if v1.ValueType == vm.Ref {
		v1, _ = lookup(v1.AsString(), localEnv, true)
	}
	if v2.ValueType == vm.Ref {
		v2, _ = lookup(v2.AsString(), localEnv, true)
	}

	switch v1.Value.(type) {
	default:
		return vm.Value{
			Value:     nil,
			ValueType: vm.Undefined,
		}
	case float64:
		return vm.Value{
			Value:     v1.AsFloat() * v2.AsFloat(),
			ValueType: vm.Float,
		}
	}
}

func div(v1, v2 vm.Value, localEnv Environment) vm.Value {
	if v1.ValueType == vm.Ref {
		v1, _ = lookup(v1.AsString(), localEnv, true)
	}

	if v2.ValueType == vm.Ref {
		v2, _ = lookup(v2.AsString(), localEnv, true)
	}

	switch v1.Value.(type) {
	default:
		return vm.Value{
			Value:     nil,
			ValueType: vm.Undefined,
		}
	case float64:
		return vm.Value{
			Value:     v1.AsFloat() / v2.AsFloat(),
			ValueType: vm.Float,
		}
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

func resolveRef(v1 vm.Value, localEnv Environment) vm.Value {
	switch {
	default:
		return v1
	case v1.ValueType == vm.Ref:
		v, _ := lookup(v1.AsString(), localEnv, true)
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
		if lhs.IsString() && rhs.IsString() {
			return lhs.AsString() > rhs.AsString()
		} else if lhs.IsFloat() && rhs.IsFloat() {
			return lhs.AsFloat() > rhs.AsFloat()
		} else {
			panic("cast to string and float was unsuccessful")
		}
	case syntax.GreaterThanOrEqual:
		if lhs.IsString() && rhs.IsString() {
			return lhs.AsString() >= rhs.AsString()
		} else if lhs.IsFloat() && rhs.IsFloat() {
			return lhs.AsFloat() >= rhs.AsFloat()
		} else {
			panic("cast to string and float was unsuccessful")
		}
	case syntax.LessThan:
		if lhs.IsString() && rhs.IsString() {
			return lhs.AsString() < rhs.AsString()
		} else if lhs.IsFloat() && rhs.IsFloat() {
			return lhs.AsFloat() < rhs.AsFloat()
		} else {
			panic("cast to string and float was unsuccessful")
		}
	case syntax.LessThanOrEqual:
		if lhs.IsString() && rhs.IsString() {
			return lhs.AsString() <= rhs.AsString()
		} else if lhs.IsFloat() && rhs.IsFloat() {
			return lhs.AsFloat() <= rhs.AsFloat()
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

func visitStmt(stmt syntax.Stmt, localEnv Environment) vm.Value {
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

func visitExpr(expr syntax.Expr, localEnv Environment) vm.Value {
	switch expr.(type) {
	default:
		return vm.Value{Value: nil}
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
func visitName(expr syntax.Expr) vm.Value {
	name := expr.(*syntax.Name)
	v := vm.Value{
		Value:     name.Value,
		ValueType: vm.Ref,
	}
	return v
}

func visitBasicLit(expr syntax.Expr) vm.Value {
	basicLit := expr.(*syntax.BasicLit)
	v := vm.Value{}
	switch basicLit.Kind {
	case syntax.FloatLit:
		v.Value, _ = strconv.ParseFloat(basicLit.Value, 64)
		v.ValueType = vm.Float
	case syntax.StringLit:
		v.Value = basicLit.Value
		v.ValueType = vm.String
	}
	return v
}

func visitOperation(expr syntax.Expr, localEnv Environment) vm.Value {
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

func visitCall(expr syntax.Expr, localEnv Environment) vm.Value {
	call := expr.(*syntax.Call)

	if isReservedFuncCall(call.CalleeName.Value) {
		// dispatch in case if call to a reserved function has been made
		dispatchReservedCall(call, localEnv)
		return vm.Value{}
	}

	value, _ := lookup(call.CalleeName.Value, nil, true)
	defValue := value.AsFunction()

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

func visitReturnStmt(stmt syntax.Stmt, localEnv Environment) vm.Value {
	returnStmt := stmt.(*syntax.ReturnStmt)
	returnValue := visitExpr(returnStmt.Value, localEnv)
	returnValue.Returned = true
	return returnValue
}

// statements
func visitVarDeclStmt(stmt syntax.Stmt, localEnv Environment) vm.Value {
	varDecl := stmt.(*syntax.VarDeclStmt)
	value := visitExpr(varDecl.Rhs, localEnv)
	value.Immutable = false
	store(varDecl.Name.Value, value, localEnv, Declare)
	return vm.Value{}
}

func visitValDeclStmt(stmt syntax.Stmt, localEnv Environment) vm.Value {
	valDecl := stmt.(*syntax.ValDeclStmt)
	value := visitExpr(valDecl.Rhs, localEnv)
	value.Immutable = true
	store(valDecl.Name.Value, value, localEnv, Declare)
	return vm.Value{}
}

func visitBlockStmt(stmt syntax.Stmt, localEnv Environment) vm.Value {
	var (
		block = stmt.(*syntax.BlockStmt)
		value vm.Value
	)

	for _, currStmt := range block.Stmts {
		value = visitStmt(currStmt, localEnv)
		if value.Returned {
			break
		}
	}

	return value
}

func visitIfStmt(stmt syntax.Stmt, localEnv Environment) vm.Value {
	var (
		value  = vm.Value{ValueType: vm.Null}
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

func visitWhileStmt(stmt syntax.Stmt, localEnv Environment) vm.Value {
	var whileStmt = stmt.(*syntax.WhileStmt)
	for isCondTrue(whileStmt.Cond, localEnv) {
		visitStmt(whileStmt.Body, localEnv)
	}
	return vm.Value{}
}

func visitAssignment(stmt syntax.Stmt, localEnv Environment) vm.Value {
	assignment := stmt.(*syntax.Assignment)
	lhsValue := visitExpr(assignment.Lhs, localEnv)
	if lhsValue.ValueType != vm.Ref {
		panic("lhs value in assignment should have a value type Ref")
	}
	rhsValue := visitExpr(assignment.Rhs, localEnv)
	if err := checkAssignmentValidity(lhsValue.AsString(), localEnv); err != nil {
		panic(err)
	}
	store(lhsValue.AsString(), rhsValue, localEnv, Assign)
	return vm.Value{}
}

func visitDefDeclStmt(stmt syntax.Stmt, localEnv Environment) vm.Value {
	var (
		defDeclStmt = stmt.(*syntax.DefDeclStmt)
		returnType  = visitExpr(defDeclStmt.ReturnType, localEnv)
		defValue    = &vm.DefValue{
			DefDeclStmt: defDeclStmt,
			ReturnType:  vm.MiniscalaTypeToValueType(returnType.AsString()),
		}
		value = vm.Value{
			ValueType: vm.Function,
		}
	)

	value.Value = defValue

	// functions reside in global environment exclusively
	store(defDeclStmt.Name.Value, value, nil, Declare)

	return vm.Value{}
}
