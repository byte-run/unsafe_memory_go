package tcmallocgo

import (
	"github.com/byte-run/unsafe_mem_go/utils"
	"sync"
	"unsafe"
)

type MemoryPool struct {
	mu       sync.Mutex
	PoolSize uintptr
	used     uintptr
}

//func (p *MemoryPool) PoolSize() int {
//	p.mu.Lock()
//	defer p.mu.Unlock()
//	return p.poolSize
//}

func (p *MemoryPool) MemoryFree() uintptr {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.PoolSize - p.used
}

//func (p *MemoryPool) SetPoolSize(size int) {
//	p.mu.Lock()
//	defer p.mu.Unlock()
//
//	p.PoolSize = size
//}

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

// AcquireMemory 申请内存，有点类似提交内存大小的申请，看pool limit够不够
func (pool *storageMemoryPool) AcquireMemory(numBytes uintptr) (uintptr, error) {
	if numBytes == 0 {
		return numBytes, utils.AcquireMemoryBytesZeroError
	}
	pool.mu.Lock()
	defer pool.mu.Unlock()

	// pool retain mem
	grant := utils.Min(numBytes, pool.MemoryFree())
	// 如果pool有空间的话
	return grant, nil
}

func (pool *storageMemoryPool) ReleaseMemory(numBytes uintptr) {

}

// shuffleMemoryPool shuffle时的内存控制，主要用于bucket数据
type shuffleMemoryPool struct {
	MemoryPool
	chuckMap map[int]unsafe.Pointer // 内存块
}

// intersectionMemoryPool 交集计算时的内存控制
type intersectionMemoryPool struct {
	MemoryPool
	chuckMap map[int]unsafe.Pointer
}
