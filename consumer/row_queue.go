package consumer

import (
	"container/list"
	tcmallocgo "github.com/byte-run/unsafe_mem_go"
	"github.com/byte-run/unsafe_mem_go/memory"
)

// RowQueue , similar to a collection of a lot of row on raw table.
type RowQueue struct {
	memoryConsumer
	//TaskMemoryManager *tcmallocgo.TaskMemoryManager
	allocatePage list.List

	CurrentPage         memory.MemBlock
	PageCursor          uintptr
	PeakMemoryUsedBytes uintptr // 使用时的峰值内存
}

func (row *RowQueue) appendRecord() {}

func NewRowQueue(taskMemoryManager tcmallocgo.TaskMemoryManager, pageSizeBytes uintptr) *RowQueue {
	var rowQueue = new(RowQueue)

	rowQueue.pageSize = uint64(pageSizeBytes)

	return rowQueue
}
