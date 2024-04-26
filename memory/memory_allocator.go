package memory

import "unsafe"

// MemAllocator 内存分配器 --- 同platform交互

type MemAllocator interface {
	Allocate(numBytes int) unsafe.Pointer
	Free(addr unsafe.Pointer)
}

type GoMemAllocator struct{}

func (allocator *GoMemAllocator) Allocate(numBytes int) unsafe.Pointer {
	address := platformInstance.allocate(uintptr(numBytes))
	return address
}

func (allocator *GoMemAllocator) Free(ptr unsafe.Pointer) {
	platformInstance.free(ptr, 0)
}

type CMemAllocator struct{}

func (allocator *CMemAllocator) Allocate(numBytes int) unsafe.Pointer {
	address := platformCInstance.allocate(uintptr(numBytes))
	return address
}

func (allocator *CMemAllocator) Free(ptr unsafe.Pointer) {
	platformInstance.free(ptr, 0)
}

var (
	UnsafeGo = &GoMemAllocator{}
	UnsafeC  = &CMemAllocator{}
)
