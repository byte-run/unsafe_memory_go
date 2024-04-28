package memory

import "unsafe"

// MemAllocator 内存分配器 --- 同platform交互

type MemAllocator interface {
	Allocate(numBytes uint64) (unsafe.Pointer, error)
	Free(addr unsafe.Pointer)
}

type GoMemAllocator struct{}

func (allocator *GoMemAllocator) Allocate(numBytes uint64) (unsafe.Pointer, error) {
	address, err := platformInstance.allocate(uintptr(numBytes))
	return address, err
}

func (allocator *GoMemAllocator) Free(ptr unsafe.Pointer) {
	platformInstance.free(ptr, 0)
}

type CMemAllocator struct{}

func (allocator *CMemAllocator) Allocate(numBytes uint64) (unsafe.Pointer, error) {
	address, err := platformCInstance.allocate(uint(uintptr(numBytes)))
	return address, err
}

func (allocator *CMemAllocator) Free(ptr unsafe.Pointer) {
	platformInstance.free(ptr, 0)
}

var (
	UnsafeGo = &GoMemAllocator{}
	UnsafeC  = &CMemAllocator{}
)
