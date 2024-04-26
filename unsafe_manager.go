package tcmallocgo

import (
	"github.com/byte-run/unsafe_mem_go/memory"
	"github.com/byte-run/unsafe_mem_go/utils"
	"sync"
	"unsafe"
)

var emptyValue = uintptr(0)

type unsafeMemory struct {
	storagePool      *storageMemoryPool
	shufflePool      *shuffleMemoryPool
	intersectionPool *intersectionMemoryPool

	mu sync.RWMutex
}

// AcquireStorageMemory 请求storage部分的内存
func (mem *unsafeMemory) AcquireStorageMemory(numBytes uintptr) (uintptr, error) {
	mem.mu.RLock()
	defer mem.mu.RUnlock()

	acquireMemory, err := mem.storagePool.AcquireMemory(numBytes)
	if err != nil {
		return emptyValue, err
	}

	if acquireMemory < numBytes {
		return emptyValue, utils.StoragePoolOutOfMemoryError
	}
	return acquireMemory, nil
}

// acquireShuffleMemory
func (mem unsafeMemory) acquireShuffleMemory(numBytes uintptr) (unsafe.Pointer, error) {
	mem.mu.RLock()
	defer mem.mu.RUnlock()

	return nil, nil
}

// acquireIntersectionMemory Intersection过程中需要什么的内存
func (mem unsafeMemory) acquireIntersectionMemory(numBytes uintptr) (unsafe.Pointer, error) {
	return nil, nil
}

func newUnsafeMemory(config MemoryConfig) *unsafeMemory {
	unsafeMem := new(unsafeMemory)

	// init memPool
	unsafeMem.storagePool = &storageMemoryPool{
		MemoryPool{PoolSize: config.GetStorageMemBytes()},
	}
	unsafeMem.shufflePool = &shuffleMemoryPool{
		MemoryPool: MemoryPool{
			PoolSize: config.GetShuffleMemBytes(),
		},
	}
	unsafeMem.intersectionPool = &intersectionMemoryPool{
		MemoryPool: MemoryPool{PoolSize: config.GetIntersectionMemBytes()},
	}

	return unsafeMem
}

type UnsafeManager struct {
	unsafeMemory
	memAllocator memory.MemAllocator
}

func (manager *UnsafeManager) Allocate(numBytes uintptr) (unsafe.Pointer, error) {
	return nil, nil
}

// 对外接口，单实列
var _unsafe *unsafeMemory = nil

func newUnsafeManager(memLimits ...int) *UnsafeManager {
	mem := new(UnsafeManager)

	mem.memAllocator = dynamicMemAllocator("C")

	return mem
}

// allocator 根据设定采用不同的内存分配实现
func dynamicMemAllocator(allocMode string) memory.MemAllocator {
	if allocMode == "C" {
		return memory.UnsafeC
	}
	return memory.UnsafeGo

}
