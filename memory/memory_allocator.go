package memory

import "unsafe"

// MemAllocator 内存分配器 --- 同platform交互

type MemAllocator interface {
	Allocate(numBytes uintptr) (unsafe.Pointer, error)
	Free(addr unsafe.Pointer, length uintptr)
}

type GoMemAllocator struct{}

func (allocator *GoMemAllocator) Allocate(numBytes uintptr) (unsafe.Pointer, error) {
	address, err := platformInstance.allocate(numBytes)
	return address, err
}

func (allocator *GoMemAllocator) Free(addr unsafe.Pointer, length uintptr) {
	platformInstance.free(addr, length)
}

// MemBlock版本
func (allocator *GoMemAllocator) AllocateBlock(numBytes uintptr) (*MemBlock, error) {
	address, err := platformInstance.allocate(numBytes)
	if err != nil {
		return nil, err
	}
	memBlock := new(MemBlock)
	memBlock.length = numBytes
	memBlock.Obj = 0
	memBlock.Offset = uintptr(address)
	return memBlock, err
}

func (allocator *GoMemAllocator) FreeBlock(memBlock *MemBlock) {
	platformInstance.free(unsafe.Pointer(memBlock.Offset), memBlock.length)
	memBlock.Offset = 0
}

type CMemAllocator struct{}

func (allocator *CMemAllocator) Allocate(numBytes uintptr) (unsafe.Pointer, error) {
	address, err := platformCInstance.allocate(numBytes)
	return address, err
}

func (allocator *CMemAllocator) Free(addr unsafe.Pointer, length uintptr) {
	platformCInstance.free(addr)
}

var (
	UnsafeGo = &GoMemAllocator{}
	UnsafeC  = &CMemAllocator{}
)
