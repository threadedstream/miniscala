package vm

import (
	"fmt"
	"github.com/ThreadedStream/miniscala/backing"
	"github.com/ThreadedStream/miniscala/syntax"
	"os"
)

type Stack [256]backing.Value

type ChainEntry struct {
	chunk Chunk
	ip    int
}

type VM struct {
	chunk        Chunk
	ip           int
	stack        Stack
	stackPtr     int
	nestingLevel int
	callChain    [256]ChainEntry
}

func isReservedCall(name string) bool {
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

func NewVM(program *syntax.Program) *VM {
	comp := newCompiler()
	comp.compile(program)
	vm := new(VM)
	vm.chunk = lookupChunk("main", true, vm.abort)
	vm.ip = 0
	vm.nestingLevel = 0
	return vm
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
			// TODO(threadedstream): replace with a call to backing.Lookup
			vm.push(vm.chunk.localValues[loadArg.RefName])
		case *InstrGreaterThan:
			var boolValue backing.Value
			secondOperand := vm.pop()
			firstOperand := vm.pop()
			if firstOperand.IsString() && secondOperand.IsString() {
				boolValue = backing.Value{
					Value:     firstOperand.AsString() > secondOperand.AsString(),
					ValueType: backing.Bool,
				}
			} else if firstOperand.IsFloat() && secondOperand.IsFloat() {
				boolValue = backing.Value{
					Value:     firstOperand.AsFloat() > secondOperand.AsFloat(),
					ValueType: backing.Bool,
				}
			}
			vm.push(boolValue)
		case *InstrLessThan:
			var boolValue backing.Value
			secondOperand := vm.pop()
			firstOperand := vm.pop()
			if firstOperand.IsString() && secondOperand.IsString() {
				boolValue = backing.Value{
					Value:     firstOperand.AsString() < secondOperand.AsString(),
					ValueType: backing.Bool,
				}
			} else if firstOperand.IsFloat() && secondOperand.IsFloat() {
				boolValue = backing.Value{
					Value:     firstOperand.AsFloat() < secondOperand.AsFloat(),
					ValueType: backing.Bool,
				}
			}
			vm.push(boolValue)
		case *InstrGreaterThanOrEqual:
			var boolValue backing.Value
			secondOperand := vm.pop()
			firstOperand := vm.pop()
			if firstOperand.IsString() && secondOperand.IsString() {
				boolValue = backing.Value{
					Value:     firstOperand.AsString() >= secondOperand.AsString(),
					ValueType: backing.Bool,
				}
			} else if firstOperand.IsFloat() && secondOperand.IsFloat() {
				boolValue = backing.Value{
					Value:     firstOperand.AsFloat() >= secondOperand.AsFloat(),
					ValueType: backing.Bool,
				}
			}
			vm.push(boolValue)
		case *InstrLessThanOrEqual:
			var boolValue backing.Value
			secondOperand := vm.pop()
			firstOperand := vm.pop()
			if firstOperand.IsString() && secondOperand.IsString() {
				boolValue = backing.Value{
					Value:     firstOperand.AsString() <= secondOperand.AsString(),
					ValueType: backing.Bool,
				}
			} else if firstOperand.IsFloat() && secondOperand.IsFloat() {
				boolValue = backing.Value{
					Value:     firstOperand.AsFloat() <= secondOperand.AsFloat(),
					ValueType: backing.Bool,
				}
			}
			vm.push(boolValue)
		case *InstrEqual:
			var boolValue backing.Value
			secondOperand := vm.pop()
			firstOperand := vm.pop()

			if firstOperand.IsString() && secondOperand.IsString() {
				boolValue = backing.Value{
					Value:     firstOperand.AsString() == secondOperand.AsString(),
					ValueType: backing.Bool,
				}
			} else if firstOperand.IsFloat() && secondOperand.IsFloat() {
				boolValue = backing.Value{
					Value:     firstOperand.AsFloat() == secondOperand.AsFloat(),
					ValueType: backing.Bool,
				}
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
			vm.push(backing.NullValue())
		case *InstrJmp:
			jmp := vm.chunk.instrStream[oldIp].(*InstrJmp)
			vm.ip += jmp.Offset
		case *InstrJmpIfFalse:
			jmpIfFalse := vm.chunk.instrStream[oldIp].(*InstrJmpIfFalse)
			operand := vm.pop()
			if !operand.AsBool() {
				vm.ip = vm.ip + jmpIfFalse.Offset
			}
		case *InstrCall:
			call := vm.chunk.instrStream[oldIp].(*InstrCall)
			if isReservedCall(call.FuncName) {
				vm.dispatchReservedCall(call.FuncName)
				break
			}
			chunk := lookupChunk(call.FuncName, true, vm.abort)
			vm.callChain[vm.nestingLevel] = ChainEntry{
				chunk: vm.chunk,
				ip:    vm.ip,
			}
			vm.chunk = chunk
			vm.nestingLevel++
			vm.ip = 0
			for i := 0; i < len(call.ArgNames); i++ {
				value := vm.pop()
				backing.StoreValue(call.ArgNames[i], value, vm.chunk.localValues, nil, backing.Assign)
			}
		case *InstrReturn:
			var returnValue backing.Value
			if vm.nestingLevel <= 0 {
				vm.normalexit()
			}
			vm.callChain[vm.nestingLevel] = ChainEntry{}
			vm.nestingLevel--
			if vm.chunk.doesReturn {
				returnValue = vm.pop()
				vm.push(returnValue)
			}
			vm.chunk = vm.callChain[vm.nestingLevel].chunk
			vm.ip = vm.callChain[vm.nestingLevel].ip
		case *InstrSetLocal:
			setLocalInstr := vm.chunk.instrStream[oldIp].(*InstrSetLocal)
			valueToAssign := vm.pop()
			valueType := backing.MiniscalaTypeToValueType(setLocalInstr.Type)
			backing.StoreType(setLocalInstr.Name, valueType, setLocalInstr.Immutable, vm.chunk.localValueTypes)
			backing.StoreValue(setLocalInstr.Name, valueToAssign, vm.chunk.localValues, vm.chunk.localValueTypes, setLocalInstr.StoringCtx)
		}
	}
}
