package tcmallocgo

import (
	"strconv"
	"strings"
)

const (
	MB_Factor = 1024 * 1024
	GB_Factor = 1024 * 1024 * 1024
)

type MemoryConfig struct {
	StorageMem      string
	ShuffleMem      string
	IntersectionMem string
}

func (conf MemoryConfig) GetStorageMemBytes() uintptr {
	numBytes, _ := convertMemBytes(conf.StorageMem)
	return numBytes
}

// GetShuffleMemBytes 获取shuffle块的内存容量，单位是bytes
func (conf MemoryConfig) GetShuffleMemBytes() uintptr {
	numBytes, _ := convertMemBytes(conf.ShuffleMem)
	return uintptr(numBytes)
}

func (conf MemoryConfig) GetIntersectionMemBytes() uintptr {
	numBytes, _ := convertMemBytes(conf.IntersectionMem)
	return uintptr(numBytes)
}

// 先限定GB单位
func convertMemBytes(memStr string) (uintptr, error) {
	split := strings.Split(memStr, "G")
	if len(split) != 2 {

	}
	v, err := strconv.Atoi(split[0])
	if err != nil {
		return 0, err
	}
	// TODO 待续

	return uintptr(v * GB_Factor), nil
}
