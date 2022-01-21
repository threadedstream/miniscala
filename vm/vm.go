package vm

import (
	"fmt"
	"github.com/ThreadedStream/miniscala/backing"
	"os"
)

type Stack [256]backing.Value

type Chunk struct {
	funcName    string
	instrStream []Instruction
	localValues map[string]backing.Value
}

type ChainEntry struct {
	chunk Chunk
	ip    int
}

type VM struct {
	chunk        Chunk
	ip           int
	chunkStore   map[string]Chunk
	stack        Stack
	stackPtr     int
	nestingLevel int
	callChain    [256]ChainEntry
}

func (vm *VM) isReservedCall(name string) bool {
	switch name {
	default:
		return false
	case "print":
		return true
	}
}

func (vm *VM) dispatchReservedCall(name string) {
	switch name {
	default:
		vm.abort("%s is not a call to reserved function", name)
	case "print":
		vm.executePrint()
	}
}

func (vm *VM) executePrint() {
	valueToPrint := vm.pop()
	fmt.Printf("%v", valueToPrint.Value)
}

func InitializeVm(code ...[]Instruction) *VM {
	vm := new(VM)
	vm.chunkStore = make(map[string]Chunk)
	vm.chunkStore["main"] = newChunk(code[0], "main")
	vm.chunkStore["fac"] = newChunk(code[1], "fac")
	vm.chunkStore["fib"] = newChunk(code[2], "fib")
	vm.chunk = vm.chunkStore["main"]
	vm.ip = 0
	vm.nestingLevel = 0
	return vm
}

func newChunk(code []Instruction, name string) Chunk {
	chunk := Chunk{}
	chunk.instrStream = code
	chunk.localValues = make(map[string]backing.Value)
	chunk.funcName = name
	return chunk
}

func (vm *VM) resetStack() {
	vm.stackPtr = 0
}

func (vm *VM) abort(format string, args ...interface{}) {
	panic(fmt.Errorf(format, args))
}

func (vm *VM) normalexit() {
	os.Exit(1)
}

func (vm *VM) lookupChunk(name string, shouldPanic bool) Chunk {
	chunk, ok := vm.chunkStore[name]
	if !ok {
		if shouldPanic {
			vm.abort("no chunk associated with name %s", name)
		}
		return Chunk{}
	}
	return chunk
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
		case *InstrLoadRef:
			loadArg := vm.chunk.instrStream[oldIp].(*InstrLoadRef)
			vm.push(vm.chunk.localValues[loadArg.RefName])
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
				ValueType: backing.Null,
			}
			vm.push(nullValue)
		case *InstrJmp:
			jmp := vm.chunk.instrStream[oldIp].(*InstrJmp)
			vm.ip = jmp.Offset
		case *InstrJmpIfFalse:
			jmpIfFalse := vm.chunk.instrStream[oldIp].(*InstrJmpIfFalse)
			operand := vm.pop()
			if !operand.AsBool() {
				vm.ip = vm.ip + jmpIfFalse.Offset
			}
		case *InstrCall:
			call := vm.chunk.instrStream[oldIp].(*InstrCall)
			if vm.isReservedCall(call.FuncName) {
				vm.dispatchReservedCall(call.FuncName)
				break
			}
			chunk := vm.lookupChunk(call.FuncName, true)
			vm.callChain[vm.nestingLevel] = ChainEntry{
				chunk: vm.chunk,
				ip:    vm.ip,
			}
			vm.chunk = chunk
			vm.nestingLevel++
			vm.chunk.localValues = make(map[string]backing.Value)
			vm.ip = 0
			for i := 0; i < len(call.ArgNames); i++ {
				value := vm.pop()
				vm.chunk.localValues[call.ArgNames[i]] = value
			}
		case *InstrReturn:
			if vm.nestingLevel <= 0 {
				vm.normalexit()
			}
			vm.callChain[vm.nestingLevel] = ChainEntry{}
			vm.nestingLevel--
			returnValue := vm.pop()
			vm.chunk = vm.callChain[vm.nestingLevel].chunk
			vm.ip = vm.callChain[vm.nestingLevel].ip
			vm.push(returnValue)
		}
	}
}
