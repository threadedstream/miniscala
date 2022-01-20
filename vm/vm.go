package vm

import (
	"fmt"
	"github.com/ThreadedStream/miniscala/backing"
	"os"
)

type Stack [256]backing.Value

type Chunk struct {
	instrStream []Instruction
	localValues map[int]backing.Value
	stackFrame  Stack
	stackPtr    int
}

type VM struct {
	chunk        *Chunk
	ip           int
	chunkStore   map[string]*Chunk
	nestingLevel int
	callChain    [256]*Chunk
}

func InitializeVm(code []Instruction) *VM {
	vm := new(VM)
	vm.chunk = newChunk(code)
	vm.ip = 0
	vm.chunkStore = make(map[string]*Chunk)
	vm.chunkStore["fac"] = vm.chunk
	return vm
}

func newChunk(code []Instruction) *Chunk {
	chunk := new(Chunk)
	chunk.instrStream = code
	chunk.stackPtr = 0
	chunk.localValues = make(map[int]backing.Value)
	// TODO(threadedstream): erase it after testing
	chunk.localValues[0] = backing.Value{
		Value:     5.0,
		ValueType: backing.Float,
	}
	return chunk
}

func (vm *VM) resetStack() {
	vm.chunk.stackPtr = 0
}

func (vm *VM) abort(format string, args ...interface{}) {
	panic(fmt.Errorf(format, args))
}

func (vm *VM) normalexit() {
	os.Exit(1)
}

func (vm *VM) lookupChunk(name string, shouldPanic bool) *Chunk {
	chunk, ok := vm.chunkStore[name]
	if !ok {
		if shouldPanic {
			vm.abort("no chunk associated with name %s", name)
		}
		return nil
	}
	return chunk
}

func (vm *VM) push(v backing.Value) {
	vm.chunk.stackFrame[vm.chunk.stackPtr] = v
	vm.chunk.stackPtr++
}

func (vm *VM) pop() backing.Value {
	vm.chunk.stackPtr--
	return vm.chunk.stackFrame[vm.chunk.stackPtr]
}

func (vm *VM) Run() {
	for vm.ip < len(vm.chunk.instrStream) {
		oldIp := vm.ip
		vm.ip++
		switch vm.chunk.instrStream[oldIp].(type) {
		default:
			vm.abort("unknown instruction %v", vm.chunk.instrStream[vm.ip])
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
		case *InstrLoadImm:
			load := vm.chunk.instrStream[oldIp].(*InstrLoadImm)
			vm.push(load.Value)
		case *InstrLoadArg:
			loadArg := vm.chunk.instrStream[oldIp].(*InstrLoadArg)
			vm.push(vm.chunk.localValues[loadArg.Idx])
		case *InstrGreaterThan:
			secondOperand := vm.pop()
			firstOperand := vm.pop()
			boolValue := backing.Value{
				Value:     firstOperand.AsFloat() > secondOperand.AsFloat(),
				ValueType: backing.Bool,
			}
			vm.push(boolValue)
		case *InstrLessThan:
			secondOperand := vm.pop()
			firstOperand := vm.pop()
			boolValue := backing.Value{
				Value:     firstOperand.AsFloat() < secondOperand.AsFloat(),
				ValueType: backing.Bool,
			}
			vm.push(boolValue)
		case *InstrGreaterThanOrEqual:
			secondOperand := vm.pop()
			firstOperand := vm.pop()
			boolValue := backing.Value{
				Value:     firstOperand.AsFloat() >= secondOperand.AsFloat(),
				ValueType: backing.Bool,
			}
			vm.push(boolValue)
		case *InstrLessThanOrEqual:
			secondOperand := vm.pop()
			firstOperand := vm.pop()
			boolValue := backing.Value{
				Value:     firstOperand.AsFloat() <= secondOperand.AsFloat(),
				ValueType: backing.Bool,
			}
			vm.push(boolValue)
		case *InstrEqual:
			secondOperand := vm.pop()
			firstOperand := vm.pop()
			boolValue := backing.Value{
				Value:     firstOperand.AsFloat() == secondOperand.AsFloat(),
				ValueType: backing.Bool,
			}
			vm.push(boolValue)
		case *InstrTrue:
			boolValue := backing.Value{
				Value:     true,
				ValueType: backing.Bool,
			}
			vm.push(boolValue)
		case *InstrFalse:
			boolValue := backing.Value{
				Value:     false,
				ValueType: backing.Bool,
			}
			vm.push(boolValue)
		case *InstrNull:
			nullValue := backing.Value{
				Value:     nil,
				ValueType: backing.Bool,
			}
			vm.push(nullValue)
		case *InstrJmp:
			jmp := vm.chunk.instrStream[oldIp].(*InstrJmp)
			vm.ip = jmp.Offset
		case *InstrJmpIfFalse:
			jmpIfFalse := vm.chunk.instrStream[oldIp].(*InstrJmpIfFalse)
			operand := vm.pop()
			if operand.AsBool() {
				vm.ip = jmpIfFalse.Offset
			}
		case *InstrCall:
			call := vm.chunk.instrStream[oldIp].(*InstrCall)
			chunk := vm.lookupChunk(call.FuncName, true)
			vm.chunk = chunk
			vm.push(backing.Value{
				Value:     vm.ip,
				ValueType: backing.Int,
			})
			vm.nestingLevel++
			vm.callChain[vm.nestingLevel] = chunk
			vm.ip = 0
		case *InstrReturn:
			if vm.nestingLevel == 0 {
				vm.normalexit()
			}
			vm.nestingLevel--
			returnValue := vm.pop()
			vm.chunk = vm.callChain[vm.nestingLevel]
			vm.ip = vm.pop().AsInt()
			vm.push(returnValue)
		}
	}
}
