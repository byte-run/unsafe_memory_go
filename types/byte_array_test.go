package types

import (
	"github.com/byte-run/unsafe_mem_go/memory"
	"reflect"
	"testing"
	"unsafe"
)

//func TestByteSize(t *testing.T) {
//	emptyByteArray := ByteArray{}
//	emptySize := unsafe.Sizeof(emptyByteArray)
//	t.Logf("empty byteArray size: %d", emptySize)
//
//	byteArray := ByteArray{10, 5, nil}
//	valueSize := unsafe.Sizeof(byteArray)
//	t.Logf("value byteArray size: %d", valueSize)
//
//	//
//}

type StructToBytes struct {
	data int64
	str  string
}

type SliceMock struct {
	addr uintptr
	len  int
	cap  int
}

type StructToByte struct {
	data int64
	str  any
}

func TestGoType(t *testing.T) {
	pointerVar := unsafe.Pointer(uintptr(int8(1)))
	pointerSize := unsafe.Sizeof(pointerVar)
	t.Logf("pointer size: %d", pointerSize)
}

func TestGoSizeof(t *testing.T) {
	// array or slice
	var testArray = [5]int{1, 2, 3, 4, 5} // int type: 8 bits
	testArraySize := unsafe.Sizeof(testArray)
	t.Logf("testArray size: %d", testArraySize) // testArray size: 40 = 8 * 5
	/*
		SliceHeader size = 24
	*/
	var testSlice = make([]int, 5)
	testSlice = append(testSlice, 1, 2, 3)
	testSliceSize := unsafe.Sizeof(testSlice)
	t.Logf("testSlice size: %d", testSliceSize) // testSlice size: 24 =

	testSlice = append(testSlice, 4, 5, 6)
	testSliceSize = unsafe.Sizeof(testSlice)
	t.Logf("testSlice size: %d", testSliceSize) // testSlice size: 24 =
}

func TestStructToBytes(t *testing.T) {
	strValue := "123"
	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&strValue))
	stringHeaderData := stringHeader.Data
	t.Log(stringHeaderData)

	var testStruct = &StructToBytes{data: 10000003, str: "123"}
	Len := unsafe.Sizeof(*testStruct)
	testBytes := &SliceMock{
		addr: uintptr(unsafe.Pointer(testStruct)),
		cap:  int(Len),
		len:  int(Len),
	}

	data := *(*[]byte)(unsafe.Pointer(testBytes))
	t.Logf("[]byte is: %v", data)

}

func TestBytesToStruct(t *testing.T) {
	allocator := memory.CMemAllocator{}
	addr, err := allocator.Allocate(24 * 5) //
	if err != nil {
		t.Logf("allocate err: %v", err)
	}
	//defer allocator.Free(addr, 0)

	// 内存切换为Go切片
	var slice []StructToBytes
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	sliceHeader.Data = uintptr(unsafe.Pointer(addr))
	sliceHeader.Len = 0
	sliceHeader.Cap = 6

	t.Logf("slice header data: %p\n", unsafe.Pointer(sliceHeader.Data))

	slice = append(slice, StructToBytes{data: 10000000, str: "123"})
	t.Logf("slice header data: %p\n", unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&slice)).Data))

	t.Log(slice)
	t.Log(slice[0])
	structModel := slice[0]
	t.Logf("structModel addr %p, data: %v, str: %v\n", &structModel, structModel.data, structModel.str)

	// for
	for i := 1; i < 5; i++ {
		intData := 10000000 + i
		slice = append(slice, StructToBytes{data: int64(intData), str: "123"})
	}

	// 扩容

	t.Logf("slice header data: %p\n", unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&slice)).Data))

	allocator.Free(addr, 0)

	t.Log(slice[0])

}

