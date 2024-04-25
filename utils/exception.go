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
	PlatformOutOfMemoryError         = MemError{msg: "platform out of memory"}
	PoolOutOfMemoryError             = MemError{msg: "pool out of memory"}
	StoragePoolOutOfMemoryError      = MemError{msg: "storage pool out of memory"}
	ShufflePoolOutOfMemoryError      = MemError{msg: "shuffle pool out of memory"}
	IntersectionPoolOutOfMemoryError = MemError{msg: "intersection pool out of memory"}
)
