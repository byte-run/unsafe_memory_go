package types

import (
	"reflect"
	"unsafe"
)

type GoSlice[T uint | int | int32 | int64 | string] struct {
	objType reflect.Kind
	len     int
	cap     int
	data    unsafe.Pointer
}

//func (gs GoSlice) allocate() uintptr {
//
//}
//
//func NewGoSlice(cap int) GoSlice {
//
//	goSlice := GoSlice{
//		cap: cap,
//		len: 0,
//	}
//
//	goSlice.data = unsafe.Pointer(goSlice.allocate())
//
//}
