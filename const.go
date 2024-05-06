package tcmallocgo

type CalcStage int8

const (
	StorageCalc CalcStage = iota + 1
	ShuffleCalc
	IntersectionCalc
)

type MemoryMode int8

const (
	MemoryMode_OnHeap MemoryMode = iota
	MemoryMode_OffHeap
)
