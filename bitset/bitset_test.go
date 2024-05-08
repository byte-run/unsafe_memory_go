package bitset

import "testing"

func TestBitsetInit(t *testing.T) {
	bitset, err := FromWithLength(256)
	if err != nil {
		t.Logf("bitset error: %v", err)
		return
	}

	for i := 0; i < len(bitset.words); i++ {
		t.Log(bitset.words[i])
	}
}
