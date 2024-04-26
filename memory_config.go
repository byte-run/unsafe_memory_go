package tcmallocgo

import (
	"strconv"
	"strings"
)

type MemoryConfig struct {
	StorageMem      string
	ShuffleMem      string
	IntersectionMem string
}

func (conf MemoryConfig) GetStorageMemBytes() uintptr {
	// TODO waiting for finish

	return 0
}

// GetShuffleMemBytes 获取shuffle块的内存容量，单位是bytes
func (conf MemoryConfig) GetShuffleMemBytes() uintptr {
	// TODO waiting for finish

	return 0
}

func (conf MemoryConfig) GetIntersectionMemBytes() uintptr {
	// TODO waiting for finish

	return 0
}

// 先限定GB单位
func convertMemBytes(memStr string) (int, error) {
	split := strings.Split(memStr, "GB")
	if len(split) != 2 {

	}
	v, err := strconv.Atoi(split[0])
	if err != nil {
		return 0, err
	}
	// TODO 待续
	return v, nil
}
