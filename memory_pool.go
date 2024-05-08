package tcmallocgo

import (
	"github.com/byte-run/unsafe_mem_go/utils"
	"sync"
)

const (
	PoolLevelOneFactor = 0.8
	PoolLevelTwoFactor = 0.9
)

type MemoryPool struct {
	condMu   sync.Cond
	PoolSize uintptr
	used     uintptr
}

func (p *MemoryPool) MemoryFree() uintptr {
	return p.PoolSize - p.used
}

func (p *MemoryPool) checkPoolCapacity(poolName string) utils.MemPoolWarn {
	levelOneThreshold := p.getMemoryPoolLevelOneThreshold()
	levelTwoThreshold := p.getMemoryPoolLevelTwoThreshold()

	if p.used >= levelTwoThreshold {
		return utils.MemoryPoolLevelTwoWarning{
			PoolName: poolName,
		}
	} else if p.used >= levelOneThreshold {
		return utils.MemoryPoolLevelOneWarning{
			PoolName: poolName,
		}
	}
	return nil
}

func (p *MemoryPool) getMemoryPoolLevelOneThreshold() uintptr {
	return uintptr(float64(p.PoolSize) * PoolLevelOneFactor)
}

func (p *MemoryPool) getMemoryPoolLevelTwoThreshold() uintptr {
	return uintptr(float64(p.PoolSize) * PoolLevelTwoFactor)
}

func (p *MemoryPool) Lock() {
	if p.condMu.L == nil {
		p.condMu = sync.Cond{L: &sync.Mutex{}}
	}
	p.condMu.L.Lock()
}

func (p *MemoryPool) Unlock() {
	if p.condMu.L == nil {
		// 说明Locker都没init
		return
	}
	p.condMu.L.Unlock()
}

// storageMemoryPool 存储内存池，管理元数据
type storageMemoryPool struct {
	*MemoryPool
}

func (pool *storageMemoryPool) PoolName() string {
	return "storage"
}

// AcquireMemory 申请内存，有点类似提交内存大小的申请，看pool limit够不够
func (pool *storageMemoryPool) AcquireMemory(numBytes uintptr) (uintptr, error) {
	if numBytes == 0 {
		return 0, utils.AcquireMemoryBytesZeroError
	}
	pool.condMu.L.Lock()
	defer pool.condMu.L.Unlock()

	// pool retain mem
	grant := utils.Min(numBytes, pool.MemoryFree())
	// 如果pool有空间的话, 更新pool的use
	//if grant == numBytes {
	//	pool.MemoryPool.used += numBytes
	//	return true, nil
	//}

	return grant, utils.StoragePoolOutOfMemoryError
}

// ReleaseMemory 释放内存
func (pool *storageMemoryPool) ReleaseMemory(numBytes uintptr) {
	//pool.condMu.L.Lock()
	//defer pool.condMu.L.Unlock()
	pool.used -= numBytes

	pool.condMu.Broadcast()
}

func (pool *storageMemoryPool) ReleaseAllMemory() {
	//pool.condMu.L.Lock()
	//defer pool.condMu.L.Unlock()

	pool.used = 0
}

func (pool *storageMemoryPool) CheckPoolCapacity() utils.MemPoolWarn {
	//pool.condMu.L.Lock()
	//defer pool.condMu.L.Unlock()

	return pool.checkPoolCapacity(pool.PoolName())
}

//type executionMemoryPool struct {
//	MemoryPool
//	cond sync.Cond
//}

// shuffleMemoryPool shuffle时的内存控制，主要用于bucket数据
type shuffleMemoryPool struct {
	*MemoryPool
}

func (pool *shuffleMemoryPool) PoolName() string {
	return "shuffle bucket"
}

func (pool *shuffleMemoryPool) acquireMemory(numBytes uintptr) (uintptr, error) {
	//pool.cond.L.Lock()
	//defer pool.cond.L.Unlock()

	if numBytes == 0 {
		return numBytes, utils.AcquireMemoryBytesZeroError
	}
	var counter = 0
	for {
		toGrant := utils.Min(numBytes, pool.MemoryFree())
		if toGrant < numBytes {
			pool.condMu.Wait()
		} else {
			pool.used += toGrant
			//atomic.AddUintptr(&pool.used, toGrant)
			return toGrant, nil
		}
		counter++
		if counter == 5 {
			return 0, utils.ShufflePoolOutOfMemoryError
		}
	}
	return 0, nil
}

func (pool *shuffleMemoryPool) ReleaseMemory(numBytes uintptr) {
	//pool.cond.L.Lock()
	//defer pool.cond.L.Unlock()

	pool.used -= numBytes
	pool.condMu.Broadcast()
}

func (pool *shuffleMemoryPool) ReleaseAllMemory() {
	pool.used = 0
}

func (pool *shuffleMemoryPool) CheckPoolCapacity() utils.MemPoolWarn {
	//pool.condMu.Lock()
	//defer pool.condMu.Unlock()

	return pool.checkPoolCapacity(pool.PoolName())
}

// intersectionMemoryPool 交集计算时的内存控制
type intersectionMemoryPool struct {
	*MemoryPool
	//chuckMap map[int]unsafe.Pointer
}

func (pool *intersectionMemoryPool) PoolName() string {
	return "intersection"
}

func (pool *intersectionMemoryPool) AcquireMemory(numBytes uintptr) (uintptr, error) {
	//pool.Lock()
	//defer pool.Unlock()

	if numBytes == 0 {
		return numBytes, utils.AcquireMemoryBytesZeroError
	}

	counter := 0
	for {
		toGrant := utils.Min(numBytes, pool.MemoryFree())
		if toGrant < numBytes {
			pool.condMu.Wait()
		} else {
			pool.used += toGrant
			return toGrant, nil
		}
		counter++
		if counter == 5 {
			return 0, utils.IntersectionPoolOutOfMemoryError
		}
	}
	return 0, nil
}

func (pool *intersectionMemoryPool) ReleaseMemory(numBytes uintptr) {
	//pool.cond.L.Lock()
	//defer pool.cond.L.Unlock()

	pool.used -= numBytes
	pool.condMu.Broadcast()
}

func (pool *intersectionMemoryPool) CheckPoolCapacity() utils.MemPoolWarn {
	//pool.condMu.L.Lock()
	//defer pool.condMu.L.Unlock()

	return pool.checkPoolCapacity(pool.PoolName())
}

func (pool *intersectionMemoryPool) ReleaseAllMemory() {
	pool.used = 0
}
