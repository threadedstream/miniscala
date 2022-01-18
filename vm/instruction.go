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
