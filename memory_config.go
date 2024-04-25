package tcmallocgo

type MemoryConfig struct {
	storageMem      string
	ShuffleMem      string
	IntersectionMem string
}

func (conf MemoryConfig) GetStorageMemBytes() int {
	// TODO waiting for finish

	return 0
}

// GetShuffleMemBytes 获取shuffle块的内存容量，单位是bytes
func (conf MemoryConfig) GetShuffleMemBytes() int {
	// TODO waiting for finish

	return 0
}

func (conf MemoryConfig) GetIntersectionMemBytes() int {
	// TODO waiting for finish

	return 0
}
