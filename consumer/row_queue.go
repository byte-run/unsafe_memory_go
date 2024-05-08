package consumer

import (
	"container/list"
	tcmallocgo "github.com/byte-run/unsafe_mem_go"
	"github.com/byte-run/unsafe_mem_go/memory"
)

// RowQueue , similar to a collection of a lot of row on raw table.
type RowQueue struct {
	MemoryConsumer
	//TaskMemoryManager *tcmallocgo.TaskMemoryManager
	allocatePage list.List

	CurrentPage         memory.MemBlock
	PageCursor          uintptr
	PeakMemoryUsedBytes uintptr // 使用时的峰值内存

	Stage tcmallocgo.CalcStage
}

func (row *RowQueue) appendRecord() {}

func (row *RowQueue) GetStage() tcmallocgo.CalcStage {
	return row.Stage
}

func NewRowQueue(taskMemoryManager tcmallocgo.TaskMemoryManager, pageSizeBytes uintptr) *RowQueue {
	var rowQueue = new(RowQueue)

	//rowQueue.PageSize = uint64(pageSizeBytes)

	return rowQueue
}
