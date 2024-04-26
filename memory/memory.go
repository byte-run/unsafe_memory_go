package memory

import "reflect"

type MemManager struct {
}

//func (mem MemManager) memoryAllocator() *MemAllocator {
//
//}

type MemLocation struct {
	Obj    reflect.Kind
	Offset uintptr // address
}
