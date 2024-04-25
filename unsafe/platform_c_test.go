package unsafe

import (
	"testing"
)

func TestCMemAllocator_Allocate(t *testing.T) {
	addr := platformCInstance.allocate(4)

	var intValue = (*int)(addr)
	*(*int)(addr) = 10

	t.Logf("addr:%v, convert *int point: %v -> %d, set value: %d", addr, intValue, intValue, *intValue)
	platformCInstance.free(addr)
}

func TestCMemAllocator_AllocateBytes(t *testing.T) {
	numBytes := 4 * 1024 * 1024 * 1024
	addr := platformCInstance.allocate(numBytes)

	var valuePointer = (*[]byte)(addr)

	//for i := 0; i < numBytes; i++ {
	//	*(*byte)(unsafe.Pointer(uintptr(addr) + uintptr(i))) = byte(i)
	//}

	t.Logf("addr: %v, valuePointer: %v", addr, valuePointer)
	t.Logf("value: %v", *valuePointer)

	platformCInstance.free(addr)
}
