package interpreter

import (
	"errors"
	"fmt"
	"github.com/ThreadedStream/miniscala/assert"
	"github.com/ThreadedStream/miniscala/backing"
	"github.com/ThreadedStream/miniscala/syntax"
	"strconv"
)

func checkOpValues(op syntax.Operator, v1, v2 backing.Value, localEnv backing.Environment) {

	if v1.ValueType == backing.Ref {
		v1, _ = backing.Lookup(v1.AsString(), localEnv, true)
	}

	if v2.ValueType == backing.Ref {
		v2, _ = backing.Lookup(v2.AsString(), localEnv, true)
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

func checkDefReturnType(v backing.Value, defReturnType, funcName string) {
	if v.ValueType == backing.MiniscalaTypeToValueType(defReturnType) {
		return
	}

	panic(fmt.Errorf("(in %s()) expected to return a backing of type %s, but it returned %s", funcName, defReturnType, v.ValueTypeToStr()))
}

func isReservedFuncCall(funcName string) bool {
	switch funcName {
	default:
		return false
	case "print":
		return true
	}
}

func unwrapValue(value backing.Value) interface{} {
	switch {
	default:
		return nil
	case value.IsString():
		return value.AsString()
	case value.IsFloat():
		return value.AsFloat()
	}
}

func callPrintFunc(call *syntax.Call, localEnv backing.Environment) {
	// TODO(threadedstream): allow more arguments to print
	assert.Assert(len(call.ArgList) == 1, "expected a single argument for print got %d", len(call.ArgList))
	value := unwrapValue(visitExpr(call.ArgList[0], localEnv))
	fmt.Printf("%v", value)
}

func dispatchReservedCall(call *syntax.Call, localEnv backing.Environment) {
	switch call.CalleeName.Value {
	case "print":
		callPrintFunc(call, localEnv)
	}
}

func checkAssignmentValidity(name string, localEnv backing.Environment) error {
	// first, check against the presence of backing associated with name
	value, _ := backing.Lookup(name, localEnv, true)

	// second, check against the possibility to change this backing
	if value.Immutable {
		return errors.New("attempt to override val backing")
	}

	return nil
}

func resolveRef(v1 backing.Value, localEnv backing.Environment) backing.Value {
	switch {
	default:
		return v1
	case v1.ValueType == backing.Ref:
		v, _ := backing.Lookup(v1.AsString(), localEnv, true)
		return v
	}
}

func isCondTrue(cond syntax.Operation, localEnv backing.Environment) bool {
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
	//environment.state()
}

func visitStmt(stmt syntax.Stmt, localEnv backing.Environment) backing.Value {
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

func visitExpr(expr syntax.Expr, localEnv backing.Environment) backing.Value {
	switch expr.(type) {
	default:
		return backing.Value{Value: nil}
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
func visitName(expr syntax.Expr) backing.Value {
	name := expr.(*syntax.Name)
	v := backing.Value{
		Value:     name.Value,
		ValueType: backing.Ref,
	}
	return v
}

func visitBasicLit(expr syntax.Expr) backing.Value {
	basicLit := expr.(*syntax.BasicLit)
	v := backing.Value{}
	switch basicLit.Kind {
	case syntax.FloatLit:
		v.Value, _ = strconv.ParseFloat(basicLit.Value, 64)
		v.ValueType = backing.Float
	case syntax.StringLit:
		v.Value = basicLit.Value
		v.ValueType = backing.String
	}
	return v
}

func visitOperation(expr syntax.Expr, localEnv backing.Environment) backing.Value {
	operation := expr.(*syntax.Operation)
	switch operation.Op {
	default:
		panic("unknown operation")
	case syntax.Plus:
		lhsValue := visitExpr(operation.Lhs, localEnv)
		rhsValue := visitExpr(operation.Rhs, localEnv)
		checkOpValues(syntax.Plus, lhsValue, rhsValue, localEnv)
		value := backing.Add(lhsValue, rhsValue, localEnv, backing.TreeWalkInterpreter)
		return value
	case syntax.Minus:
		lhsValue := visitExpr(operation.Lhs, localEnv)
		rhsValue := visitExpr(operation.Rhs, localEnv)
		checkOpValues(syntax.Minus, lhsValue, rhsValue, localEnv)
		value := backing.Sub(lhsValue, rhsValue, localEnv, backing.TreeWalkInterpreter)
		return value
	case syntax.Mul:
		lhsValue := visitExpr(operation.Lhs, localEnv)
		rhsValue := visitExpr(operation.Rhs, localEnv)
		checkOpValues(syntax.Mul, lhsValue, rhsValue, localEnv)
		value := backing.Mul(lhsValue, rhsValue, localEnv, backing.TreeWalkInterpreter)
		return value
	case syntax.Div:
		lhsValue := visitExpr(operation.Lhs, localEnv)
		rhsValue := visitExpr(operation.Rhs, localEnv)
		checkOpValues(syntax.Div, lhsValue, rhsValue, localEnv)
		value := backing.Div(lhsValue, rhsValue, localEnv, backing.TreeWalkInterpreter)
		return value
	}
}

func visitCall(expr syntax.Expr, localEnv backing.Environment) backing.Value {
	call := expr.(*syntax.Call)

	if isReservedFuncCall(call.CalleeName.Value) {
		// dispatch in case if call to a reserved function has been made
		dispatchReservedCall(call, localEnv)
		return backing.Value{}
	}

	value, _ := backing.Lookup(call.CalleeName.Value, nil, true)
	defValue := value.AsFunction()

	var funcFrame = make(backing.Environment)
	// TODO(threadedstream): do some checks regarding the number of passed arguments
	for idx, arg := range call.ArgList {
		argValue := visitExpr(arg, localEnv)
		paramName := defValue.DefDeclStmt.ParamList[idx].Name.Value
		backing.Store(paramName, argValue, funcFrame, backing.Assign)
	}

	returnValue := visitBlockStmt(defValue.DefDeclStmt.Body, funcFrame)

	checkDefReturnType(returnValue, defValue.DefDeclStmt.ReturnType.(*syntax.Name).Value, defValue.DefDeclStmt.Name.Value)

	// destroy local environment
	funcFrame = nil

	return returnValue
}

func visitReturnStmt(stmt syntax.Stmt, localEnv backing.Environment) backing.Value {
	returnStmt := stmt.(*syntax.ReturnStmt)
	returnValue := visitExpr(returnStmt.Value, localEnv)
	returnValue.Returned = true
	return returnValue
}

// statements
func visitVarDeclStmt(stmt syntax.Stmt, localEnv backing.Environment) backing.Value {
	varDecl := stmt.(*syntax.VarDeclStmt)
	value := visitExpr(varDecl.Rhs, localEnv)
	value.Immutable = false
	backing.Store(varDecl.Name.Value, value, localEnv, backing.Declare)
	return backing.Value{}
}

func visitValDeclStmt(stmt syntax.Stmt, localEnv backing.Environment) backing.Value {
	valDecl := stmt.(*syntax.ValDeclStmt)
	value := visitExpr(valDecl.Rhs, localEnv)
	value.Immutable = true
	backing.Store(valDecl.Name.Value, value, localEnv, backing.Declare)
	return backing.Value{}
}

func visitBlockStmt(stmt syntax.Stmt, localEnv backing.Environment) backing.Value {
	var (
		block = stmt.(*syntax.BlockStmt)
		value backing.Value
	)

	for _, currStmt := range block.Stmts {
		value = visitStmt(currStmt, localEnv)
		if value.Returned {
			break
		}
	}

	return value
}

func visitIfStmt(stmt syntax.Stmt, localEnv backing.Environment) backing.Value {
	var (
		value  = backing.Value{ValueType: backing.Null}
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

func visitWhileStmt(stmt syntax.Stmt, localEnv backing.Environment) backing.Value {
	var whileStmt = stmt.(*syntax.WhileStmt)
	for isCondTrue(whileStmt.Cond, localEnv) {
		visitStmt(whileStmt.Body, localEnv)
	}
	return backing.Value{}
}

func visitAssignment(stmt syntax.Stmt, localEnv backing.Environment) backing.Value {
	assignment := stmt.(*syntax.Assignment)
	lhsValue := visitExpr(assignment.Lhs, localEnv)
	if lhsValue.ValueType != backing.Ref {
		panic("lhs backing in assignment should have a backing type Ref")
	}
	rhsValue := visitExpr(assignment.Rhs, localEnv)
	if err := checkAssignmentValidity(lhsValue.AsString(), localEnv); err != nil {
		panic(err)
	}
	backing.Store(lhsValue.AsString(), rhsValue, localEnv, backing.Assign)
	return backing.Value{}
}

func visitDefDeclStmt(stmt syntax.Stmt, localEnv backing.Environment) backing.Value {
	var (
		defDeclStmt = stmt.(*syntax.DefDeclStmt)
		returnType  = visitExpr(defDeclStmt.ReturnType, localEnv)
		defValue    = &backing.DefValue{
			DefDeclStmt: defDeclStmt,
			ReturnType:  backing.MiniscalaTypeToValueType(returnType.AsString()),
		}
		value = backing.Value{
			ValueType: backing.Function,
		}
	)

	value.Value = defValue

	// functions reside in global environment exclusively
	backing.Store(defDeclStmt.Name.Value, value, nil, backing.Declare)

	return backing.Value{}
}
