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

	InstrMod struct {
		instr
	}

	InstrLogicalAnd struct {
		instr
	}

	InstrLogicalOr struct {
		instr
	}

	InstrLogicalNot struct {
		instr
	}

	InstrLoadImm struct {
		backing.Value
		instr
	}

	InstrLoadRef struct {
		RefName string
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

	InstrSetLocal struct {
		Name       string
		Type       string
		StoringCtx backing.StoringContext
		Immutable  bool
		instr
	}

	InstrCall struct {
		FuncName string
		ArgNames []string
		instr
	}

	InstrReturn struct {
		instr
	}

	instr struct {
		text string
	}
)

func (i instr) Str() string {
	return i.text
}
