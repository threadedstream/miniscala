package backing

import (
	"github.com/ThreadedStream/miniscala/assert"
	"github.com/ThreadedStream/miniscala/syntax"
)

type ExecutionContext int

const (
	TreeWalkInterpreter ExecutionContext = iota
	Vm
)

type (
	Value struct {
		Value     interface{}
		ValueType ValueType
		Returned  bool
	}

	// DefValue or Callable backing
	DefValue struct {
		DefDeclStmt *syntax.DefDeclStmt
		ReturnType  ValueType
	}
)

func NullValue() Value {
	return Value{
		ValueType: Null,
	}
}

func (v Value) AsFloat() float64 {
	assert.Assert(v.IsFloat(), "cannot cast value type %s to float", ValueTypeToStr(v.ValueType))
	return v.Value.(float64)
}

func (v Value) AsInt() int64 {
	assert.Assert(v.IsInt(), "cannot cast value type %s to int", ValueTypeToStr(v.ValueType))
	return v.Value.(int64)
}

func (v Value) AsString() string {
	assert.Assert(v.IsString() || v.IsRef(), "cannot cast value type %s to string", ValueTypeToStr(v.ValueType))
	return v.Value.(string)
}

func (v Value) AsBool() bool {
	assert.Assert(v.IsBool(), "cannot cast value type %s to bool", ValueTypeToStr(v.ValueType))
	return v.Value.(bool)
}

func (v Value) AsFunction() *DefValue {
	assert.Assert(v.IsFunction(), "cannot cast value type %s to function", ValueTypeToStr(v.ValueType))
	return v.Value.(*DefValue)
}

// change to v.ValueType == Float
func (v Value) IsFloat() bool {
	_, ok := v.Value.(float64)
	return ok
}

func (v Value) IsInt() bool {
	//_, ok := v.Value.(int)
	//return ok
	return v.ValueType == Int
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

func Add(v1, v2 Value, localEnv ValueEnv, ctx ExecutionContext) Value {
	if ctx == TreeWalkInterpreter {
		if v1.ValueType == Ref {
			v1, _ = LookupValue(v1.AsString(), localEnv, true)
		}
		if v2.ValueType == Ref {
			v2, _ = LookupValue(v2.AsString(), localEnv, true)
		}
	}

	switch {
	default:
		return Value{
			Value:     nil,
			ValueType: Undefined,
		}
	case v1.IsFloat() && v2.IsFloat():
		return Value{
			Value:     v1.AsFloat() + v2.AsFloat(),
			ValueType: Float,
		}
	case v1.IsString() && v2.IsString():
		return Value{
			Value:     v1.AsString() + v2.AsString(),
			ValueType: String,
		}
	case v1.IsInt() && v2.IsInt():
		return Value{
			Value:     v1.AsInt() + v2.AsInt(),
			ValueType: Int,
		}
	case v1.IsFloat() && v2.IsInt():
		return Value{
			Value:     v1.AsFloat() + float64(v2.AsInt()),
			ValueType: Float,
		}
	case v1.IsInt() && v2.IsFloat():
		return Value{
			Value:     float64(v1.AsInt()) + v2.AsFloat(),
			ValueType: Float,
		}
	}
}

func Sub(v1, v2 Value, localEnv ValueEnv, ctx ExecutionContext) Value {
	if ctx == TreeWalkInterpreter {
		if v1.ValueType == Ref {
			v1, _ = LookupValue(v1.AsString(), localEnv, true)
		}
		if v2.ValueType == Ref {
			v2, _ = LookupValue(v2.AsString(), localEnv, true)
		}
	}

	switch {
	default:
		return Value{
			Value:     nil,
			ValueType: Undefined,
		}
	case v1.IsFloat() && v2.IsFloat():
		return Value{
			Value:     v1.AsFloat() - v2.AsFloat(),
			ValueType: Float,
		}
	case v1.IsInt() && v2.IsInt():
		return Value{
			Value:     v1.AsInt() - v2.AsInt(),
			ValueType: Int,
		}
	case v1.IsFloat() && v2.IsInt():
		return Value{
			Value:     v1.AsFloat() - float64(v2.AsInt()),
			ValueType: Float,
		}
	case v1.IsInt() && v2.IsFloat():
		return Value{
			Value:     float64(v1.AsInt()) - v2.AsFloat(),
			ValueType: Float,
		}
	}
}

func Mul(v1, v2 Value, localEnv ValueEnv, ctx ExecutionContext) Value {
	if ctx == TreeWalkInterpreter {
		if v1.ValueType == Ref {
			v1, _ = LookupValue(v1.AsString(), localEnv, true)
		}
		if v2.ValueType == Ref {
			v2, _ = LookupValue(v2.AsString(), localEnv, true)
		}
	}

	switch {
	default:
		return Value{
			Value:     nil,
			ValueType: Undefined,
		}
	case v1.IsFloat() && v2.IsFloat():
		return Value{
			Value:     v1.AsFloat() * v2.AsFloat(),
			ValueType: Float,
		}
	case v1.IsInt() && v2.IsInt():
		return Value{
			Value:     v1.AsInt() * v2.AsInt(),
			ValueType: Int,
		}
	case v1.IsFloat() && v2.IsInt():
		return Value{
			Value:     v1.AsFloat() * float64(v2.AsInt()),
			ValueType: Float,
		}
	case v1.IsInt() && v2.IsFloat():
		return Value{
			Value:     float64(v1.AsInt()) * v2.AsFloat(),
			ValueType: Float,
		}
	}
}

func Div(v1, v2 Value, localEnv ValueEnv, ctx ExecutionContext) Value {
	if ctx == TreeWalkInterpreter {
		if v1.ValueType == Ref {
			v1, _ = LookupValue(v1.AsString(), localEnv, true)
		}
		if v2.ValueType == Ref {
			v2, _ = LookupValue(v2.AsString(), localEnv, true)
		}
	}

	switch {
	default:
		return Value{
			Value:     nil,
			ValueType: Undefined,
		}
	case v1.IsFloat() && v2.IsFloat():
		return Value{
			Value:     v1.AsFloat() / v2.AsFloat(),
			ValueType: Float,
		}
	case v1.IsInt() && v2.IsInt():
		return Value{
			Value:     float64(v1.AsInt()) / float64(v2.AsInt()),
			ValueType: Float,
		}
	case v1.IsFloat() && v2.IsInt():
		return Value{
			Value:     v1.AsFloat() / float64(v2.AsInt()),
			ValueType: Float,
		}
	case v1.IsInt() && v2.IsFloat():
		return Value{
			Value:     float64(v1.AsInt()) / v2.AsFloat(),
			ValueType: Float,
		}
	}
}

func Mod(v1, v2 Value, localEnv ValueEnv, ctx ExecutionContext) Value {
	switch {
	default:
		return Value{
			Value:     nil,
			ValueType: Undefined,
		}
	case v1.IsInt() && v2.IsInt():
		return Value{
			Value:     v1.AsInt() % v2.AsInt(),
			ValueType: Int,
		}
	}
}
