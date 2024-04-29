package tcmallocgo

import (
	"github.com/byte-run/unsafe_mem_go/utils"
	"sync"
)

var emptyValue = uintptr(0)

// unsafeMemory 业务分块内存的内存管理器
type staticMemoryManage struct {
	storagePool      *storageMemoryPool
	shufflePool      *shuffleMemoryPool
	intersectionPool *intersectionMemoryPool

	mu sync.RWMutex
}

// AcquireStorageMemory 请求storage部分的内存
func (mem *staticMemoryManage) AcquireStorageMemory(numBytes uintptr) (bool, utils.MemPoolWarn, error) {
	mem.mu.RLock()
	defer mem.mu.RUnlock()

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
	mem.mu.Lock()
	defer mem.mu.Unlock()

	mem.storagePool.ReleaseMemory(numBytes)
}

func (mem *staticMemoryManage) ReleaseAllStorageMemory() {
	mem.mu.Lock()
	defer mem.mu.Unlock()

	mem.storagePool.ReleaseAllMemory()
}

// acquireShuffleMemory
func (mem *staticMemoryManage) acquireShuffleMemory(numBytes uintptr) (uintptr, utils.MemPoolWarn, error) {
	mem.mu.RLock()
	defer mem.mu.RUnlock()

	if numBytes > mem.shufflePool.PoolSize {
		return emptyValue, nil, utils.ShufflePoolOutOfMemoryError
	}

	memory, err := mem.shufflePool.AcquireMemory(numBytes)
	if err != nil {
		return emptyValue, nil, err
	}
	poolStatus := mem.shufflePool.CheckPoolCapacity()

	return memory, poolStatus, nil
}

func (mem *staticMemoryManage) ReleaseShuffleMemory(numBytes uintptr) {
	mem.mu.Lock()
	defer mem.mu.Unlock()

	mem.shufflePool.ReleaseMemory(numBytes)
}

// acquireIntersectionMemory Intersection过程中需要什么的内存
func (mem *staticMemoryManage) acquireIntersectionMemory(numBytes uintptr) (uintptr, utils.MemPoolWarn, error) {
	mem.mu.Lock()
	defer mem.mu.Unlock()

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
	mem.mu.Lock()
	defer mem.mu.Unlock()

	mem.intersectionPool.ReleaseMemory(numBytes)
}

func (mem *staticMemoryManage) ResetPoolUsed() {
	mem.storagePool.used = 0
	mem.shufflePool.used = 0
	mem.intersectionPool.used = 0
}

// newStaticMemoryManage init
func newStaticMemoryManage(config *MemoryConfig) *staticMemoryManage {
	unsafeMem := new(staticMemoryManage)

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