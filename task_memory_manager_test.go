package tcmallocgo

import (
	"github.com/byte-run/unsafe_mem_go/utils"
	"testing"
	"time"
)

func TestMemoryManager_Init(t *testing.T) {
	memManager := InitMemoryManager(&MemoryConfig{ShuffleMem: "5G", IntersectionMem: "5G", StorageMem: "5G"})
	t.Logf("memory manager is %v", memManager)
}

// ---------------------------- That part for storage memory test ----------------------------
func TestMemoryManager_AllocateStoragePool(t *testing.T) {
	memManager := InitMemoryManager(&MemoryConfig{ShuffleMem: "5G", IntersectionMem: "5G", StorageMem: "5G"})

	var storageNumBytes = uintptr(3 * 1024 * 1024) // 3MB
	isSuccess, warnInfo, err := memManager.AcquireStorageMemory(storageNumBytes)
	if err != nil {
		t.Logf("AcquireStorageMemory err %v", err)
		return
	}
	if warnInfo != nil {
		//warnInfo.type()
		t.Logf("AcquireStorageMemory warnInfo %v", warnInfo.Warning())
		return
	}
	t.Logf("AcquireStorageMemory isSuccess %v", isSuccess)
}

func TestMemoryManager_AllocateStoragePool_UsageWarn(t *testing.T) {
	memManager := InitMemoryManager(&MemoryConfig{ShuffleMem: "5G", IntersectionMem: "5G", StorageMem: "5G"})

	var storageNumBytes = uintptr(4*GB_Factor + 3*MB_Factor)
	isSuccess, warnInfo, err := memManager.AcquireStorageMemory(storageNumBytes)
	if err != nil {
		t.Logf("AcquireStorageMemory err %v", err)
		return
	}
	if warnInfo != nil {
		//warnInfo.type()
		switch warnInfo.(type) {
		case utils.MemoryPoolLevelOneWarning:
			t.Logf("storageMemory level one warnInfo %s", warnInfo.Warning())
		case utils.MemoryPoolLevelTwoWarning:
			t.Logf("storageMemory level two warn: %s", warnInfo.Warning())
		default:
			t.Log("not match warnType")
		}
		return
	}
	t.Logf("AcquireStorageMemory isSuccess %v", isSuccess)
}

func TestMemoryManager_AllocateStoragePool_Release(t *testing.T) {
	// 申请
	memManager := InitMemoryManager(&MemoryConfig{ShuffleMem: "5G", IntersectionMem: "5G", StorageMem: "5G"})

	var storageNumBytes = uintptr(4*GB_Factor + 3*MB_Factor)
	isSuccess, warnInfo, err := memManager.AcquireStorageMemory(storageNumBytes)
	if err != nil {
		t.Logf("AcquireStorageMemory err %v", err)
		return
	}
	if warnInfo != nil {
		//warnInfo.type()
		switch warnInfo.(type) {
		case utils.MemoryPoolLevelOneWarning:
			t.Logf("storageMemory level one warnInfo %s", warnInfo.Warning())
		case utils.MemoryPoolLevelTwoWarning:
			t.Logf("storageMemory level two warn: %s", warnInfo.Warning())
		default:
			t.Log("not match warnType")
		}
	}
	t.Logf("AcquireStorageMemory isSuccess %v", isSuccess)

	// 释放
	err = memManager.ReleaseStorageMemory(4*GB_Factor + 3*MB_Factor)
	if err != nil {
		t.Logf("ReleaseStorageMemory err %v", err)
	}
}

func TestMemoryManager_AllocateStoragePool_AndMemory(t *testing.T) {
	memManager := InitMemoryManager(&MemoryConfig{ShuffleMem: "5G", IntersectionMem: "5G", StorageMem: "5G"})
	// 申请
	// 先申请memPool
	var storageNumBytes = uintptr(3 * GB_Factor)
	isSuccess, warnInfo, err := memManager.AcquireStorageMemory(storageNumBytes)
	if err != nil {
		t.Logf("AcquireStorageMemory err %v", err)
		return
	}
	if warnInfo != nil {
		//warnInfo.type()
		switch warnInfo.(type) {
		case utils.MemoryPoolLevelOneWarning:
			t.Logf("storageMemory level one warnInfo %s", warnInfo.Warning())
		case utils.MemoryPoolLevelTwoWarning:
			t.Logf("storageMemory level two warn: %s", warnInfo.Warning())
		default:
			t.Log("not match warnType")
		}
	}
	t.Logf("AcquireStorageMemory isSuccess %v", isSuccess)
	// 再申请page
	allocatePage, err := memManager.AllocatePage(storageNumBytes)
	if err != nil {
		t.Logf("Allocate err %v", err)
		return
	}

	time.Sleep(30 * time.Second) // 假设休眠30s

	// 释放
	// 先释放page - 再释放memPool
	memManager.FreePage(uintptr(allocatePage), storageNumBytes)
	err = memManager.ReleaseStorageMemory(4*GB_Factor + 3*MB_Factor)
	if err != nil {
		t.Logf("ReleaseStorageMemory err %v", err)
	}
}

