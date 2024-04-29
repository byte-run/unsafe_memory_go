package tcmallocgo

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func Test_memory_string(t *testing.T) {
	slice, _ := NewRockSlice(STRING)
	// 添加一些元素
	for i := 0; i < 10; i++ {
		a := "hellohellohellohellohellohellohellohellohellohellohellohello" + strconv.Itoa(i)
		fmt.Println(a)
		err := slice.Append(a)
		if err != nil {
			t.Logf("error type, target type: %v, value type: %v", reflect.String, reflect.TypeOf(a))
		}
	}
	// 打印结果
	slice.Print()
}

func Test_memory_int(t *testing.T) {
	slice, _ := NewRockSlice(INT)

	// 添加一些元素
	for i := 0; i < 5; i++ {
		a := i
		fmt.Println(a)
		err := slice.Append(a)
		if err != nil {
			t.Logf("error type, target type: %v, value type: %v", reflect.Int, reflect.TypeOf(a))
		}
	}
	// 打印结果
	slice.Print()
}

func Test_memory_int32(t *testing.T) {
	slice, _ := NewRockSlice(INT32)

	// 添加一些元素
	for i := 0; i < 5; i++ {
		a := int32(i)
		fmt.Println(a)
		err := slice.Append(a)
		if err != nil {
			t.Logf("error type, target type: %v, value type: %v", reflect.Int32, reflect.TypeOf(a))
		}
	}
	// 打印结果
	slice.Print()
}

func Test_memory_float64(t *testing.T) {
	slice, _ := NewRockSlice(FLOAT64)

	// 添加一些元素
	for i := 0; i < 5; i++ {
		a := float64(i)
		fmt.Println(a)
		err := slice.Append(a)
		if err != nil {
			t.Logf("error type, target type: %v, value type: %v", reflect.Float64, reflect.TypeOf(a))
		}
	}
	// 打印结果
	slice.Print()
}

func Test_memory_bool(t *testing.T) {
	slice, _ := NewRockSlice(BOOL)

	// 添加一些元素
	for i := 0; i < 5; i++ {
		a := i%2 == 1
		fmt.Println(a)
		err := slice.Append(a)
		if err != nil {
			t.Logf("error type, target type: %v, value type: %v", reflect.Bool, reflect.TypeOf(a))
		}
	}
	// 打印结果
	slice.Print()
}

func Test_memory_time(t *testing.T) {
	slice, _ := NewRockSlice(TIME)

	// 添加一些元素
	for i := 0; i < 5; i++ {
		a := time.Now()
		time.Sleep(1 * time.Second)
		fmt.Println(a)
		err := slice.Append(a)
		if err != nil {
			t.Logf("error type, target type: %v, value type: %v", reflect.TypeOf(time.Time{}), reflect.TypeOf(a))
		}
	}
	// 打印结果
	slice.Print()
}

func Test_interface(t *testing.T) {
	slice, _ := NewRockSlice(INTERFACE)

	// 添加一些元素
	for i := 0; i < 10; i++ {
		var a interface{} = "hellohellohellohellohellohellohellohellohellohellohellohello" + strconv.Itoa(i)
		fmt.Println(a)
		err := slice.Append(a)
		if err != nil {
			t.Logf("error type, target type: %v, value type: %v", reflect.Interface, reflect.TypeOf(a))
		}
	}
	// 打印结果
	slice.Print()
}

func Test_memory_mistype(t *testing.T) {
	slice, _ := NewRockSlice(INT)

	// 添加一些元素
	for i := 0; i < 5; i++ {
		a := i%2 == 1
		fmt.Println(a)
		err := slice.Append(a)
		if err != nil {
			t.Logf("error type, target type: %v, value type: %v", reflect.Int, reflect.TypeOf(a))
		}
	}
	// 打印结果
	slice.Print()
}

func Test_intersection(t *testing.T) {
	var intersection [][]*RockSlice

	leftPli := make(map[int32][]int32)
	leftPli[0] = []int32{0, 1, 2}
	leftPli[1] = []int32{3, 5}
	leftPli[2] = []int32{4, 6, 7}
	rightPli := make(map[int32][]int32)
	rightPli[0] = []int32{3, 5, 6}
	rightPli[1] = []int32{1, 2}
	rightPli[2] = []int32{0, 4}
	leftTableIndex := 0
	rightTableIndex := 1
	maxTableIndex := rightTableIndex

	isSameColumn := false
	for value, leftIds := range leftPli {
		idPairs := make([]*RockSlice, maxTableIndex+1)
		if isSameColumn {
			//idPairs[leftTableIndex] = leftIds
			//idPairs[rightTableIndex] = leftIds
			leftIdList, err := NewRockSlice(INT32)
			if err != nil {
				return
			}
			for _, id := range leftIds {
				leftIdList.Append(id)
			}
			idPairs[leftTableIndex] = leftIdList
			idPairs[rightTableIndex] = leftIdList
		} else {
			rightIds := rightPli[value]
			if len(rightIds) > 0 {
				leftIdList, err := NewRockSlice(INT32)
				if err != nil {
					return
				}
				for _, id := range leftIds {
					leftIdList.Append(id)
				}
				idPairs[leftTableIndex] = leftIdList

				rightIdList, err := NewRockSlice(INT32)
				if err != nil {
					return
				}
				for _, id := range rightIds {
					rightIdList.Append(id)
				}
				idPairs[rightTableIndex] = rightIdList
			} else {
				continue
			}
		}
		intersection = append(intersection, idPairs)
	}

	for index, idPairs := range intersection {
		t.Logf("--------------------")
		t.Logf("value group:%v", index)

		for i, ids := range idPairs {
			t.Logf("==table i:%v==", i)
			ids.Print()
		}
		t.Logf("--------------------")
	}
}
