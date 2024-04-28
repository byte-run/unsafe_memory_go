package memory

// C代码集成

/*
#include <string.h>
#include <stdlib.h>

void *customMalloc(size_t size) {
	return malloc(size);
}
*/
import "C"
import (
	"github.com/byte-run/unsafe_mem_go/utils"
	"unsafe"
)

type platformC struct{}

// malloc 申请
func (p platformC) allocate(size uint) (unsafe.Pointer, error) {
	addr := C.customMalloc(C.size_t(size))
	if addr == nil {
		return nil, utils.PlatformOutOfMemoryError
	}
	return unsafe.Pointer(addr), nil
}

// free 释放
func (p platformC) free(addr unsafe.Pointer) {
	C.free(addr)
}

// memMove 移动
func (p platformC) memMove(dst, src unsafe.Pointer, length uintptr) {
	C.memmove(dst, src, C.size_t(length))
}

// memCpy 拷贝
func (p platformC) memCpy(dst, src unsafe.Pointer, length uintptr) {
	//C.CBytes()
	C.memcpy(dst, src, C.size_t(length))
}

var platformCInstance = new(platformC)
