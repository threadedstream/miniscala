package backing

import (
	"fmt"
	"github.com/ThreadedStream/miniscala/assert"
)

func IsRuntimeCall(name string) bool {
	switch name {
	default:
		return false
	case "print", "to_string", "array_new", "array_set", "array_get", "array_size":
		return true
	}
}

func DispatchRuntimeFuncCall(name string, args ...Value) Value {
	switch name {
	case "print":
		callPrint(args[0])
	case "to_string":
		return callToString(args[0])
	case "array_new":
		return callArrayNew(args[0], args[1])
	case "array_set":
		callArraySet(args[0], args[1], args[2])
	case "array_get":
		return callArrayGet(args[0], args[1])
	case "array_size":
		return callArraySize(args[0])
	}
	return Value{
		ValueType: Unit,
	}
}

func callPrint(val Value) {
	assert.Assert(val.ValueType == String, "print requires string type as an only argument")
	fmt.Printf("%s", val.AsString())
}

func callToString(val Value) Value {
	strValue := fmt.Sprintf("%v", val.Value)
	return Value{
		Value:     strValue,
		ValueType: String,
	}
}

func callArrayNew(numberOfElements, typeOfElements Value) Value {
	assert.Assert(numberOfElements.IsInt(), "1st argument to array_new must be an integer")
	assert.Assert(typeOfElements.IsString(), "2nd argument to array_new must be a string")
	ty := MiniscalaTypeToValueType(typeOfElements.AsString())
	arrValue := ArrayValue{
		Arr:         ArrayOfValues(int(numberOfElements.AsInt()), ty),
		ElementType: ty,
	}
	return Value{
		Value:     arrValue,
		ValueType: Array,
	}
}

func callArraySize(arrPtr Value) Value {
	assert.Assert(arrPtr.IsArray(), "1st argument to array_size must be an array")
	arrValue := arrPtr.Value.(ArrayValue)
	return Value{
		Value:     int64(len(arrValue.Arr)),
		ValueType: Int,
	}
}

func callArraySet(arrPtr, idx, value Value) {
	assert.Assert(arrPtr.IsArray(), "1st argument to array_set must be an array")
	assert.Assert(idx.IsInt(), "2nd argument to array_set must be an integer")
	// TODO(threadedstream): do proper typechecking here, check accordance of type of the value in respect to
	// the type mandated by arrPtr
	arrValue := arrPtr.Value.(ArrayValue)
	assert.Assert(
		value.ValueType == arrValue.ElementType,
		"array expected type %s, but got %s",
		ValueTypeToStr(value.ValueType),
		ValueTypeToStr(arrValue.ElementType),
	)
	arrValue.Arr[idx.AsInt()] = value
}

func callArrayGet(arrPtr, idx Value) Value {
	assert.Assert(arrPtr.IsArray(), "1st argument to array_get must be an array")
	assert.Assert(idx.IsInt(), "2nd argument to array_get must be an integer")
	arrValue := arrPtr.Value.(ArrayValue)
	return arrValue.Arr[idx.AsInt()]
}
