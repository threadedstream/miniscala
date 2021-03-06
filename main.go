package main

import "C"
import (
	"github.com/ThreadedStream/miniscala/syntax"
	"github.com/ThreadedStream/miniscala/typecheck"
	"github.com/ThreadedStream/miniscala/vm"
)

func main() {
	path := "sources/sort.miniscala"
	program, hadErrors := syntax.Parse(path)
	if hadErrors {
		return
	}
	hadErrors = typecheck.Typecheck(program)
	if hadErrors {
		return
	}

	vmHandle := vm.NewVM(program)
	vmHandle.Run()
}