func TestBytesToStruct_OutOfBound(t *testing.T) {
	// 1.尝试赋值给slice
	allocator := memory.CMemAllocator{}
	addr, err := allocator.Allocate(24 * 5) //120 bytes
	if err != nil {
		t.Logf("allocate err: %v", err)
	}
	defer allocator.Free(addr, 0)

	// 内存切换为Go切片
	var slice []StructToBytes
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	sliceHeader.Data = uintptr(unsafe.Pointer(addr))
	sliceHeader.Len = 0
	sliceHeader.Cap = 6

	t.Logf("slice header data: %p\n", unsafe.Pointer(sliceHeader.Data)) // 数据地址

	slice = append(slice, StructToBytes{data: 10000000, str: "123"})
	t.Logf("slice header data: %p\n", unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&slice)).Data)) // 数据地址是否变动

	// 填满slice
	t.Log(slice)
	for i := 1; i < 6; i++ {
		intData := 10000000 + i
		slice = append(slice, StructToBytes{data: int64(intData), str: "123"})
	}
	t.Logf("slice header data: %p\n", unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&slice)).Data)) // 数据地址是否变动

	// 遍历,访问最后一个元素
	//lastElem := slice[15]
	//t.Logf("the last elem: %v", lastElem)

	// 探寻元素存放
	// 遍历地址
	var byteCon []byte
	for i := 0; i < 120; i++ {
		byteCon = append(byteCon, *(*byte)(unsafe.Pointer(uintptr(addr) + uintptr(i))))
	}
	t.Logf("byteCon is: %v", byteCon)
}

func TestStructToByte(t *testing.T) {
	strValue := "123"
	stringHeader := (*[]byte)(unsafe.Pointer(&strValue))
	//stringHeaderData := stringHeader.Data
	t.Log(stringHeader)

	var testStruct = &StructToByte{data: 10000003, str: []byte("123")}
	Len := unsafe.Sizeof(*testStruct)
	testBytes := &SliceMock{
		addr: uintptr(unsafe.Pointer(testStruct)),
		cap:  int(Len),
		len:  int(Len),
	}

	data := *(*[]byte)(unsafe.Pointer(testBytes))
	t.Logf("[]byte is: %v", data)

}

func TestBytesToStruct_multiDim(t *testing.T) {
	// 1.尝试赋值给slice
	allocator := memory.CMemAllocator{}
	addr, err := allocator.Allocate(24 * 5) //120 bytes
	if err != nil {
		t.Logf("allocate err: %v", err)
	}
	defer allocator.Free(addr, 0)

	// 内存切换为Go切片
	var slice [][]int
	// 大小计算:
	var oneSlice []int
	slice = append(slice, oneSlice)
	addr1, err := allocator.Allocate(24 * 5) //120 bytes
	if err != nil {
		t.Logf("allocate err: %v", err)
	}
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&oneSlice))
	sliceHeader.Data = uintptr(unsafe.Pointer(addr1))
	sliceHeader.Len = 0
	sliceHeader.Cap = 6
	t.Logf("slice header data: %p\n", unsafe.Pointer(sliceHeader.Data))
	slice[0] = append(slice[0], 1, 2)

	var twoSlice []int
	slice = append(slice, twoSlice)
	addr2, err := allocator.Allocate(24 * 5) //120 bytes
	if err != nil {
		t.Logf("allocate err: %v", err)
	}
	sliceHeader2 := (*reflect.SliceHeader)(unsafe.Pointer(&twoSlice))
	sliceHeader2.Data = uintptr(unsafe.Pointer(addr2))
	sliceHeader2.Len = 0
	sliceHeader2.Cap = 6
	t.Logf("slice header data: %p\n", unsafe.Pointer(sliceHeader2.Data))
	slice[1] = append(slice[1], 3, 4)

	t.Logf("slice : %v", slice)

	//(*reflect.SliceHeader)(unsafe.Pointer(&oneSlice)).Data
	sliceHeader.Data = uintptr(0)
	allocator.Free(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&oneSlice)).Data), 0)
	allocator.Free(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&twoSlice)).Data), 0)

	t.Logf("%v", slice[0][0])

}
