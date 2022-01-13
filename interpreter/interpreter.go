package interpreter

import (
	"errors"
	"fmt"
	"github.com/ThreadedStream/miniscala/syntax"
	"reflect"
	"strconv"
)

func checkOpValues(op syntax.Operator, v1, v2 Value) {

	if v1.ValueType == Ref {
		v1 = environment[v1.asString()]
	}

	if v2.ValueType == Ref {
		v2 = environment[v2.asString()]
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

func checkAssignmentValidity(name string) error {
	// first, check against the presence of value associated with name
	value, ok := environment[name]
	if !ok {
		return fmt.Errorf("no entry with name %s was found", name)
	}

	// second, check against the possibility to change this value
	if value.Immutable {
		return errors.New("attempt to override val value")
	}

	return nil
}

func resolveRef(v1 Value) Value {
	switch {
	default:
		return v1
	case v1.ValueType == Ref:
		return lookup(v1.asString())
	}
}

func isCondTrue(cond syntax.Operation) bool {
	syntax.IsComparisonOp(cond.Op)
	var (
		lhs = visitExpr(cond.Lhs)
		rhs = visitExpr(cond.Rhs)
	)

	lhs = resolveRef(lhs)
	rhs = resolveRef(rhs)

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
			return lhs.asFloat() <= rhs.asFloat()
		} else {
			panic("cast to string and float was unsuccessful")
		}
	case syntax.LessThanOrEqual:
		if lhs.isString() && rhs.isString() {
			return lhs.asString() < rhs.asString()
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
		visitStmt(stmt)
	}
}

func DumpEnvState() {
	state()
}

func visitStmt(stmt syntax.Stmt) Value {
	switch stmt.(type) {
	default:
		panic(fmt.Errorf("unknown node type %v", reflect.TypeOf(stmt)))
	case *syntax.VarDeclStmt:
		return visitVarDeclStmt(stmt)
	case *syntax.ValDeclStmt:
		return visitValDeclStmt(stmt)
	case *syntax.BlockStmt:
		return visitBlockStmt(stmt)
	case *syntax.IfStmt:
		return visitIfStmt(stmt)
	case *syntax.WhileStmt:
		return visitWhileStmt(stmt)
	case *syntax.Assignment:
		return visitAssignment(stmt)
	}
}

func visitExpr(expr syntax.Expr) Value {
	switch expr.(type) {
	default:
		return Value{Value: nil}
	case *syntax.Name:
		return visitName(expr)
	case *syntax.BasicLit:
		return visitBasicLit(expr)
	case *syntax.Operation:
		return visitOperation(expr)
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

func visitOperation(expr syntax.Expr) Value {
	operation := expr.(*syntax.Operation)
	switch operation.Op {
	default:
		panic("unknown operation")
	case syntax.Plus:
		lhsValue := visitExpr(operation.Lhs)
		rhsValue := visitExpr(operation.Rhs)
		checkOpValues(syntax.Plus, lhsValue, rhsValue)
		value := add(lhsValue, rhsValue)
		return value
	case syntax.Minus:
		lhsValue := visitExpr(operation.Lhs)
		rhsValue := visitExpr(operation.Rhs)
		checkOpValues(syntax.Minus, lhsValue, rhsValue)
		value := sub(lhsValue, rhsValue)
		return value
	case syntax.Mul:
		lhsValue := visitExpr(operation.Lhs)
		rhsValue := visitExpr(operation.Rhs)
		checkOpValues(syntax.Mul, lhsValue, rhsValue)
		value := mul(lhsValue, rhsValue)
		return value
	case syntax.Div:
		lhsValue := visitExpr(operation.Lhs)
		rhsValue := visitExpr(operation.Rhs)
		checkOpValues(syntax.Div, lhsValue, rhsValue)
		value := div(lhsValue, rhsValue)
		return value
	}
}

// statements
func visitVarDeclStmt(stmt syntax.Stmt) Value {
	varDecl := stmt.(*syntax.VarDeclStmt)
	value := visitExpr(varDecl.Rhs)
	value.Immutable = false
	store(varDecl.Name.Value, value)
	return Value{}
}

func visitValDeclStmt(stmt syntax.Stmt) Value {
	valDecl := stmt.(*syntax.ValDeclStmt)
	value := visitExpr(valDecl.Rhs)
	value.Immutable = true
	store(valDecl.Name.Value, value)
	return Value{}
}

func visitBlockStmt(stmt syntax.Stmt) Value {
	var block = stmt.(*syntax.BlockStmt)
	for _, stmt := range block.Stmts {
		visitStmt(stmt)
	}
	return Value{}
}

func visitIfStmt(stmt syntax.Stmt) Value {
	var ifStmt = stmt.(*syntax.IfStmt)
	if isCondTrue(ifStmt.Cond) {
		visitStmt(ifStmt.Body)
	} else {
		visitStmt(ifStmt.ElseBody)
	}
	return Value{}
}

func visitWhileStmt(stmt syntax.Stmt) Value {
	var whileStmt = stmt.(*syntax.WhileStmt)
	for isCondTrue(whileStmt.Cond) {
		visitStmt(whileStmt.Body)
	}
	return Value{}
}

func visitAssignment(stmt syntax.Stmt) Value {
	assignment := stmt.(*syntax.Assignment)
	lhsValue := visitExpr(assignment.Lhs)
	if lhsValue.ValueType != Ref {
		panic("lhs value in assignment should have a value type Ref")
	}
	rhsValue := visitExpr(assignment.Rhs)
	if err := checkAssignmentValidity(lhsValue.asString()); err != nil {
		panic(err)
	}
	store(lhsValue.asString(), rhsValue)
	return Value{}
}
