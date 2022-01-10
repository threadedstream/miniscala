package main

import (
	"fmt"
	"strconv"
)

var (
	environment = make(map[string]Value)
)

func add(v1, v2 Value) Value {
	var (
		realValue1, realValue2 Value
		ok                     bool
	)

	if v1.valueType == Ref {
		realValue1, ok = environment[v1.asString()]
		if !ok {
			panic(fmt.Errorf("no entry with name %s was found in a global object map", v1.value.(string)))
		}
	} else {
		realValue1 = v1
	}

	if v2.valueType == Ref {
		realValue2, ok = environment[v2.asString()]
		if !ok {
			panic(fmt.Errorf("no entry with name %s was found in a global object map", v2.value.(string)))
		}
	} else {
		realValue2 = v2
	}

	switch realValue1.value.(type) {
	case float64:
		return Value{
			value:     realValue1.asFloat() + realValue2.asFloat(),
			valueType: Literal,
		}
	case string:
		return Value{
			value:     realValue1.asString() + realValue2.asString(),
			valueType: Literal,
		}
	default:
		return Value{
			value:     nil,
			valueType: Undefined,
		}
	}
}

func sub(v1, v2 Value) Value {
	var (
		realValue1, realValue2 Value
		ok                     bool
	)

	if v1.valueType == Ref {
		realValue1, ok = environment[v1.asString()]
		if !ok {
			panic(fmt.Errorf("no entry with name %s was found in a global object map", v1.asString()))
		}
	} else {
		realValue1 = v1
	}

	if v2.valueType == Ref {
		realValue2, ok = environment[v2.asString()]
		if !ok {
			panic(fmt.Errorf("no entry with name %s was found in a global object map", v2.asString()))
		}
	} else {
		realValue2 = v2
	}

	switch realValue1.value.(type) {
	case float64:
		return Value{
			value:     realValue1.asFloat() - realValue2.asFloat(),
			valueType: Literal,
		}
	default:
		return Value{
			value:     nil,
			valueType: Undefined,
		}
	}
}

func mul(v1, v2 Value) Value {
	var (
		realValue1, realValue2 Value
		ok                     bool
	)

	if v1.valueType == Ref {
		realValue1, ok = environment[v1.asString()]
		if !ok {
			panic(fmt.Errorf("no entry with name %s was found in a global object map", v1.asString()))
		}
	} else {
		realValue1 = v1
	}

	if v2.valueType == Ref {
		realValue2, ok = environment[v2.asString()]
		if !ok {
			panic(fmt.Errorf("no entry with name %s was found in a global object map", v2.asString()))
		}
	} else {
		realValue2 = v2
	}

	switch realValue1.value.(type) {
	case float64:
		return Value{
			value:     realValue1.asFloat() * realValue2.asFloat(),
			valueType: Literal,
		}
	default:
		return Value{
			value:     nil,
			valueType: Undefined,
		}
	}
}

func execute(program Program) {
	for _, node := range program.nodeList {
		switch node.(type) {
		case *VarDecl, *ValDecl:
			visitDecl(node)
		}
	}
}

func visitDecl(node Node) {
	switch node.(type) {
	case *VarDecl:
		varDecl := node.(*VarDecl)
		value := visitExpr(varDecl.rhs)
		value.immutable = false
		environment[varDecl.name.value] = value
	case *ValDecl:
		valDecl := node.(*ValDecl)
		value := visitExpr(valDecl.rhs)
		value.immutable = true
		environment[valDecl.name.value] = value
	}
}

func visitExpr(expr Expr) Value {
	switch expr.(type) {
	case *Name:
		name := expr.(*Name)
		v := Value{
			value:     name.value,
			valueType: Ref,
		}
		return v
	case *BasicLit:
		basicLit := expr.(*BasicLit)
		v := Value{
			valueType: Literal,
		}
		switch basicLit.kind {
		case FloatLit:
			v.value, _ = strconv.ParseFloat(basicLit.value, 32)
		case StringLit:
			v.value = basicLit.value
		}
		return v
	case *IfStmt:
		// TODO(threadedstream):
		return Value{}
	case *WhileStmt:
		//TODO(threadedstream):
		return Value{}
	case *Assignment:
		assignment := expr.(*Assignment)
		lhsValue := visitExpr(assignment.rhs)
		if lhsValue.valueType != Ref {
			panic("lhs value in assignment should have a value type Ref")
		}
		rhsValue := visitExpr(assignment.lhs)
		if err := checkAssignmentValidity(lhsValue.value.(string)); err != nil {
			panic(err)
		}
		environment[lhsValue.value.(string)] = rhsValue
		// for now, return a dummy value object
		return Value{}
	case *Operation:
		operation := expr.(*Operation)
		switch operation.op {
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
		default:
			panic("unknown operation")
		}
	default:
		return Value{value: nil}
	}
}
