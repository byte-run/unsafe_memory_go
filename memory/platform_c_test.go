package memory

import (
	"testing"
)

func TestCMemAllocator_Allocate(t *testing.T) {
	addr := platformCInstance.allocate(4 * 1024 * 1024 * 1024)

	var intValue = (*[]byte)(addr)
	*(*int)(addr) = 10

	t.Logf("addr:%v, convert *int point: %v -> %d, set value: %d", addr, intValue, intValue, *intValue)
	platformCInstance.free(addr)
}
