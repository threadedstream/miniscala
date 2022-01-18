package vm

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

func MiniscalaTypeToValueType(typ string) ValueType {
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

func (v Value) ValueTypeToStr() string {
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

func NullValue() Value {
	return Value{
		ValueType: Null,
	}
}

func (v Value) AsFloat() float64 {
	assert.Assert(v.IsFloat(), "cannot cast value type %s to float", v.ValueTypeToStr())
	return v.Value.(float64)
}

func (v Value) AsString() string {
	assert.Assert(v.IsString() || v.IsRef(), "cannot cast value type %s to string", v.ValueTypeToStr())
	return v.Value.(string)
}

func (v Value) AsBool() bool {
	assert.Assert(v.IsBool(), "cannot cast value type %s to bool", v.ValueTypeToStr())
	return v.Value.(bool)
}

func (v Value) AsFunction() *DefValue {
	assert.Assert(v.IsFunction(), "cannot cast value type %s to function", v.ValueTypeToStr())
	return v.Value.(*DefValue)
}

// change to v.ValueType == Float
func (v Value) IsFloat() bool {
	_, ok := v.Value.(float64)
	return ok
}

// thought it might be worthwhile putting it here
func (v Value) IsZero() bool {
	assert.Assert(v.IsFloat(), "cannot call isZero() on something other than float")
	return v.AsFloat() == 0
}

// change to v.ValueType == String
func (v Value) IsString() bool {
	_, ok := v.Value.(string)
	return ok
}

func (v Value) IsRef() bool {
	return v.ValueType == Ref
}

func (v Value) IsBool() bool {
	return v.ValueType == Bool
}

func (v Value) IsFunction() bool {
	return v.ValueType == Function
}

func (v Value) IsNull() bool {
	return v.ValueType == Null
}

func (v Value) IsUndefined() bool {
	return v.ValueType == Undefined
}
