package bitset

import (
	"errors"
	"math"
)

const addressBitsPerWord uint = 6

// wordMask used to shift left or right for partial word mask
const wordMask uint64 = 0xffffffffffffffff

// BitSet 位图
type BitSet struct {
	wordInUse uint
	words     []uint64
}

func (bit *BitSet) initWords(size uint) {
	arrayLength := WordIndex(size-1) + 1
	bit.words = make([]uint64, arrayLength)
}

func FromWithLength(length uint) (*BitSet, error) {
	if length < 0 {
		return nil, errors.New("length must be greater than or equal to 0")
	}

	bitset := new(BitSet)
	bitset.initWords(length)
	bitset.wordInUse = 0
	return bitset, nil
}

func WordIndex(index uint) uint {
	return index >> addressBitsPerWord
}

func (bit *BitSet) checkInvariants() error {
	if bit.wordInUse == 0 || bit.words[bit.wordInUse-1] != 0 {
		return nil
	}
	if bit.wordInUse >= 0 && bit.wordInUse <= uint(cap(bit.words)) {
		return nil
	}
	if bit.wordInUse == uint(cap(bit.words)) || bit.words[bit.wordInUse] == 0 {
		return nil
	}
	return errors.New("the number of word in use error")
}

func (bit *BitSet) ensureCapacity(requiredWords uint) {
	setCurrentCapacity := uint(cap(bit.words))
	if setCurrentCapacity < requiredWords {
		// Allocate a larger of double capacity or required size
		newCapacity := uint(math.Max(float64(2*setCurrentCapacity), float64(requiredWords)))
		newWords := make([]uint64, newCapacity)
		copy(newWords, bit.words)
		bit.words = newWords
	}
}

func (bit *BitSet) expendTo(wordIndex uint) {
	var requiredWords = wordIndex + 1
	if requiredWords > bit.wordInUse {
		bit.ensureCapacity(requiredWords)
		bit.wordInUse = requiredWords
	}
}

func (bit *BitSet) Set(bitIndex uint) error {
	if bitIndex >= bit.wordInUse {
		return errors.New("index out of range")
	}
	wordIndex := WordIndex(bitIndex)
	bit.expendTo(wordIndex)

	bit.words[wordIndex] |= (1 << bitIndex)
	return bit.checkInvariants()
}

func (bit *BitSet) SetWithValue(word uint64, value bool) {}

func (bit *BitSet) Clear(index uint) {
}
