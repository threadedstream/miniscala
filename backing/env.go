package backing

// This is an experimental environment implementation which is borrowed
// from Andrew Appel's book entitled "Modern compiler implementation in C"
//and adapted to the current project environment

type EntryKind int

const (
	EntryVar EntryKind = iota
	EntryFun
)

var (
	outermostLevel *Level
)

type Level struct {
	parent *Level
}

type EnvEntry struct {
	Kind       EntryKind
	Label      string
	Level      *Level
	ParamTypes []ValueType
	ResultType ValueType
	Immutable  bool
}

func OutermostLevel() *Level {
	if outermostLevel == nil {
		outermostLevel = new(Level)
		outermostLevel.parent = nil
	}
	return outermostLevel
}

func MakeVarEntry(label string, level *Level, resultType ValueType) *EnvEntry {
	entry := new(EnvEntry)
	entry.Kind = EntryVar
	entry.Level = level
	entry.Label = label
	entry.ResultType = resultType
	return entry
}

func MakeFunEntry(label string, paramTypes []ValueType, level *Level, resultType ValueType) *EnvEntry {
	entry := new(EnvEntry)
	entry.Kind = EntryFun
	entry.Label = label
	entry.Level = level
	entry.ParamTypes = paramTypes
	entry.ResultType = resultType
	return entry
}

func BaseTypeEnv() SymbolTable {
	var symTable = SEmpty()

	SEnter(symTable, SSymbol("Int"), Int)
	SEnter(symTable, SSymbol("Function"), Function)
	SEnter(symTable, SSymbol("Unit"), Unit)
	SEnter(symTable, SSymbol("Float"), Float)
	SEnter(symTable, SSymbol("Bool"), Bool)
	SEnter(symTable, SSymbol("Null"), Null)
	SEnter(symTable, SSymbol("Any"), Any)
	SEnter(symTable, SSymbol("Undefined"), Undefined)

	return symTable
}

func BaseValueEnv() SymbolTable {
	var symTable = SEmpty()

	SEnter(symTable, SSymbol("print"), MakeFunEntry(
		"print",
		[]ValueType{String},
		OutermostLevel(),
		Unit))

	return symTable
}
