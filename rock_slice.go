package tcmallocgo

import "C"
import (
	"fmt"
	//"gitlab.grandhoo.com/rock/rock-share/base/logger"
	"reflect"
	"time"
	"unsafe"
)

const (
	STRING = iota
	INT
	INT32
	FLOAT64
	BOOL
	TIME
	INTERFACE
)

const DefaultCapacity = 20

// RockSlice 是一个自定义的切片类型，用于管理内存
type RockSlice struct {
	Addr        unsafe.Pointer // 首元素位置
	Size        int            // 切片的当前大小
	Capacity    int            // 切片的容量
	ByteSize    int            //占用字节数，为string类型设计
	ElementLens *RockSlice     //记录每个元素的占用大小，string专用
	ElementType reflect.Type   //元素类型
}

// NewRockSlice 创建一个 RockSlice 实例
func NewRockSlice(elementType int) (*RockSlice, error) {
	return NewRockSliceWithCapacity(elementType, DefaultCapacity)
}

// NewRockSliceWithCapacity 创建一个指定容量的 RockSlice 实例
func NewRockSliceWithCapacity(elementType int, initialCapacity int) (*RockSlice, error) {
	// 根据类型编号确定 reflect.Type
	var reflectType reflect.Type
	var elementLens *RockSlice
	var err error
	switch elementType {
	case INT:
		reflectType = reflect.TypeOf(0)
	case INT32:
		reflectType = reflect.TypeOf(int32(0))
	case FLOAT64:
		reflectType = reflect.TypeOf(float64(0))
	case STRING:
		reflectType = reflect.TypeOf("")
		elementLens, err = NewRockSliceWithCapacity(INT, DefaultCapacity)
		if err != nil {
			return nil, fmt.Errorf("initial elementLens failed")
		}
	case BOOL:
		reflectType = reflect.TypeOf(true)
	case TIME:
		reflectType = reflect.TypeOf(time.Time{})
	case INTERFACE:
		reflectType = reflect.TypeOf(interface{}(nil)) //会返回nil

	default:
		return nil, fmt.Errorf("unsupported type: %d", elementType)
	}

	var address unsafe.Pointer
	//todo jl
	//1.内存告警检查

	address, err = allocMemory(reflectType, initialCapacity)
	if err != nil {
		return nil, err
	}
	//fmt.Println(uintptr(address))
	return &RockSlice{
		Addr:        address,
		Size:        0,
		Capacity:    initialCapacity,
		ElementType: reflectType,
		ByteSize:    0,
		ElementLens: elementLens,
	}, nil
}

func allocMemory(reflectType reflect.Type, capacity int) (unsafe.Pointer, error) {
	if reflectType == nil { //interface{}在添加第一个元素时确定类型再申请内存，但不支持同一列有不同类型。若要支持，则要改成类似于string的变长
		return nil, nil
	}
	totalSize := reflectType.Size() * uintptr(capacity)
	return C.malloc(C.size_t(totalSize)), nil
}

// Append 向 RockSlice 中添加新元素
func (s *RockSlice) Append(value any) (err error) {
	//判断类型
	valueType := reflect.TypeOf(value)
	if s.Size == 0 && s.ElementType == nil { //interface{}第一个元素确定类型
		s.ElementType = valueType
		s.Addr, err = allocMemory(s.ElementType, s.Capacity)
		if err != nil {
			return err
		}
		if valueType.Kind() == reflect.String { //字符串类型
			s.ElementLens, err = NewRockSliceWithCapacity(INT, DefaultCapacity)
			if err != nil {
				return err
			}
		}
	}
	if valueType.Kind() != s.ElementType.Kind() || valueType.Name() != s.ElementType.Name() {
		return fmt.Errorf("type mismatch: expected %s, got %s", s.ElementType, valueType)
	}

	// 检查是否需要扩容
	if valueType.Kind() != reflect.String && s.Size == s.Capacity {
		err = s.grow()
		if err != nil {
			return err
		}
	}

	//添加元素
	//uintptr 类型是不能存储在临时变量中的。因为从 GC 的角度来看，uintptr 类型的临时变量只是一个无符号整数，并不知道它是一个指针地址。
	valuePtr := unsafe.Pointer(uintptr(s.Addr) + s.ElementType.Size()*uintptr(s.Size))
	// 使用 unsafe.Pointer 转换类型并赋值
	switch s.ElementType.Kind() {
	case reflect.Int:
		*((*int)(valuePtr)) = value.(int)
	case reflect.Int32:
		*((*int32)(valuePtr)) = value.(int32)
	case reflect.Float64:
		*((*float64)(valuePtr)) = value.(float64)
	case reflect.String:
		byteSlice := []byte(value.(string))
		// 检查是否需要扩容
		for uintptr(s.ByteSize)+uintptr(len(byteSlice)) >= s.ElementType.Size()*uintptr(s.Capacity) {
			err = s.grow()
			if err != nil {
				return err
			}
		}
		valuePtr = unsafe.Pointer(uintptr(s.Addr) + uintptr(s.ByteSize))
		//copy((*(*[1 << 30]byte)(valuePtr))[:len(byteSlice)], byteSlice)
		for i := 0; i < len(byteSlice); i++ {
			*(*byte)(unsafe.Pointer(uintptr(valuePtr) + uintptr(i))) = byteSlice[i]
		}

		s.ByteSize += len(byteSlice)
		err = s.ElementLens.Append(len(byteSlice))
		if err != nil {
			return fmt.Errorf("record length of %v failed", value)
		}
		//if s.Size == s.Capacity {//todo怎么处理，size>capacity的情况
		//	s.Capacity++
		//}
	case reflect.Bool:
		*((*bool)(valuePtr)) = value.(bool)
	case reflect.TypeOf(time.Time{}).Kind():
		*((*time.Time)(valuePtr)) = value.(time.Time)
	default:
		return fmt.Errorf("unsupported type: %s", s.ElementType)
	}
	s.Size++
	return nil
}

