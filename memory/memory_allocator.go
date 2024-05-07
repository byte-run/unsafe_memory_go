package memory

import "unsafe"

// MemAllocator 内存分配器 --- 同platform交互

type MemAllocator interface {
	Allocate(numBytes uintptr) (unsafe.Pointer, error)
	Free(addr unsafe.Pointer, length uintptr)

	AllocateBlock(numBytes uintptr) (*MemBlock, error)
	FreeBlock(*MemBlock)
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
	memBlock.obj = 0
	memBlock.offset = uintptr(address)
	return memBlock, err
}

func (allocator *GoMemAllocator) FreeBlock(page *MemBlock) {
	platformInstance.free(unsafe.Pointer(page.offset), page.length)
	page.offset = 0
}

type CMemAllocator struct{}

func (allocator *CMemAllocator) Allocate(numBytes uintptr) (unsafe.Pointer, error) {
	address, err := platformCInstance.allocate(numBytes)
	return address, err
}

func (allocator *CMemAllocator) Free(addr unsafe.Pointer, length uintptr) {
	platformCInstance.free(addr)
}

func (allocator *CMemAllocator) AllocateBlock(numBytes uintptr) (*MemBlock, error) {
	address, err := platformInstance.allocate(numBytes)
	if err != nil {
		return nil, err
	}
	memBlock := new(MemBlock)
	memBlock.length = numBytes
	memBlock.obj = 0
	memBlock.offset = uintptr(address)
	return memBlock, err
}

func (allocator *CMemAllocator) FreeBlock(page *MemBlock) {
	platformInstance.free(unsafe.Pointer(page.offset), page.length)
	page.offset = 0
}

var (
	UnsafeGo = &GoMemAllocator{}
	UnsafeC  = &CMemAllocator{}
)
