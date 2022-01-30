package backing

import "unsafe"

const (
	Size = 131
)

var (
	Venv      SymbolTable
	Tenv      SymbolTable
	hashtable [Size]*Symbol
	markSym   = Symbol{"<mark>", nil}
)

type Symbol struct {
	Name string
	Next *Symbol
}

type SymbolTable *TabTable

func MakeSymbol(name string, next *Symbol) *Symbol {
	symbol := new(Symbol)
	symbol.Name = name
	symbol.Next = next
	return symbol
}

func hash(s string) uint {
	var h uint
	for _, c := range s {
		h = h*65599 + uint(c)
	}
	return h
}

func SSymbol(name string) *Symbol {
	index := hash(name) % Size
	syms := hashtable[index]
	var sym *Symbol
	for sym = syms; sym != nil; sym = sym.Next {
		if name == sym.Name {
			return sym
		}
	}
	sym = MakeSymbol(name, syms)
	hashtable[index] = sym
	return sym
}

func SEmpty() SymbolTable {
	return TabEmpty()
}

func SEnter(table SymbolTable, sym *Symbol, value interface{}) {
	TabEnter(table, unsafe.Pointer(sym), value)
}

func SLook(table SymbolTable, sym *Symbol) interface{} {
	return TabLook(table, unsafe.Pointer(sym))
}

func SBeginScope(table SymbolTable) {
	SEnter(table, &markSym, nil)
}

func SEndScope(table SymbolTable) {
	var sym unsafe.Pointer

	for sym != unsafe.Pointer(&markSym) {
		sym = TabPop(table).(unsafe.Pointer)
	}
}
