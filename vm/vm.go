package vm

import (
	"github.com/ThreadedStream/miniscala/backing"
)

type Chunk struct {
	instrStream []Instruction
}

type VM struct {
	chunk    *Chunk
	ip       int
	stack    [256]backing.Value
	stackPtr int8
}

func InitializeVm(code []Instruction) *VM {
	vm := new(VM)
	vm.chunk = newChunk(code)
	vm.ip = 0
	vm.stackPtr = 0
	return vm
}

func newChunk(code []Instruction) *Chunk {
	chunk := new(Chunk)
	chunk.instrStream = code
	return chunk
}

func (vm *VM) resetStack() {
	vm.stackPtr = 0
}

func (vm *VM) push(v backing.Value) {
	vm.stack[vm.stackPtr] = v
	vm.stackPtr++
}

func (vm *VM) pop() backing.Value {
	vm.stackPtr--
	return vm.stack[vm.stackPtr]
}

func (vm *VM) Run() {
	for vm.ip < len(vm.chunk.instrStream) {
		switch vm.chunk.instrStream[vm.ip].(type) {
		case *InstrAdd:
			secondOperand := vm.pop()
			firstOperand := vm.pop()
			vm.push(backing.Add(firstOperand, secondOperand, nil, backing.Vm))
		case *InstrSub:
			secondOperand := vm.pop()
			firstOperand := vm.pop()
			vm.push(backing.Sub(firstOperand, secondOperand, nil, backing.Vm))
		case *InstrMul:
			secondOperand := vm.pop()
			firstOperand := vm.pop()
			vm.push(backing.Mul(firstOperand, secondOperand, nil, backing.Vm))
		case *InstrDiv:
			secondOperand := vm.pop()
			firstOperand := vm.pop()
			vm.push(backing.Div(firstOperand, secondOperand, nil, backing.Vm))
		case *InstrLoad:
			load := vm.chunk.instrStream[vm.ip].(*InstrLoad)
			vm.push(load.Value)
		}
		vm.ip++
	}
}
