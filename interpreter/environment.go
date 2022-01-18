package interpreter

import (
	"fmt"
	"github.com/ThreadedStream/miniscala/vm"
)

type Environment map[string]vm.Value

type StoringContext int

const (
	Assign StoringContext = iota
	Declare
)

var (
	environment = make(Environment)
)

// consider the following scenario
// what most awful is that currently language permits it
// def func(x: Int): Unit {
// 		var x = 54
// }
// a possible solution is to add a storing context
func store(name string, value vm.Value, localEnv Environment, ctx StoringContext) {
	if ctx == Declare {
		_, ok := lookup(name, localEnv, false)
		if ok {
			panic(fmt.Errorf("name %s has already got entry in an environment", name))
		}
	}
	if localEnv != nil {
		localEnv[name] = value
	} else {
		environment[name] = value
	}

}

func lookup(name string, localEnv Environment, shouldPanic bool) (vm.Value, bool) {
	var (
		val vm.Value
		ok  bool
	)

	// search for the value in local environment first
	val, ok = localEnv[name]
	if !ok {
		// in case of failure, try seeking value in the global one
		val, ok = environment[name]
		if !ok {
			if shouldPanic {
				panic(fmt.Errorf("no entry associated with name %s was found", name))
			}
			return vm.NullValue(), false
		}
	}

	return val, true
}

func state() {

}
