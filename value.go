package main

type ValueType int

const (
	Literal ValueType = iota
	Ref
	Null
	Undefined
)

type Value struct {
	value     interface{}
	valueType ValueType
	immutable bool
}

func (v *Value) asFloat() float64 {
	return v.value.(float64)
}

func (v *Value) asString() string {
	return v.value.(string)
}

func (v *Value) isFloat() bool {
	_, ok := v.value.(float64)
	return ok
}

func (v *Value) isString() bool {
	_, ok := v.value.(string)
	return ok
}

func (v *Value) isNull() bool {
	return v.valueType == Null
}

func (v *Value) isUndefined() bool {
	return v.valueType == Undefined
}
