package tcmallocgo

import (
	"strings"
	"testing"
)

func TestMemStringSplit(t *testing.T) {
	testMem := "10GB2"
	splitArray := strings.Split(testMem, "GB")
	t.Logf("splitArray: %v", splitArray)
	t.Logf("splitArray length: %v", len(splitArray))
}
