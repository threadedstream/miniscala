package vm

type Chunk struct {
}

type VM struct {
	chunk    *Chunk
	ip       *int8
	stack    []Value
	stackTop *Value
}
