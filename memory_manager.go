package tcmallocgo

import (
	"github.com/byte-run/unsafe_mem_go/memory"
	"github.com/byte-run/unsafe_mem_go/utils"
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
	staticPool   *staticMemoryManage
	conf         *MemoryConfig
	memAllocator memory.MemAllocator

	pageTable map[uintptr]uintptr
}

// 所有的操作方法都需要检查unsafe
func (mem *MemoryManager) checkUnsafeIsNil() bool {
	return mem.staticPool == nil
}

func (mem MemoryManager) AcquireStorageMemory(numBytes uintptr) (bool, utils.MemPoolWarn, error) {
	if numBytes < 0 {
		return false, nil, utils.AcquireMemoryBytesZeroError
	}

	return mem.staticPool.AcquireStorageMemory(uintptr(numBytes))
}

func (mem MemoryManager) ReleaseStorageMemory(numBytes uintptr) error {
	if numBytes < 0 {
		return utils.AcquireMemoryBytesZeroError
	}

	mem.staticPool.ReleaseStorageMemory(numBytes)
	return nil
}

func (mem MemoryManager) ReleaseAllStorageMemory() {
	mem.staticPool.ReleaseAllStorageMemory()
}

func (mem *MemoryManager) AcquireShuffleMemory(numBytes uintptr) (uintptr, utils.MemPoolWarn, error) {
	if numBytes < 0 {
		return emptyValue, nil, utils.AcquireMemoryBytesZeroError
	}

	// 从memory pool中获取可用的memory size
	//mem.unsafe.
	return mem.staticPool.acquireShuffleMemory(numBytes)
}

func (mem *MemoryManager) ReleaseShuffleMemory(numBytes uintptr) error {
	if numBytes < 0 {
		return utils.AcquireMemoryBytesZeroError
	}

	mem.staticPool.ReleaseShuffleMemory(numBytes)
	return nil
}

func (mem *MemoryManager) AcquireIntersectionMemory(numBytes uintptr) (uintptr, utils.MemPoolWarn, error) {
	if numBytes < 0 {
		return emptyValue, nil, utils.AcquireMemoryBytesZeroError
	}

	return mem.staticPool.acquireIntersectionMemory(numBytes)
}

func (mem *MemoryManager) ReleaseIntersectionMemory(numBytes uintptr) error {
	if numBytes < 0 {
		return utils.AcquireMemoryBytesZeroError
	}
	mem.staticPool.ReleaseIntersectionMemory(numBytes)
	return nil
}

func (mem *MemoryManager) AllocatePage(numBytes uintptr) (uintptr, error) {
	// 当前不加page size limit
	if numBytes < 0 {
		return emptyValue, utils.AcquireMemoryBytesZeroError
	}

	addr, err := mem.memAllocator.Allocate(uintptr(numBytes))
	if err != nil {
		return 0, err
	}
	return uintptr(addr), nil
}

func (mem *MemoryManager) FreePage(addr uintptr, numBytes uintptr) {
	mem.memAllocator.Free(unsafe.Pointer(addr), numBytes)
}

//func (mem *MemoryManager) AllocateStoragePage(numBytes uintptr) (uintptr, error) {
//	if numBytes < 0 {
//		return emptyValue, utils.AcquireMemoryBytesZeroError
//	}
//
//	mem.mu.Lock()
//	defer mem.mu.Unlock()
//
//	addr, err := mem.memAllocator.Allocate(numBytes)
//	if err != nil {
//		return 0, err
//	}
//	return uintptr(addr), nil
//}
//
//func (mem *MemoryManager) FreeStoragePage(addr uintptr, numBytes uintptr) {
//	mem.mu.Lock()
//	defer mem.mu.Unlock()
//
//	mem.memAllocator.Free(unsafe.Pointer(addr), 0)
//	// 再由unsafe -> pool 释放
//	mem.staticPool.ReleaseStorageMemory(numBytes)
//}

// Destory 释放所有分配的内存
func (mem *MemoryManager) CleanAllAllocatedMemory() {
	for size, addrValue := range mem.pageTable {
		mem.memAllocator.Free(unsafe.Pointer(addrValue), size)
	}
	mem.staticPool.ResetPoolUsed()
}

func InitMemoryManager(config *MemoryConfig) *MemoryManager {
	manager := new(MemoryManager)
	manager.conf = config
	manager.staticPool = newStaticMemoryManage(config)
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

//var MemoryManagerInstance *MemoryManager = InitMemoryManager(&MemoryConfig{StorageMem: "5G", ShuffleMem: "5G", IntersectionMem: "5G"})
