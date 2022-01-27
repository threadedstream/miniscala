package backing

import "github.com/ThreadedStream/miniscala/assert"

const (
	TabSize = 131
)

type Binder struct {
	Key     interface{}
	Value   interface{}
	Next    *Binder
	PrevTop interface{}
}

type TabTable struct {
	Table [TabSize]*Binder
	Top   interface{}
}

func NewBinder(key interface{}, value interface{}, next *Binder, prevTop interface{}) *Binder {
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

func TabEnter(table *TabTable, key interface{}, value interface{}) {
	index := key.(uintptr) % TabSize
	table.Table[index] = NewBinder(key, value, table.Table[index], table.Top)
	table.Top = key
}

func TabLook(table *TabTable, key interface{}) interface{} {
	index := key.(uintptr) % TabSize
	for b := table.Table[index]; b != nil; b = b.Next {
		if b.Key == key {
			return b.Value
		}
	}
	return nil
}

func TabPop(table *TabTable) interface{} {
	k := table.Top
	index := k.(uintptr) % TabSize
	binder := table.Table[index]
	assert.Assert(binder != nil, "[in TabPop()] expected non-nil value of binder")
	table.Table[index] = binder.Next
	table.Top = binder.PrevTop
	return binder.Key
}

func TabDump(table *TabTable, show func(interface{}, interface{})) {
	key := table.Top
	index := key.(uintptr) % TabSize
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
