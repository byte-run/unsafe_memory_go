package tcmallocgo

// 对rock业务来说
//type MemoryConsumer interface {
//	acquireMemory(size uintptr) uintptr
//	releaseMemory(size uintptr)
//}

// TODO 为array或slice申请空间, 返回值未定
//func (consumer *MemoryConsumer) allocateArray(numByte uintptr) {}

//type ArrayMemConsumer interface {
//	allocateArray(numBytes uintptr) uintptr
//	freeArray(numBytes uintptr)
//}

type MemoryConsumer struct {
	taskMemoryManager *TaskMemoryManager
	pageSize          uint64
	used              uint64
}

func (consumer *MemoryConsumer) AllocateArray(size uintptr) uintptr {
	// TODO waiting to finish
	return 0
}

func (consumer *MemoryConsumer) FreeArray(size uintptr) {
	// TODO waiting to finish
}

func (consumer *MemoryConsumer) acquireMemory(size uintptr) uintptr {
	//TODO waiting to finish
	return 0
}

func (consumer *MemoryConsumer) releaseMemory(size uintptr) {
	// TODO waiting to finish
}

// 释放当前consumer下的所有内存占用
func (consumer *MemoryConsumer) FreeMemory() {
	// TODO waiting to finish
}

//func (consumer *MemoryConsumer) ThrowOomError() uintptr {}

func NewMemoryConsumer(manager TaskMemoryManager) *MemoryConsumer {
	var memoryConsumer = new(MemoryConsumer)
	memoryConsumer.taskMemoryManager = &manager
	memoryConsumer.used = 0

	return memoryConsumer
}
