package vm

import (
	"github.com/ThreadedStream/miniscala/backing"
)

type (
	Instruction interface {
		Str() string
	}

	InstrAdd struct {
		instr
	}

	InstrSub struct {
		instr
	}

	InstrMul struct {
		instr
	}

	InstrDiv struct {
		instr
	}

	InstrLoad struct {
		backing.Value
		instr
	}

	InstrGreaterThan struct {
		instr
	}

	InstrGreaterThanOrEqual struct {
		instr
	}

	InstrLessThan struct {
		instr
	}

	InstrLessThanOrEqual struct {
		instr
	}

	InstrEqual struct {
		instr
	}

	InstrTrue struct {
		instr
	}

	InstrFalse struct {
		instr
	}

	InstrNull struct {
		instr
	}
	
    InstrJmp struct {
		Offset int
		instr
	}

	InstrJmpIfFalse struct {
		Offset int
		instr
	}

	instr struct {
		text string
	}
)

func (i instr) Str() string {
	return i.text
}

//const (
//	InstrAdd Instruction = iota
//	InstrSub
//	InstrMul
//	InstrDiv
//	InstrLoad
//	InstrStore
//)