/*扩展条件：
string：ByteSize + len(NewStringByteLen) >= capacity*16，capacity不由size而变化，因此size可能大于capacity
其他：size==capacity
*/
// grow 扩展切片的容量
func (s *RockSlice) grow() (err error) {
	newCapacity := s.Capacity * 2 // 倍增
	var address unsafe.Pointer

	//todo jl
	//1.内存告警检查
	totalSize := s.ElementType.Size() * uintptr(newCapacity)
	//2.申请内存
	address = C.malloc(C.size_t(totalSize))
	fmt.Printf("grow capacity from %v to %v", s.Capacity, s.Capacity*2)
	//fmt.Println(uintptr(address))
	//数据迁移
	oldSize := s.Size * int(s.ElementType.Size()) // 计算旧数据的总大小
	if s.ElementType.Kind() == reflect.String {
		oldSize = s.ByteSize
	}
	oldAddr := s.Addr
	for i := 0; i < oldSize; i++ {
		// 获取旧地址上的数据
		oldValue := *(*byte)(unsafe.Pointer(uintptr(oldAddr) + uintptr(i)))
		// 将数据写入新地址
		newValueAddr := uintptr(address) + uintptr(i)
		*(*byte)(unsafe.Pointer(newValueAddr)) = oldValue
	}
	// 释放旧的内存空间
	//todo jl
	//C.free(s.Addr)

	s.Addr = address
	s.Capacity = newCapacity
	return nil
}

func (s *RockSlice) Get(index int) (value any, err error) {
	if index < 0 || index >= s.Size {
		return nil, fmt.Errorf("index out of range")
	}
	addr := s.Addr
	//a := uintptr(Addr) + s.ElementType.Len()*uintptr(index)
	//valuePtr := unsafe.Pointer(a)
	valuePtr := unsafe.Pointer(uintptr(addr) + s.ElementType.Size()*uintptr(index))
	switch s.ElementType.Kind() {
	case reflect.Int:
		value = *((*int)(valuePtr))
	case reflect.Int32:
		value = *((*int32)(valuePtr))
	case reflect.Float64:
		value = *((*float64)(valuePtr))
	case reflect.String:
		valuePtr = addr
		for i := 0; i < index; i++ {
			elementLen, err := s.ElementLens.Get(i)
			if err != nil {
				return nil, fmt.Errorf("get length of No.%v failed", i)
			}
			valuePtr = unsafe.Pointer(uintptr(valuePtr) + uintptr(elementLen.(int)))
		}
		//fmt.Println(uintptr(valuePtr))
		elementLen, err := s.ElementLens.Get(index)
		if err != nil {
			return nil, fmt.Errorf("get length of target index %v failed", index)
		}
		dataSlice := make([]byte, elementLen.(int))
		//copy(dataSlice, (*(*[1 << 30]byte)(valuePtr))[:len.(int)])
		for i := 0; i < elementLen.(int); i++ {
			dataSlice[i] = *(*byte)(unsafe.Pointer(uintptr(valuePtr) + uintptr(i)))
		}
		value = string(dataSlice)
	case reflect.Bool:
		value = *((*bool)(valuePtr))
	case reflect.TypeOf(time.Time{}).Kind():
		value = *((*time.Time)(valuePtr))
	default:
		return nil, fmt.Errorf("unsupported type: %s", s.ElementType)
	}
	return value, nil

}

func (s *RockSlice) Len() int {
	return s.Size
}

func (s *RockSlice) Free() error {
	//todo jl
	return nil
}

// Print 打印 RockSlice 的内容
func (s *RockSlice) Print() {
	fmt.Printf("Len: %d, Capacity: %d\n", s.Size, s.Capacity)

	fmt.Printf("data:\n")
	for i := 0; i < s.Size; i++ {
		a, err := s.Get(i)
		if err != nil {
			fmt.Printf("get  No.%v value error: %s", i, err)
		}
		fmt.Printf("i: %v, value: %v", i, a)
	}

}
