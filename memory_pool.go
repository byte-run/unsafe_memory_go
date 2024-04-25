package tcmallocgo

import (
	"sync"
	"unsafe"
)

type MemoryPool struct {
	mu       sync.Mutex
	poolSize int
	used     int
}

func (p *MemoryPool) PoolSize() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.poolSize
}

func (p *MemoryPool) MemoryFree() int {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.poolSize - p.used
}

func (p *MemoryPool) SetPoolSize(size int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.poolSize = size
}

// storageMemoryPool 存储内存池，管理元数据
type storageMemoryPool struct {
	MemoryPool
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
