package assert

import "fmt"

func Assert(cond bool, format string, args ...interface{}) {
	if !cond {
		panic(fmt.Errorf(format, args))
	}
}

func AssertCallback(cond bool, callback func(string, ...interface{}), format string, args ...interface{}) {
	if !cond {
		callback(format, args)
	}
}