// ---------------------------- That part for shuffle memory test ----------------------------
func TestMemoryManager_AllocateShuffleMemPool(t *testing.T) {
	memManager := InitMemoryManager(&MemoryConfig{ShuffleMem: "5G", IntersectionMem: "5G", StorageMem: "5G"})

	// 申请size from shufflePool
	shuffleNumBytes := uintptr(3 * GB_Factor)
	acquireBytes, poolStatus, err := memManager.AcquireShuffleMemory(shuffleNumBytes)
	if err != nil {
		t.Logf("AcquireShuffleMemory err %v", err)
		return
	}
	if poolStatus != nil {
		switch poolStatus.(type) {
		case utils.MemoryPoolLevelOneWarning:
			t.Logf("shuffleMemory level one warnInfo: %s", poolStatus.Warning())
		case utils.MemoryPoolLevelTwoWarning:
			t.Logf("shuffleMemory level two warnInfo: %s", poolStatus.Warning())
		default:
			t.Log("not match warnType")
		}
		return
	}
	t.Logf("AcquireShuffleMemory isSuccess %v", acquireBytes)
}

func TestMemoryManager_AllocateShuffleMemPool_UsageWarn(t *testing.T) {
	memManager := InitMemoryManager(&MemoryConfig{ShuffleMem: "5G", IntersectionMem: "5G", StorageMem: "5G"})

	// 申请size from shufflePool
	shuffleNumBytes := uintptr(4*GB_Factor + 4*MB_Factor)
	acquireBytes, poolStatus, err := memManager.AcquireShuffleMemory(shuffleNumBytes)
	if err != nil {
		t.Logf("AcquireShuffleMemory err %v", err)
		return
	}
	if poolStatus != nil {
		switch poolStatus.(type) {
		case utils.MemoryPoolLevelOneWarning:
			t.Logf("shuffleMemory level one warnInfo: %s", poolStatus.Warning())
		case utils.MemoryPoolLevelTwoWarning:
			t.Logf("shuffleMemory level two warnInfo: %s", poolStatus.Warning())
		default:
			t.Log("not match warnType")
		}
		return
	}
	t.Logf("AcquireShuffleMemory isSuccess %v", acquireBytes)
}

func TestMemoryManager_AllocateShuffleMemPool_AndRelease(t *testing.T) {
	memManager := InitMemoryManager(&MemoryConfig{ShuffleMem: "5G", IntersectionMem: "5G", StorageMem: "5G"})

	// 申请size from shufflePool
	shuffleNumBytes := uintptr(3 * GB_Factor)
	acquireBytes, poolStatus, err := memManager.AcquireShuffleMemory(shuffleNumBytes)
	if err != nil {
		t.Logf("AcquireShuffleMemory err %v", err)
		return
	}
	if poolStatus != nil {
		switch poolStatus.(type) {
		case utils.MemoryPoolLevelOneWarning:
			t.Logf("shuffleMemory level one warnInfo: %s", poolStatus.Warning())
		case utils.MemoryPoolLevelTwoWarning:
			t.Logf("shuffleMemory level two warnInfo: %s", poolStatus.Warning())
		default:
			t.Log("not match warnType")
		}
		return
	}
	t.Logf("AcquireShuffleMemory isSuccess %v", acquireBytes)

	// 释放
	err = memManager.ReleaseShuffleMemory(shuffleNumBytes)
	if err != nil {
		t.Logf("ReleaseShuffleMemory err %v", err)
		return
	}
	t.Log("ReleaseShuffleMemory")
}

