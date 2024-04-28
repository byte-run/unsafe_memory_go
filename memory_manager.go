package tcmallocgo

import (
	"github.com/byte-run/unsafe_mem_go/memory"
	"github.com/byte-run/unsafe_mem_go/utils"
	"sync"
	"unsafe"
)

//type MemoryManager struct {
//}
//
//// AcquireMemory 从mem-pool中申请合适的内存块
//func (mem MemoryManager) AcquireMemory(numbytes uintptr) unsafe.Pointer {
//	// TODO Unimplement
//	return unsafe.Pointer(uintptr(0))
//}
//
//// ReleaseMemory 将使用完的内存块释放回mem-pool中
//func (mem MemoryManager) ReleaseMemory(offset unsafe.Pointer, numbytes uintptr) {
//
//}
//
//// 单例控制
//var memoryManager *MemoryManager
//
//func NewMemoryManager() *MemoryManager {
//	// TODO Unimplement
//
//	if memoryManager != nil {
//		return memoryManager
//	}
//
//	manager := new(MemoryManager)
//	memoryManager = manager
//	return manager
//}

type MemoryManager struct {
	unsafe       *unsafeMemory
	conf         *MemoryConfig
	memAllocator memory.MemAllocator

	mu sync.Mutex //
}

// 所有的操作方法都需要检查unsafe
func (mem MemoryManager) checkUnsafeIsNil() bool {
	return mem.unsafe == nil
}

func (mem MemoryManager) AcquireStorageMemory(numBytes uint64) (bool, error) {
	if numBytes < 0 {
		return false, utils.AcquireMemoryBytesZeroError
	}

	// 不加锁
	return mem.unsafe.AcquireStorageMemory(uintptr(numBytes))
}

func (mem MemoryManager) ReleaseStorageMemory() {
	mem.mu.Lock()
	defer mem.mu.Unlock()

	mem.unsafe.ReleaseStorageMemory()
}

func (mem MemoryManager) AcquireComputeMemory(numBytes uint64) (uintptr, error) {
	if numBytes < 0 {
		return emptyValue, utils.AcquireMemoryBytesZeroError
	}

	mem.mu.Lock()
	// 从memory pool中获取可用的memory size
	//mem.unsafe.

	return 0, nil
}

func (mem MemoryManager) ReleaseComputeMemory(numBytes uint64) {

}

func (mem MemoryManager) AllocateComputePage(numBytes uint64) uintptr {
	// 当前不加page size limit

	return 0

}

func (mem MemoryManager) AllocateStoragePage(numBytes uint64) (uintptr, error) {
	if numBytes < 0 {
		return emptyValue, utils.AcquireMemoryBytesZeroError
	}

	mem.mu.Lock()
	defer mem.mu.Unlock()

	mem.memAllocator.Allocate(numBytes)

	return 0, nil
}

func (mem MemoryManager) FreeStoragePage(addr uintptr) {
	mem.mu.Lock()
	defer mem.mu.Unlock()

	mem.memAllocator.Free(unsafe.Pointer(addr))
	// 再由unsafe -> pool 释放
	mem.unsafe.ReleaseStorageMemory()
}

func InitMemoryManager(config *MemoryConfig) *MemoryManager {
	manager := new(MemoryManager)
	manager.conf = config
	manager.unsafe = newUnsafeMemory(config)
	manager.memAllocator = dynamicMemAllocator("C")

	return manager
}

// allocator 根据设定采用不同的内存分配实现
func dynamicMemAllocator(allocMode string) memory.MemAllocator {
	if allocMode == "C" {
		return memory.UnsafeC
	}
	return memory.UnsafeGo

}
