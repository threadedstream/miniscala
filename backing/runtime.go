package backing

import (
	"fmt"
	"github.com/ThreadedStream/miniscala/assert"
)

func IsRuntimeCall(name string) bool {
	switch name {
	default:
		return false
	case "print", "to_string":
		return true
	}
}

func DispatchRuntimeFuncCall(name string, args ...Value) {
	switch name {
	case "print":
		callPrint(args[0])
	case "to_string":
		callToString(args[0])
	}
}

func callPrint(val Value) {
	assert.Assert(val.ValueType == String, "print requires string type as an only argument")
	fmt.Printf("%s", val.Value)
}

func callToString(val Value) Value {
	strValue := fmt.Sprintf("%v", val.Value)
	return Value{
		Value:     strValue,
		ValueType: String,
	}
}
