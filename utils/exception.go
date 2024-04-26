package utils

import "fmt"

type MemError struct {
	msg string
}

func (e MemError) Error() string {
	return fmt.Sprintf("error: %s", e.msg)
}

// 定义error常量
var (
	//	内存不足场景

	PlatformOutOfMemoryError         = MemError{msg: "platform out of memory"}
	PoolOutOfMemoryError             = MemError{msg: "pool out of memory"}
	StoragePoolOutOfMemoryError      = MemError{msg: "storage pool out of memory"}
	ShufflePoolOutOfMemoryError      = MemError{msg: "shuffle pool out of memory"}
	IntersectionPoolOutOfMemoryError = MemError{msg: "intersection pool out of memory"}

	// 参数异常场景
	AcquireMemoryBytesZeroError = MemError{msg: "acquire memory bytes is zero"}
)
