package memory

// C代码集成

/*
#include <string.h>
#include <stdlib.h>
*/
import "C"
import "unsafe"

// malloc 申请
func malloc(size int) unsafe.Pointer {
	return C.malloc(C.size_t(size))
}

// free 释放
func free(p unsafe.Pointer) {
	C.free(p)
}

// memMove 移动
func memMove(dst, src unsafe.Pointer, length uintptr) {
	C.memmove(dst, src, C.size_t(length))
}

// memCpy 拷贝
func memCpy(dst, src unsafe.Pointer, length uintptr) {
	//C.CBytes()
	C.memcpy(dst, src, C.size_t(length))
}

type platformC struct{}

func (p platformC) allocate(size uintptr) unsafe.Pointer {
	return malloc(int(size))
}

func (p platformC) free(addr unsafe.Pointer) {
	free(addr)
}

func (p platformC) memCpy(dst, src unsafe.Pointer, length uintptr) {}

var platformCInstance = new(platformC)
