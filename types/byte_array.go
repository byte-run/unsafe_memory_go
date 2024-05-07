package types

import (
	"github.com/byte-run/unsafe_mem_go/memory"
	"unsafe"
)

//type byteArray struct {
//	cap uintptr // 8bits
//	len uintptr // 8bits
//	data unsafe.Pointer // 8bits
//}

// ByteArray api
//type ByteArray byteArray
//
//func (b *ByteArray) Cap() int {
//	return int(b.cap)
//}
//
//func (b *ByteArray) Len() int {
//	return int(b.len)
//}
//
//func (b *ByteArray) growSizeIfNecessary(size uintptr) {
//	if b.len+size < b.cap {
//		return
//	}
//}
//
//func (b *ByteArray) locEnd() uintptr {
//	return uintptr(b.data) + b.len
//}
//
//func (b *ByteArray) append(bytes []byte) {
//	b.growSizeIfNecessary(uintptr(len(bytes)))
//
//	for i := 0; i < len(bytes); i++ {
//
//	}
//}

/*
type byte uint8 (in Go)
*/
const byteWordWidth = 8

type byteArray struct {
	memoryBlock *memory.MemBlock
	baseObj     any
	baseOffset  uintptr
	length      uintptr
}

func (array *byteArray) Length() int64 {
	return int64(array.length)
}

func (array *byteArray) SetValue(index uintptr, value any) {

}

func (array *byteArray) GetValue(index uintptr) byte {
	loc := array.baseOffset + index*byteWordWidth
	return *(*byte)(unsafe.Pointer(loc))
}

// GetWordWidth	根据baseObj获取
//func (array *byteArray) GetWordWidth() uintptr {
//
//	switch array.baseObj.(type) {
//	case reflect.Kind:
//
//	}
//}

type ByteArray byteArray

func (array ByteArray) Length() int64 {
	return int64(array.length / byteWordWidth)
}

func (array *ByteArray) SetValue(index uintptr, value byte) {

}

func NewByteArray(memBlock *memory.MemBlock) *ByteArray {
	var array = new(ByteArray)
	array.memoryBlock = memBlock
	array.baseOffset = memBlock.GetOffset()
	array.baseObj = memBlock.GetObj()
	array.length = memBlock.Size() // TODO 要不要去除以8？这里是byteArray形式，应该不用

	return array
}
