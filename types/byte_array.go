package types

import "unsafe"

type byteArray struct {
	cap  uintptr        // 8bits
	len  uintptr        // 8bits
	data unsafe.Pointer // 8bits
}

// ByteArray api
type ByteArray byteArray

func (b *ByteArray) Cap() int {
	return int(b.cap)
}

func (b *ByteArray) Len() int {
	return int(b.len)
}

func (b *ByteArray) growSizeIfNecessary(size uintptr) {
	if b.len+size < b.cap {
		return
	}
}

func (b *ByteArray) locEnd() uintptr {
	return uintptr(b.data) + b.len
}

func (b *ByteArray) append(bytes []byte) {
	b.growSizeIfNecessary(uintptr(len(bytes)))

	for i := 0; i < len(bytes); i++ {

	}
}

type StructToBytes struct {
	data int64
	str  string
}

type SliceMock struct {
	addr uintptr
	len  int
	cap  int
}

type StructToByte struct {
	data int64
	str  any
}
