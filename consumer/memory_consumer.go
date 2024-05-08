package consumer

import (
	"github.com/byte-run/unsafe_mem_go"
	"github.com/byte-run/unsafe_mem_go/memory"
	"github.com/byte-run/unsafe_mem_go/utils"
)

// 对rock业务来说

type MemoryConsumer interface {
	AllocateArray(size uintptr) uintptr
	FreeArray(memBlock *memory.MemBlock)
	AllocatePage(numBytes uintptr) (*memory.MemBlock, error)
	FreePage(page *memory.MemBlock)
	//acquireMemory(size uintptr) uintptr
	//releaseMemory(size uintptr)
	FreeMemory()

	GetStage() tcmallocgo.CalcStage
}

type memoryConsumer struct {
	taskMemoryManager *tcmallocgo.TaskMemoryManager
	pageSize          uint64
	used              uintptr
}

func (consumer *memoryConsumer) AllocateArray(size uintptr) uintptr {
	// TODO waiting to finish

	return 0
}

func (consumer *memoryConsumer) FreeArray(memBlock *memory.MemBlock) {
	// TODO waiting to finish
	consumer.FreePage(memBlock)
}

func (consumer *memoryConsumer) AllocatePage(numBytes uintptr, obj MemoryConsumer) (*memory.MemBlock, error) {

	page, err := consumer.taskMemoryManager.AllocatePage(numBytes, obj)
	if err != nil {
		return nil, err
	}
	if page == nil || page.Size() < numBytes {
		// TODO log records the kind of pool and out to log file
		return nil, utils.PoolOutOfMemoryError
	}
	consumer.used += page.Size()
	return page, nil
}

func (consumer *memoryConsumer) FreePage(page *memory.MemBlock) {
	consumer.used -= page.Size()
	consumer.taskMemoryManager.FreeBlockPage(page)
}

func (consumer *memoryConsumer) acquireMemory(size uintptr) uintptr {
	//TODO waiting to finish
	return 0
}

func (consumer *memoryConsumer) releaseMemory(size uintptr) {
	// TODO waiting to finish
}

// 释放当前consumer下的所有内存占用
func (consumer *memoryConsumer) FreeMemory() {
	// TODO waiting to finish
}

//func (consumer *MemoryConsumer) ThrowOomError() uintptr {}

//func NewMemoryConsumer(manager tcmallocgo.TaskMemoryManager) MemoryConsumer {
//	var memoryConsumer = new(memoryConsumer)
//	memoryConsumer.taskMemoryManager = &manager
//	memoryConsumer.used = 0
//
//	return memoryConsumer
//}
