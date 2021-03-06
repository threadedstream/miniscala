package vm

import (
	"github.com/ThreadedStream/miniscala/backing"
	"runtime/debug"
)

var (
	chunkStore = make(map[string]Chunk)
)

type Chunk struct {
	funcName    string
	instrStream []Instruction
	argNames    []string // temporary solution
	localVars   map[string]backing.Value
	argPool     map[string]backing.Value
	doesReturn  bool
}

func newChunk(code []Instruction, name string) Chunk {
	chunk := Chunk{}
	chunk.instrStream = code
	//chunk.argPool = make(map[string]backing.Value)
	chunk.argNames = make([]string, 0)
	chunk.funcName = name
	return chunk
}

func lookupChunk(name string, shouldPanic bool, abort func(format string, args ...interface{})) Chunk {
	chunk, ok := chunkStore[name]
	if !ok {
		if shouldPanic {
			if abort != nil {
				abort("no chunk associated with name %s", name)
			} else {
				debug.PrintStack()
				panic("no abort function was passed, so i invoked myself")
			}
		}
		return Chunk{}
	}
	return chunk
}
