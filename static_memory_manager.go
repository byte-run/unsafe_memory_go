package tcmallocgo

import (
	"github.com/byte-run/unsafe_mem_go/memory"
	"github.com/byte-run/unsafe_mem_go/utils"
	"log"
)

var emptyValue = uintptr(0)

// unsafeMemory 业务分块内存的内存管理器
type staticMemoryManage struct {
	storagePool      *storageMemoryPool
	shufflePool      *shuffleMemoryPool
	intersectionPool *intersectionMemoryPool

	MemoryMode MemoryMode
}

// AcquireStorageMemory 请求storage部分的内存
func (mem *staticMemoryManage) AcquireStorageMemory(numBytes uintptr) (bool, utils.MemPoolWarn, error) {
	mem.storagePool.Lock()
	defer mem.storagePool.Unlock()

	if numBytes > mem.storagePool.PoolSize {
		return false, nil, utils.StoragePoolOutOfMemoryError
	}

	memory, err := mem.storagePool.AcquireMemory(numBytes)
	if err != nil {
		return false, nil, err
	}
	poolStatus := mem.storagePool.CheckPoolCapacity()
	return memory, poolStatus, nil
}

func (mem *staticMemoryManage) ReleaseStorageMemory(numBytes uintptr) {
	mem.storagePool.Lock()
	defer mem.storagePool.Unlock()

	mem.storagePool.ReleaseMemory(numBytes)
}

func (mem *staticMemoryManage) ReleaseAllStorageMemory() {
	mem.storagePool.Lock()
	defer mem.storagePool.Unlock()

	mem.storagePool.ReleaseAllMemory()
}

// acquireShuffleMemory
func (mem *staticMemoryManage) acquireShuffleMemory(numBytes uintptr) (uintptr, utils.MemPoolWarn, error) {
	mem.shufflePool.Lock()
	defer mem.shufflePool.Unlock()

	if numBytes > mem.shufflePool.PoolSize {
		return emptyValue, nil, utils.ShufflePoolOutOfMemoryError
	}

	memory, err := mem.shufflePool.acquireMemory(numBytes)
	if err != nil {
		return emptyValue, nil, err
	}
	poolStatus := mem.shufflePool.CheckPoolCapacity()

	return memory, poolStatus, nil
}

func (mem *staticMemoryManage) ReleaseShuffleMemory(numBytes uintptr) {
	mem.shufflePool.Lock()
	defer mem.shufflePool.Unlock()

	mem.shufflePool.ReleaseMemory(numBytes)
}

// acquireIntersectionMemory Intersection过程中需要什么的内存
func (mem *staticMemoryManage) acquireIntersectionMemory(numBytes uintptr) (uintptr, utils.MemPoolWarn, error) {
	mem.intersectionPool.Lock()
	defer mem.intersectionPool.Unlock()

	// param check
	if numBytes > mem.intersectionPool.PoolSize {
		return emptyValue, nil, utils.ShufflePoolOutOfMemoryError
	}

	// attempt request numBytes memory space from the business pool
	memory, err := mem.intersectionPool.AcquireMemory(numBytes)
	if err != nil {
		return emptyValue, nil, err
	}
	// then, Checking the business pool usage
	poolStatus := mem.intersectionPool.CheckPoolCapacity()

	return memory, poolStatus, nil
}

func (mem *staticMemoryManage) ReleaseIntersectionMemory(numBytes uintptr) {
	mem.intersectionPool.Lock()
	defer mem.intersectionPool.Unlock()

	mem.intersectionPool.ReleaseMemory(numBytes)
}

func (mem *staticMemoryManage) ResetPoolUsed() {
	mem.storagePool.used = 0
	mem.shufflePool.used = 0
	mem.intersectionPool.used = 0
}

func (mem *staticMemoryManage) setMemoryMode() {
	mem.MemoryMode = MemoryMode_OffHeap
}

// allocator 根据设定采用不同的内存分配实现
func (mem *staticMemoryManage) DynamicMemAllocator() memory.MemAllocator {
	switch mem.MemoryMode {
	case MemoryMode_OnHeap:
		return memory.UnsafeGo
	case MemoryMode_OffHeap:
		return memory.UnsafeC
	default:
		log.Fatalf("Not Supported MemoryMode: %v", mem.MemoryMode)
		return nil
	}
}

// newStaticMemoryManage init
func newStaticMemoryManage(config *MemoryConfig) *staticMemoryManage {
	memManager := new(staticMemoryManage)

	memManager.setMemoryMode()
	// init memPool
	memManager.storagePool = &storageMemoryPool{
		&MemoryPool{PoolSize: config.GetStorageMemBytes()},
	}
	memManager.shufflePool = &shuffleMemoryPool{
		MemoryPool: &MemoryPool{
			PoolSize: config.GetShuffleMemBytes(),
		},
	}
	memManager.intersectionPool = &intersectionMemoryPool{
		MemoryPool: &MemoryPool{PoolSize: config.GetIntersectionMemBytes()},
	}

	return memManager
}
