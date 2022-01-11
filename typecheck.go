package main

import (
	"errors"
	"fmt"
)

func checkOpValues(op Operator, v1, v2 Value) {

	if v1.valueType == Ref {
		v1 = environment[v1.asString()]
	}

	if v2.valueType == Ref {
		v2 = environment[v2.asString()]
	}

	switch op {
	default:
		panic("unknown operation")
	case PlusOp:
		if (v1.isString() && v2.isString()) || (v1.isFloat() && v2.isFloat()) {
			return
		}
		panic("v1 and v2 must both be of type string or float")
	case MinusOp:
		if v1.isFloat() && v2.isFloat() {
			return
		}
		panic("v1 and v2 must both be of type float")
	case MulOp:
		if v1.isFloat() && v2.isFloat() {
			return
		}
		panic("v1 and v2 must both be of type float")
	}
}

func checkAssignmentValidity(name string) error {
	// first, check against the presence of value associated with name
	value, ok := environment[name]
	if !ok {
		return fmt.Errorf("no entry with name %s was found", name)
	}

	// second, check against the possibility to change this value
	if value.immutable {
		return errors.New("attempt to override val value")
	}

	return nil
}

func resolveRef(v1 Value) Value {
	switch {
	default:
		return v1
	case v1.valueType == Ref:
		return lookup(v1.asString())
	}
}
