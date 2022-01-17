package interpreter

import (
	"github.com/ThreadedStream/miniscala/assert"
	"github.com/ThreadedStream/miniscala/syntax"
)

type ValueType int

const (
	Float ValueType = iota
	String
	Unit // void
	Bool // add functionality for bool later on
	Function
	Ref
	Null
	Undefined
)

type (
	Value struct {
		Value     interface{}
		ValueType ValueType
		Immutable bool
		Returned  bool
	}

	// DefValue or Callable value
	DefValue struct {
		DefDeclStmt *syntax.DefDeclStmt
		ReturnType  ValueType
	}
)

func miniscalaTypeToValueType(typ string) ValueType {
	switch typ {
	default:
		return Undefined
	case "Int", "Float":
		return Float
	case "String":
		return String
	case "Unit":
		return Unit
	case "Bool":
		return Bool
	}
}

func (v Value) valueTypeToStr() string {
	switch v.ValueType {
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

func nullValue() Value {
	return Value{
		ValueType: Null,
	}
}

func (v Value) asFloat() float64 {
	assert.Assert(v.isFloat(), "cannot cast value type %s to float", v.valueTypeToStr())
	return v.Value.(float64)
}

func (v Value) asString() string {
	assert.Assert(v.isString() || v.isRef(), "cannot cast value type %s to string", v.valueTypeToStr())
	return v.Value.(string)
}

func (v Value) asFunction() *DefValue {
	assert.Assert(v.isFunction(), "cannot cast value type %s to function", v.valueTypeToStr())
	return v.Value.(*DefValue)
}

// change to v.ValueType == Float
func (v Value) isFloat() bool {
	_, ok := v.Value.(float64)
	return ok
}

// thought it might be worthwhile putting it here
func (v Value) isZero() bool {
	assert.Assert(v.isFloat(), "cannot call isZero() on something other than float")
	return v.asFloat() == 0
}

// change to v.ValueType == String
func (v Value) isString() bool {
	_, ok := v.Value.(string)
	return ok
}

func (v Value) isRef() bool {
	return v.ValueType == Ref
}

func (v Value) isBool() bool {
	return v.ValueType == Bool
}

func (v Value) isFunction() bool {
	return v.ValueType == Function
}

func (v Value) isNull() bool {
	return v.ValueType == Null
}

func (v Value) isUndefined() bool {
	return v.ValueType == Undefined
}

func add(v1, v2 Value, localEnv Environment) Value {
	if v1.ValueType == Ref {
		v1, _ = lookup(v1.asString(), localEnv, true)
	}
	if v2.ValueType == Ref {
		v2, _ = lookup(v2.asString(), localEnv, true)
	}

	switch v1.Value.(type) {
	default:
		return Value{
			Value:     nil,
			ValueType: Undefined,
		}
	case float64:
		return Value{
			Value:     v1.asFloat() + v2.asFloat(),
			ValueType: Float,
		}
	case string:
		return Value{
			Value:     v1.asString() + v2.asString(),
			ValueType: String,
		}
	}
}

func sub(v1, v2 Value, localEnv Environment) Value {
	if v1.ValueType == Ref {
		v1, _ = lookup(v1.asString(), localEnv, true)
	}
	if v2.ValueType == Ref {
		v2, _ = lookup(v2.asString(), localEnv, true)
	}

	switch v1.Value.(type) {
	case float64:
		return Value{
			Value:     v1.asFloat() - v2.asFloat(),
			ValueType: Float,
		}
	default:
		return Value{
			Value:     nil,
			ValueType: Undefined,
		}
	}
}

func mul(v1, v2 Value, localEnv Environment) Value {
	if v1.ValueType == Ref {
		v1, _ = lookup(v1.asString(), localEnv, true)
	}
	if v2.ValueType == Ref {
		v2, _ = lookup(v2.asString(), localEnv, true)
	}

	switch v1.Value.(type) {
	default:
		return Value{
			Value:     nil,
			ValueType: Undefined,
		}
	case float64:
		return Value{
			Value:     v1.asFloat() * v2.asFloat(),
			ValueType: Float,
		}
	}
}

func div(v1, v2 Value, localEnv Environment) Value {
	if v1.ValueType == Ref {
		v1, _ = lookup(v1.asString(), localEnv, true)
	}

	if v2.ValueType == Ref {
		v2, _ = lookup(v2.asString(), localEnv, true)
	}

	switch v1.Value.(type) {
	default:
		return Value{
			Value:     nil,
			ValueType: Undefined,
		}
	case float64:
		return Value{
			Value:     v1.asFloat() / v2.asFloat(),
			ValueType: Float,
		}
	}
}
