package memory

import (
	"testing"
	"unsafe"
)

func TestCMemAllocator_Allocate(t *testing.T) {
	addr, err := platformCInstance.allocate(4 * 1024 * 1024 * 1024)
	if err != nil {
		t.Logf("allocate err:%v", err)
	}

	var intValue = (*int)(addr)
	*(*int)(addr) = 10

	t.Logf("addr:%v, convert *int point: %v -> %d, set value: %d", addr, intValue, intValue, *intValue)
	platformCInstance.free(addr)
}

func TestGoMemAllocator_Allocate(t *testing.T) {
	needBytes := uint((10 + 7) & ^(7))
	addr, err := platformCInstance.allocate(needBytes)
	defer platformInstance.free(addr, uintptr(needBytes))

	if err != nil {
		t.Logf("allocate err:%v", err)
		return
	}

	//var pageCur uintptr = uintptr(addr)
	//var pageEnd uintptr = uintptr(addr) + uintptr(needBytes)
	var intValue = (*[]byte)(addr)

	// 赋值
	t.Logf("%v", unsafe.Pointer(intValue))

}

func TestCMemAllocator_OutOfMemory(t *testing.T) {
	needNumBytes := uint64(80 * 1024 * 1024 * 1024) // 80GB
	allocate, err := platformCInstance.allocate(uint(uintptr(needNumBytes)))
	if err != nil {
		t.Log("allocate fail: Out Of Memory")
		return
	}
	t.Logf("allocate address: %v", allocate)
	platformCInstance.free(allocate)
}

func TestGoMemAllocator_OutOfMemory(t *testing.T) {
	needNumBytes := uint64(80 * 1024 * 1024 * 1024) // 80GB
	allocate, err := platformInstance.allocate(uintptr(needNumBytes))
	if err != nil {
		t.Log("allocate fail: Out Of Memory")
		return
	}
	t.Logf("allocate address: %v", allocate)
	platformInstance.free(allocate, uintptr(needNumBytes))
}

func TestIntConvertToByte(t *testing.T) {
	testValue := 10
	byteValue := byte(testValue)

	t.Logf("byte value: %v", byteValue)
}
