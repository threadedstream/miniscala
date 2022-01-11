package main

import (
	"fmt"
	"reflect"
	"strconv"
)

func isCondTrue(cond Operation) bool {
	isComparisonOp(cond.op)
	var (
		lhs = visitExpr(cond.lhs)
		rhs = visitExpr(cond.rhs)
	)

	lhs = resolveRef(lhs)
	rhs = resolveRef(rhs)

	switch cond.op {
	default:
		return false
	case GreaterThan:
		if lhs.isString() && rhs.isString() {
			return lhs.asString() > rhs.asString()
		} else if lhs.isFloat() && rhs.isFloat() {
			return lhs.asFloat() > rhs.asFloat()
		} else {
			panic("cast to string and float was unsuccessful")
		}
	case GreaterThanOrEqual:
		if lhs.isString() && rhs.isString() {
			return lhs.asString() >= rhs.asString()
		} else if lhs.isFloat() && rhs.isFloat() {
			return lhs.asFloat() >= rhs.asFloat()
		} else {
			panic("cast to string and float was unsuccessful")
		}
	case LessThan:
		if lhs.isString() && rhs.isString() {
			return lhs.asString() < rhs.asString()
		} else if lhs.isFloat() && rhs.isFloat() {
			return lhs.asFloat() <= rhs.asFloat()
		} else {
			panic("cast to string and float was unsuccessful")
		}
	case LessThanOrEqual:
		if lhs.isString() && rhs.isString() {
			return lhs.asString() < rhs.asString()
		} else if lhs.isFloat() && rhs.isFloat() {
			return lhs.asFloat() <= rhs.asFloat()
		} else {
			panic("cast to string and float was unsuccessful")
		}
	case Equal:
		return lhs.value == rhs.value
	case NotEqual:
		return lhs.value != rhs.value
	}
}

func add(v1, v2 Value) Value {
	if v1.valueType == Ref {
		v1 = lookup(v1.asString())
	}
	if v2.valueType == Ref {
		v2 = lookup(v2.asString())
	}

	switch v1.value.(type) {
	default:
		return Value{
			value:     nil,
			valueType: Undefined,
		}
	case float64:
		return Value{
			value:     v1.asFloat() + v2.asFloat(),
			valueType: Float,
		}
	case string:
		return Value{
			value:     v1.asString() + v2.asString(),
			valueType: String,
		}
	}
}

func sub(v1, v2 Value) Value {
	if v1.valueType == Ref {
		v1 = lookup(v1.asString())
	}
	if v2.valueType == Ref {
		v2 = lookup(v2.asString())
	}

	switch v1.value.(type) {
	case float64:
		return Value{
			value:     v1.asFloat() - v2.asFloat(),
			valueType: Float,
		}
	default:
		return Value{
			value:     nil,
			valueType: Undefined,
		}
	}
}

func mul(v1, v2 Value) Value {
	if v1.valueType == Ref {
		v1 = lookup(v1.asString())
	}
	if v2.valueType == Ref {
		v2 = lookup(v2.asString())
	}

	switch v1.value.(type) {
	case float64:
		return Value{
			value:     v1.asFloat() * v2.asFloat(),
			valueType: Float,
		}
	default:
		return Value{
			value:     nil,
			valueType: Undefined,
		}
	}
}

func execute(program Program) {
	for _, stmt := range program.stmtList {
		visitStmt(stmt)
	}
}

func visitStmt(node Stmt) Value {
	switch node.(type) {
	default:
		panic(fmt.Errorf("unknown node type %v", reflect.TypeOf(node)))
	case *VarDeclStmt:
		varDecl := node.(*VarDeclStmt)
		value := visitExpr(varDecl.rhs)
		value.immutable = false
		store(varDecl.name.value, value)
		return Value{}
	case *ValDeclStmt:
		valDecl := node.(*ValDeclStmt)
		value := visitExpr(valDecl.rhs)
		value.immutable = true
		store(valDecl.name.value, value)
		return Value{}
	case *IfStmt:
		var ifStmt = node.(*IfStmt)
		if isCondTrue(ifStmt.cond) {
			visitStmt(ifStmt.body)
		} else {
			visitStmt(ifStmt.elseBody)
		}
		return Value{}
	case *WhileStmt:
		var whileStmt = node.(*WhileStmt)
		for isCondTrue(whileStmt.cond) {
			visitStmt(whileStmt.body)
		}
		return Value{}
	case *Assignment:
		assignment := node.(*Assignment)
		lhsValue := visitExpr(assignment.lhs)
		if lhsValue.valueType != Ref {
			panic("lhs value in assignment should have a value type Ref")
		}
		rhsValue := visitExpr(assignment.rhs)
		if err := checkAssignmentValidity(lhsValue.asString()); err != nil {
			panic(err)
		}
		store(lhsValue.asString(), rhsValue)
		return Value{}
	}
}

func visitExpr(expr Expr) Value {
	switch expr.(type) {
	default:
		return Value{value: nil}
	case *Name:
		name := expr.(*Name)
		v := Value{
			value:     name.value,
			valueType: Ref,
		}
		return v
	case *BasicLit:
		basicLit := expr.(*BasicLit)
		v := Value{}
		switch basicLit.kind {
		case FloatLit:
			v.value, _ = strconv.ParseFloat(basicLit.value, 64)
			v.valueType = Float
		case StringLit:
			v.value = basicLit.value
			v.valueType = String
		}
		return v
	case *Operation:
		operation := expr.(*Operation)
		switch operation.op {
		default:
			panic("unknown operation")
		case PlusOp:
			lhsValue := visitExpr(operation.lhs)
			rhsValue := visitExpr(operation.rhs)
			checkOpValues(PlusOp, lhsValue, rhsValue)
			value := add(lhsValue, rhsValue)
			return value
		case MinusOp:
			lhsValue := visitExpr(operation.lhs)
			rhsValue := visitExpr(operation.rhs)
			checkOpValues(MinusOp, lhsValue, rhsValue)
			value := sub(lhsValue, rhsValue)
			return value
		case MulOp:
			lhsValue := visitExpr(operation.lhs)
			rhsValue := visitExpr(operation.rhs)
			checkOpValues(MulOp, lhsValue, rhsValue)
			value := mul(lhsValue, rhsValue)
			return value
		}
	}
}
