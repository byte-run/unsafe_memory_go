package tcmallocgo

import (
	"fmt"
	"github.com/byte-run/unsafe_mem_go/bitset"
	"github.com/byte-run/unsafe_mem_go/memory"
	"github.com/byte-run/unsafe_mem_go/utils"
	"sync"
	"unsafe"
)

type TaskMemoryManager struct {
	staticPool   *staticMemoryManage
	memAllocator memory.MemAllocator

	pageTable      []*memory.MemBlock
	allocatedPages bitset.BitSet

	// lock
	lock sync.RWMutex
}

// 所有的操作方法都需要检查unsafe
func (tmm *TaskMemoryManager) checkUnsafeIsNil() bool {
	return tmm.staticPool == nil
}

func (tmm TaskMemoryManager) AcquireStorageMemory(numBytes uintptr) (bool, utils.MemPoolWarn, error) {
	if numBytes < 0 {
		return false, nil, utils.AcquireMemoryBytesZeroError
	}

	return tmm.staticPool.AcquireStorageMemory(uintptr(numBytes))
}

func (tmm TaskMemoryManager) ReleaseStorageMemory(numBytes uintptr) error {
	if numBytes < 0 {
		return utils.AcquireMemoryBytesZeroError
	}

	tmm.staticPool.ReleaseStorageMemory(numBytes)
	return nil
}

func (tmm TaskMemoryManager) ReleaseAllStorageMemory() {
	tmm.staticPool.ReleaseAllStorageMemory()
}

// --------------------- update -------------------------------

// --------------------- update finish ------------------------
func (tmm *TaskMemoryManager) AcquireShuffleMemory(numBytes uintptr) (uintptr, utils.MemPoolWarn, error) {
	if numBytes < 0 {
		return emptyValue, nil, utils.AcquireMemoryBytesZeroError
	}

	// 从memory pool中获取可用的memory size
	//tmm.unsafe.
	return tmm.staticPool.acquireShuffleMemory(numBytes)
}

func (tmm *TaskMemoryManager) ReleaseShuffleMemory(numBytes uintptr) error {
	if numBytes < 0 {
		return utils.AcquireMemoryBytesZeroError
	}

	tmm.staticPool.ReleaseShuffleMemory(numBytes)
	return nil
}

func (tmm *TaskMemoryManager) AcquireIntersectionMemory(numBytes uintptr) (uintptr, utils.MemPoolWarn, error) {
	if numBytes < 0 {
		return emptyValue, nil, utils.AcquireMemoryBytesZeroError
	}

	return tmm.staticPool.acquireIntersectionMemory(numBytes)
}

func (tmm *TaskMemoryManager) ReleaseIntersectionMemory(numBytes uintptr) error {
	if numBytes < 0 {
		return utils.AcquireMemoryBytesZeroError
	}
	tmm.staticPool.ReleaseIntersectionMemory(numBytes)
	return nil
}

func (tmm *TaskMemoryManager) AllocatePage(numBytes uintptr) (uintptr, error) {
	// 当前不加page size limit
	if numBytes < 0 {
		return emptyValue, utils.AcquireMemoryBytesZeroError
	}

	addr, err := tmm.memAllocator.Allocate(uintptr(numBytes))
	if err != nil {
		return 0, err
	}
	return uintptr(addr), nil
}

func (tmm *TaskMemoryManager) FreePage(addr uintptr, numBytes uintptr) {
	tmm.memAllocator.Free(unsafe.Pointer(addr), numBytes)
}

// -------------- insert content
func (tmm *TaskMemoryManager) FreeBlockPage(page *memory.MemBlock) {
	// TODO assert

	tmm.pageTable[page.PageNumber] = nil
	tmm.lock.Lock()
	tmm.allocatedPages.Clear(page.PageNumber)
	tmm.lock.Unlock()
	// TODO waiting to Logger
	fmt.Printf("Free page number %d (%d bytes)", page.PageNumber, page.Size)

	pageSize := page.Size()
	page.PageNumber = memory.FreedInTMMPageNumber
	tmm.memAllocator.FreeBlock(page)
	tmm.ReleaseShuffleMemory(pageSize)

}

//func (tmm *TaskMemoryManager) AllocateStoragePage(numBytes uintptr) (uintptr, error) {
//	if numBytes < 0 {
//		return emptyValue, utils.AcquireMemoryBytesZeroError
//	}
//
//	tmm.mu.Lock()
//	defer tmm.mu.Unlock()
//
//	addr, err := tmm.memAllocator.Allocate(numBytes)
//	if err != nil {
//		return 0, err
//	}
//	return uintptr(addr), nil
//}
//
//func (tmm *TaskMemoryManager) FreeStoragePage(addr uintptr, numBytes uintptr) {
//	tmm.mu.Lock()
//	defer tmm.mu.Unlock()
//
//	tmm.memAllocator.Free(unsafe.Pointer(addr), 0)
//	// 再由unsafe -> pool 释放
//	tmm.staticPool.ReleaseStorageMemory(numBytes)
//}

// Destory 释放所有分配的内存
func (tmm *TaskMemoryManager) CleanAllAllocatedMemory() {
	//for size, addrValue := range tmm.pageTable {
	//	tmm.memAllocator.Free(unsafe.Pointer(addrValue), size)
	//}
	//tmm.staticPool.ResetPoolUsed()
}

func InitTaskMemoryManager(config *MemoryConfig) *TaskMemoryManager {
	manager := new(TaskMemoryManager)
	manager.staticPool = newStaticMemoryManage(config)
	manager.memAllocator = manager.staticPool.DynamicMemAllocator()

	return manager
}

//var TaskMemoryManagerInstance *TaskMemoryManager = InitTaskMemoryManager(&MemoryConfig{StorageMem: "5G", ShuffleMem: "5G", IntersectionMem: "5G"})
