package tcmallocgo

// 对rock业务来说
type MemoryConsumer struct {
}

// TODO 为array或slice申请空间, 返回值未定
func (consumer *MemoryConsumer) allocateArray(numByte uintptr) {}
