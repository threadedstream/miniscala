package backing

import (
	"github.com/ThreadedStream/miniscala/assert"
	"unsafe"
)

const (
	TabSize = 131
)

type Binder struct {
	Key     unsafe.Pointer
	Value   interface{}
	Next    *Binder
	PrevTop unsafe.Pointer
}

type TabTable struct {
	Table [TabSize]*Binder
	Top   unsafe.Pointer
}

func NewBinder(key unsafe.Pointer, value interface{}, next *Binder, prevTop unsafe.Pointer) *Binder {
	binder := new(Binder)
	binder.Key = key
	binder.Value = value
	binder.Next = next
	binder.PrevTop = prevTop
	return binder
}

func TabEmpty() *TabTable {
	table := new(TabTable)
	return table
}

func TabEnter(table *TabTable, key unsafe.Pointer, value interface{}) {
	index := uintptr(key) % TabSize
	table.Table[index] = NewBinder(key, value, table.Table[index], table.Top)
	table.Top = key
}

func TabLook(table *TabTable, key unsafe.Pointer) interface{} {
	index := uintptr(key) % TabSize
	for b := table.Table[index]; b != nil; b = b.Next {
		if b.Key == key {
			return b.Value
		}
	}
	return nil
}

func TabPop(table *TabTable) interface{} {
	k := table.Top
	index := uintptr(k) % TabSize
	binder := table.Table[index]
	assert.Assert(binder != nil, "[in TabPop()] expected non-nil value of binder")
	table.Table[index] = binder.Next
	table.Top = binder.PrevTop
	return binder.Key
}

func TabDump(table *TabTable, show func(interface{}, interface{})) {
	key := table.Top
	index := uintptr(key) % TabSize
	binder := table.Table[index]
	if binder == nil {
		return
	}
	table.Top = binder.PrevTop
	show(binder.Key, binder.Value)
	TabDump(table, show)
	assert.Assert(table.Top == binder.PrevTop && table.Table[index] == binder.Next, "[in TabDump()] table.Top == binder.PrevTop && table.Table[index] == binder.Next")
	table.Top = key
	table.Table[index] = binder
}
