package backing

import (
	"fmt"
)

type (
	ValueEnvironment map[string]Value
	TypeEnvironment  map[string]Type
)

type Type struct {
	// TODO(threadedstream)
}

type StoringContext int

const (
	Assign StoringContext = iota
	Declare
)

var (
	valueEnvironment = make(ValueEnvironment)
)

func Store(name string, value Value, localEnv ValueEnvironment, ctx StoringContext) {
	if ctx == Declare {
		_, ok := lookup(name, localEnv, false)
		if ok {
			panic(fmt.Errorf("name %s has already got entry in an environment", name))
		}
	} else if ctx == Assign {
		_, ok := lookup(name, localEnv, false)
		if !ok {
			panic(fmt.Errorf("undefined reference to name %s", name))
		}
		if value.Immutable {
			panic(fmt.Errorf("attempt to assign to immutable memory cell"))
		}
	}

	if localEnv != nil {
		localEnv[name] = value
	} else {
		valueEnvironment[name] = value
	}
}

func Lookup(name string, localEnv ValueEnvironment, shouldPanic bool) (Value, bool) {
	return lookup(name, localEnv, shouldPanic)
}

func lookup(name string, localEnv ValueEnvironment, shouldPanic bool) (Value, bool) {
	var (
		val Value
		ok  bool
	)

	// search for the backing in local environment first
	val, ok = localEnv[name]
	if !ok {
		// in case of failure, try seeking backing in the global one
		val, ok = valueEnvironment[name]
		if !ok {
			if shouldPanic {
				panic(fmt.Errorf("no entry associated with name %s was found", name))
			}
			return NullValue(), false
		}
	}

	return val, true
}
