package interpreter

import (
	"fmt"
)

var (
	environment = make(map[string]Value)
)

func store(name string, value Value) {
	environment[name] = value
}

func lookup(name string) Value {
	val, ok := environment[name]
	if !ok {
		panic(fmt.Errorf("no entry with name %s was found in an environment", name))
	}
	return val
}

func state() {
	for k, v := range environment {
		fmt.Printf("%s: %v\n", k, v)
	}
}
