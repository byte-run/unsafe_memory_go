package memory

import "reflect"

const (
	PageWordSize   = 8
	LargePageShift = 17
	LargePageSize  = 1 << LargePageShift // 125K
	PageSizeShift  = 1 << PageWordSize   // 256byte
)

type MemLocation struct {
	Obj    reflect.Kind
	Offset uintptr // address
}

func (loc *MemLocation) ClearObjAndOffset() {
	loc.Obj = 0
	loc.Offset = uintptr(0)
}

type MemBlock struct {
	MemLocation
	length uintptr
}
