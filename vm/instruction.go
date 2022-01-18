package vm

type Instruction int8

const (
	InstrAdd Instruction = iota
	InstrSub
	InstrMul
	InstrDiv
	InstrConst
)
