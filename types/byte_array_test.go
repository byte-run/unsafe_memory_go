package types

import (
	"testing"
	"unsafe"
)

func TestByteSize(t *testing.T) {
	emptyByteArray := ByteArray{}
	emptySize := unsafe.Sizeof(emptyByteArray)
	t.Logf("empty byteArray size: %d", emptySize)

	byteArray := ByteArray{10, 5, nil}
	valueSize := unsafe.Sizeof(byteArray)
	t.Logf("value byteArray size: %d", valueSize)

	//
}

func TestGoType(t *testing.T) {
	pointerVar := unsafe.Pointer(uintptr(int8(1)))
	pointerSize := unsafe.Sizeof(pointerVar)
	t.Logf("pointer size: %d", pointerSize)
}
