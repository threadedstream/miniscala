package backing

import "github.com/ThreadedStream/miniscala/syntax"

type ValueType int

const (
	Float ValueType = iota
	Int
	String
	Unit
	Bool
	Function
	Ref
	Any
	Null
	Undefined
)

type TypeInfo struct {
	ValueType
	Immutable  bool
	ParamTypes []ValueType // for functions
}

func MiniscalaTypeToValueType(typ string) ValueType {
	switch typ {
	default:
		return Undefined
	case "Float":
		return Float
	case "Int":
		return Int
	case "String":
		return String
	case "Unit":
		return Unit
	case "Bool":
		return Bool
	}
}

func LitKindToValueType(kind syntax.LitKind) ValueType {
	switch kind {
	default:
		return Undefined
	case syntax.FloatLit:
		return Float
	case syntax.IntLit:
		return Int
	case syntax.StringLit:
		return String
	case syntax.BoolLit:
		return Bool
	}
}

func TypesEqual(t1, t2 ValueType) bool {
	if t1 != Any && t2 != Any {
		return t1 == t2
	}
	// ValueType::Any == (whatsoever ValueType was passed) evaluates to true
	// Although I've got to make sure this scheme is suitable for each possible corner-case
	return true
}

func ValueTypeToStr(valueType ValueType) string {
	switch valueType {
	default:
		return "Unknown"
	case Float:
		return "Float"
	case Int:
		return "Int"
	case String:
		return "String"
	case Bool:
		return "Bool"
	case Function:
		return "Function"
	case Ref:
		return "Ref"
	case Null:
		return "Null"
	case Undefined:
		return "Undefined"
	}
}