func TestMemoryManager_AllocateShuffleMemPool_AndRelease_TwoAllocate(t *testing.T) {
	memManager := InitMemoryManager(&MemoryConfig{ShuffleMem: "5G", IntersectionMem: "5G", StorageMem: "5G"})

	// 申请size from shufflePool
	shuffleNumBytes := uintptr(1 * GB_Factor)
	acquireBytes, poolStatus, err := memManager.AcquireShuffleMemory(shuffleNumBytes)
	if err != nil {
		t.Logf("AcquireShuffleMemory err %v", err)
		return
	}

	// other allocate request
	go func() {
		shuffleNumBytes := uintptr(2 * GB_Factor)
		acquireBytes, poolStatus, err := memManager.AcquireShuffleMemory(shuffleNumBytes)
		if err != nil {
			t.Logf("AcquireShuffleMemory2 err %v", err)
			return
		}
		if poolStatus != nil {
			switch poolStatus.(type) {
			case utils.MemoryPoolLevelOneWarning:
				t.Logf("shuffleMemory2 level one warnInfo: %s", poolStatus.Warning())
			case utils.MemoryPoolLevelTwoWarning:
				t.Logf("shuffleMemory2 level two warnInfo: %s", poolStatus.Warning())
			default:
				t.Log("2not match warnType")
			}
			return
		}
		t.Logf("AcquireShuffleMemory isSuccess %v", acquireBytes)

		// 休眠
		time.Sleep(30 * time.Second)

		// 释放
		err = memManager.ReleaseShuffleMemory(shuffleNumBytes)
		if err != nil {
			t.Logf("ReleaseShuffleMemory2 err %v", err)
			return
		}
		t.Log("ReleaseShuffleMemory2")
	}()

	if poolStatus != nil {
		switch poolStatus.(type) {
		case utils.MemoryPoolLevelOneWarning:
			t.Logf("shuffleMemory level one warnInfo: %s", poolStatus.Warning())
		case utils.MemoryPoolLevelTwoWarning:
			t.Logf("shuffleMemory level two warnInfo: %s", poolStatus.Warning())
		default:
			t.Log("not match warnType")
		}
		return
	}
	t.Logf("AcquireShuffleMemory isSuccess %v", acquireBytes)

	// 休眠
	time.Sleep(20 * time.Second)

	// 释放
	err = memManager.ReleaseShuffleMemory(shuffleNumBytes)
	if err != nil {
		t.Logf("ReleaseShuffleMemory err %v", err)
		return
	}
	t.Log("ReleaseShuffleMemory")
	time.Sleep(20 * time.Second)
}

func TestMemoryManager_AllocateShuffleMemPool_Memory(t *testing.T) {
	memManager := InitMemoryManager(&MemoryConfig{ShuffleMem: "5G", IntersectionMem: "5G", StorageMem: "5G"})

	// 申请size from shufflePool
	shuffleNumBytes := uintptr(3 * GB_Factor)
	acquireBytes, poolStatus, err := memManager.AcquireShuffleMemory(shuffleNumBytes)
	if err != nil {
		t.Logf("AcquireShuffleMemory err %v", err)
		return
	}
	if poolStatus != nil {
		switch poolStatus.(type) {
		case utils.MemoryPoolLevelOneWarning:
			t.Logf("shuffleMemory level one warnInfo: %s", poolStatus.Warning())
		case utils.MemoryPoolLevelTwoWarning:
			t.Logf("shuffleMemory level two warnInfo: %s", poolStatus.Warning())
		default:
			t.Log("not match warnType")
		}
		return
	}
	t.Logf("AcquireShuffleMemory isSuccess %v", acquireBytes)
	// 再申请内存空间
	addr, err := memManager.memAllocator.Allocate(acquireBytes)
	if err != nil {
		t.Logf("Allocate memory from system err %v", err)
		return
	}

	time.Sleep(30 * time.Second)

	// 释放
	// 释放内存
	memManager.FreePage(uintptr(addr), acquireBytes)
	// 再释放内存池
	err = memManager.ReleaseShuffleMemory(shuffleNumBytes)
	if err != nil {
		t.Logf("ReleaseShuffleMemory err %v", err)
		return
	}
	t.Log("ReleaseShuffleMemory")
}

// ---------------------------- That part for shuffle memory test ----------------------------
func TestMemoryManager_AllocateIntersectionMemPool(t *testing.T) {
	memManager := InitMemoryManager(&MemoryConfig{ShuffleMem: "5G", IntersectionMem: "5G", StorageMem: "5G"})

	// 申请
	intersectionNumBytes := uintptr(3 * GB_Factor)
	acquired, poolWarn, err := memManager.AcquireIntersectionMemory(intersectionNumBytes)
	if err != nil {
		t.Logf("AcquireIntersectionMemory err %v", err)
	}
	if poolWarn != nil {
		switch poolWarn.(type) {
		case utils.MemoryPoolLevelOneWarning:
			t.Logf("intersectionMemory level one warnInfo: %s", poolWarn.Warning())
		case utils.MemoryPoolLevelTwoWarning:
			t.Logf("intersectionMemory level two warnInfo: %s", poolWarn.Warning())
		default:
			t.Logf("not match warnType")
		}
	}
	t.Logf("AcquireIntersectionMemory isSuccess %v", acquired)

	// 释放pool

	err = memManager.ReleaseIntersectionMemory(acquired)
	if err != nil {
		t.Logf("ReleaseIntersectionMemory err %v", err)
	}

	t.Log("ReleaseIntersectionMemory")
}
