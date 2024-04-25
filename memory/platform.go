package memory

import "unsafe"

type platform struct{}

func (p *platform) allocate(n uintptr) unsafe.Pointer {
	return sysAllocOS(n)
}

func (p *platform) free(address unsafe.Pointer, n uintptr) {
	sysFreeOS(address, n)
}

var platformInstance = new(platform)

/**
通过使用go:linkname编译指令，链接到runtime包中申请内存的私有方法
*/
//go:linkname sysAllocOS runtime.sysAllocOS
//go:noescape
func sysAllocOS(n uintptr) unsafe.Pointer

//go:linkname sysFreeOS runtime.sysFreeOS
//go:noescape
func sysFreeOS(v unsafe.Pointer, n uintptr)

//go:linkname memmoveNoHeapPointers reflect.memmove
//go:noescape
func memmoveNoHeapPointers(dst, src unsafe.Pointer, n uintptr)
