package memory

const (
	PageWordSize   = 8
	LargePageShift = 17
	LargePageSize  = 1 << LargePageShift // 125K
	PageSizeShift  = 1 << PageWordSize   // 256byte
)

type memLocation struct {
	obj    any
	offset uintptr // address
}

func (loc *memLocation) ClearObjAndOffset() {
	loc.obj = nil
	loc.offset = uintptr(0)
}

func (loc *memLocation) GetObj() any {
	return loc.obj
}

func (loc *memLocation) GetOffset() uintptr {
	return loc.offset
}

const (
	NoPageNumber = iota - 3
	FreedInTMMPageNumber
	FreedInAllocatorPageNumber
)

// MemBlock MemBlock的大小[0, MaxPageSize]
type MemBlock struct {
	*memLocation
	length     uintptr // request size
	PageNumber uintptr // allocated pageNumber #{pageTable}

}

func (block *MemBlock) Size() uintptr {
	return block.length
}

func NewMemBlock(obj any, offset uintptr, length uintptr) *MemBlock {
	block := new(MemBlock)
	block.obj = obj
	block.offset = offset
	block.length = length
	block.PageNumber = NoPageNumber

	return block
}
