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

func TestGoSizeof(t *testing.T) {
	// array or slice
	var testArray = [5]int{1, 2, 3, 4, 5} // int type: 8 bits
	testArraySize := unsafe.Sizeof(testArray)
	t.Logf("testArray size: %d", testArraySize) // testArray size: 40 = 8 * 5

	/*
		SliceHeader size = 24
	*/
	var testSlice = make([]int, 5)
	testSlice = append(testSlice, 1, 2, 3)
	testSliceSize := unsafe.Sizeof(testSlice)
	t.Logf("testSlice size: %d", testSliceSize) // testSlice size: 24 =

	testSlice = append(testSlice, 4, 5, 6)
	testSliceSize = unsafe.Sizeof(testSlice)
	t.Logf("testSlice size: %d", testSliceSize) // testSlice size: 24 =
}
