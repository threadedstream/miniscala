package assert

import "fmt"

func Assert(cond bool, format string, args ...interface{}) {
	if !cond {
		panic(fmt.Errorf(format, args))
	}
}
