package tcmallocgo

import "github.com/byte-run/unsafe_mem_go/memory"

type unsafeMemory struct {
	allocatorMode    string
	storagePool      *storageMemoryPool
	shufflePool      *shuffleMemoryPool
	intersectionPool *intersectionMemoryPool
}

func newUnsafeMemory() *unsafeMemory {
	unsafeMem := new(unsafeMemory)

	// init allocator
	//mem.allocator = memory.MemAllocator{}

	// init memPool

	return unsafeMem
}

type UnsafeManager struct {
	unsafeMemory
}

// 对外接口，单实列
var _unsafe *unsafeMemory = nil

var (
	unsafeGo = memory.GoMemAllocator{}
	unsafeC  = memory.CMemAllocator{}
)

func newUnsafeManager(memLimits ...int) *UnsafeManager {
	mem := new(UnsafeManager)

	return mem
}

// allocator 根据设定采用不同的内存分配实现
func dynamicMemAllocator() *memory.MemAllocator {

	return nil
}
