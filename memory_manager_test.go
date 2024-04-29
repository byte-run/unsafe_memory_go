package tcmallocgo

import (
	"github.com/byte-run/unsafe_mem_go/utils"
	"testing"
)

func TestMemoryManager_Init(t *testing.T) {
	memManager := InitMemoryManager(&MemoryConfig{ShuffleMem: "5G", IntersectionMem: "5G", StorageMem: "5G"})
	t.Logf("memory manager is %v", memManager)
}

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

func TestMemoryManager_AllocateStoragePool_AndMemory(t *testing.T) {}
