package main

import "fmt"

type ValueType int

const (
	Float ValueType = iota
	String
	Bool // add functionality for bool later on
	Function
	Ref
	Null
	Undefined
)

type Value struct {
	value     interface{}
	valueType ValueType
	immutable bool
}

func (v Value) valueTypeToStr() string {
	switch v.valueType {
	default:
		return "VtUnknown"
	case Float:
		return "VtFloat"
	case String:
		return "VtString"
	case Bool:
		return "VtBool"
	case Function:
		return "VtFunction"
	case Ref:
		return "VtRef"
	case Null:
		return "VtNull"
	case Undefined:
		return "Undefined"
	}
}

func (v Value) asFloat() float64 {
	assert(v.valueType == Float, func() {
		panic(fmt.Errorf("cannot cast value type %s to float", v.valueTypeToStr()))
	})
	return v.value.(float64)
}

func (v Value) asString() string {
	assert(v.valueType == String, func() {
		panic(fmt.Errorf("cannot cast value type %s to string", v.valueTypeToStr()))
	})
	return v.value.(string)
}

func (v Value) isFloat() bool {
	_, ok := v.value.(float64)
	return ok
}

func (v Value) isString() bool {
	_, ok := v.value.(string)
	return ok
}

func (v Value) isBool() bool {
	return v.valueType == Bool
}

func (v Value) isFunction() bool {
	return v.valueType == Function
}

func (v Value) isNull() bool {
	return v.valueType == Null
}

func (v Value) isUndefined() bool {
	return v.valueType == Undefined
}
