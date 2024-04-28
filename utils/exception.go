package utils

import "fmt"

type MemError struct {
	msg string
}

func (e MemError) Error() string {
	return fmt.Sprintf("error: %s", e.msg)
}

type MemWarn struct {
	msg string
}

func (e MemWarn) Error() string {
	return fmt.Sprintf("error: %s", e.msg)
}

// 内存不足场景
var (
	PlatformOutOfMemoryError         = MemError{msg: "platform out of memory"}
	PoolOutOfMemoryError             = MemError{msg: "pool out of memory"}
	StoragePoolOutOfMemoryError      = MemError{msg: "storage pool out of memory"}
	ShufflePoolOutOfMemoryError      = MemError{msg: "shuffle pool out of memory"}
	IntersectionPoolOutOfMemoryError = MemError{msg: "intersection pool out of memory"}
)

// 参数问题场景
var (
	AcquireMemoryBytesZeroError = MemError{msg: "acquire memory bytes is zero"}
)

// 警告场景
var (
	StoragePoolLevelOneMemoryWarning      = MemError{msg: "storage memory pool use 80%"}
	StoragePoolLevelTwoMemoryWarning      = MemError{msg: "storage memory pool use 90%"}
	ShufflePoolLevelOneMemoryWarning      = MemError{msg: "shuffle memory pool use 80%"}
	ShufflePoolLevelTwoMemoryWarning      = MemError{msg: "shuffle memory pool use 90%"}
	IntersectionPoolLevelOneMemoryWarning = MemError{msg: "intersection memory pool use 80%"}
	IntersectionPoolLevelTwoMemoryWarning = MemError{msg: "intersection memory pool use 90%"}
)
