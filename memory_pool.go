package tcmallocgo

import (
	"github.com/byte-run/unsafe_mem_go/utils"
	"sync"
)

type MemoryPool struct {
	mu       sync.Mutex
	PoolSize uintptr
	used     uintptr
}

func (p *MemoryPool) MemoryFree() uintptr {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.PoolSize - p.used
}

func (p *MemoryPool) IncrementPoolSize(size uintptr) {
	p.mu.Lock()
	defer p.mu.Unlock()
	// check

	p.used -= size
}

type memChuck struct {
}

// storageMemoryPool 存储内存池，管理元数据
type storageMemoryPool struct {
	MemoryPool
}

func (pool *storageMemoryPool) PoolName() string {
	return "storage"
}

// AcquireMemory 申请内存，有点类似提交内存大小的申请，看pool limit够不够
func (pool *storageMemoryPool) AcquireMemory(numBytes uintptr) (uintptr, error) {
	if numBytes == 0 {
		return numBytes, utils.AcquireMemoryBytesZeroError
	}
	pool.mu.Lock()
	defer pool.mu.Unlock()

	// pool retain mem
	grant := utils.Min(numBytes, pool.MemoryFree())
	// 如果pool有空间的话, 更新pool的use
	if grant == numBytes {
		pool.MemoryPool.used += numBytes
	}

	return grant, nil
}

// ReleaseMemory 释放内存
func (pool *storageMemoryPool) ReleaseMemory() {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	pool.used = 0
}

// shuffleMemoryPool shuffle时的内存控制，主要用于bucket数据
type shuffleMemoryPool struct {
	MemoryPool
	cond sync.Cond
	//chuckMap map[int]unsafe.Pointer // 内存块
}

func (pool *shuffleMemoryPool) PoolName() string {
	return "shuffle bucket"
}

func (pool *shuffleMemoryPool) AcquireMemory(numBytes uintptr) (uintptr, error) {
	pool.cond.L.Lock()
	defer pool.cond.L.Unlock()

	if numBytes == 0 {
		return numBytes, utils.AcquireMemoryBytesZeroError
	}

	toGrant := utils.Min(numBytes, pool.MemoryFree())
	if toGrant < numBytes {
		pool.cond.Wait()
	}
	return toGrant, nil
}

func (pool *shuffleMemoryPool) ReleaseMemory(numBytes uintptr) {
	pool.cond.L.Lock()
	defer pool.cond.L.Unlock()

	pool.used -= numBytes
	pool.cond.Broadcast()
}

// intersectionMemoryPool 交集计算时的内存控制
type intersectionMemoryPool struct {
	MemoryPool
	cond sync.Cond
	//chuckMap map[int]unsafe.Pointer
}

func (pool *intersectionMemoryPool) PoolName() string {
	return "intersection"
}

func (pool *intersectionMemoryPool) AcquireMemory(numBytes uintptr) (uintptr, error) {
	pool.cond.L.Lock()
	defer pool.cond.L.Unlock()

	if numBytes == 0 {
		return numBytes, utils.AcquireMemoryBytesZeroError
	}

	toGrant := utils.Min(numBytes, pool.MemoryFree())
	if toGrant < numBytes {
		pool.cond.Wait()
	}
	return toGrant, nil
}

func (pool *intersectionMemoryPool) ReleaseMemory(numBytes uintptr) {
	pool.cond.L.Lock()
	defer pool.cond.L.Unlock()

	pool.used -= numBytes
	pool.cond.Broadcast()
}
