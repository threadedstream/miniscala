package backing

import (
	"fmt"
)

type (
	ValueEnv map[string]Value
	TypeEnv  map[string]TypeInfo
)

type StoringContext int

const (
	Assign StoringContext = iota
	Declare
)

type LookupContext int

const (
	Val  LookupContext = iota
	Type LookupContext = iota
)

var (
	valueEnv = make(ValueEnv)
	typeEnv  = make(TypeEnv)
)

func StoreValue(name string, value Value, localValEnv ValueEnv, localTypeEnv TypeEnv, ctx StoringContext) {
	if ctx == Declare {
		_, ok := LookupValue(name, localValEnv, false)
		if ok {
			panic(fmt.Errorf("name %s has already got entry in a value env", name))
		}
	} else if ctx == Assign {
		_, ok := LookupValue(name, localValEnv, false)
		if !ok {
			panic(fmt.Errorf("undefined reference to name %s", name))
		}
		if localTypeEnv != nil {
			valTypeInfo, _ := LookupType(name, localTypeEnv, true)
			if valTypeInfo.Immutable {
				panic(fmt.Errorf("attempt to assign to immutable memory cell"))
			}
		}
	}

	switch {
	case localValEnv != nil:
		localValEnv[name] = value
	default:
		valueEnv[name] = value
	}
}

func StoreType(name string, valueType ValueType, immutable bool, paramTypes []ValueType, localEnv TypeEnv) {
	if localEnv != nil {
		localEnv[name] = TypeInfo{
			Immutable:  immutable,
			ValueType:  valueType,
			ParamTypes: paramTypes,
		}
	} else {
		typeEnv[name] = TypeInfo{
			Immutable:  immutable,
			ValueType:  valueType,
			ParamTypes: paramTypes,
		}
	}
}

func LookupType(name string, localEnv TypeEnv, shouldPanic bool) (TypeInfo, bool) {
	typ, ok := lookup(name, localEnv, shouldPanic, Type)
	return typ.(TypeInfo), ok
}

func LookupValue(name string, localEnv ValueEnv, shouldPanic bool) (Value, bool) {
	val, ok := lookup(name, localEnv, shouldPanic, Val)
	return val.(Value), ok
}

func lookup(name string, localEnv interface{}, shouldPanic bool, lookupCtx LookupContext) (interface{}, bool) {
	switch lookupCtx {
	default:
		panic("unknown lookup ctx was passed")
	case Val:
		var (
			val Value
			ok  bool
		)

		// search for the backing in local environment first
		val, ok = localEnv.(ValueEnv)[name]
		if !ok {
			// in case of failure, try seeking backing in the global one
			val, ok = valueEnv[name]
			if !ok {
				if shouldPanic {
					panic(fmt.Errorf("no value associated with name %s was found", name))
				}
				return NullValue(), false
			}
		}

		return val, true

	case Type:
		var (
			typeInfo TypeInfo
			ok       bool
		)

		typeInfo, ok = localEnv.(TypeEnv)[name]
		if !ok {
			typeInfo, ok = typeEnv[name]
			if !ok {
				if shouldPanic {
					panic(fmt.Errorf("no type info associated with name %s was found", name))
				}
			}
			return TypeInfo{}, false
		}
		return typeInfo, true
	}
}
