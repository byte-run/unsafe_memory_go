package tcmallocgo

import (
	"math"
	"testing"
)

func TestPoolThresholdCompute(t *testing.T) {
	var poolSize uint64 = 100 * 1024 * 1024 * 1024 * 1024 * 1024 // 100TB = 112589990684262400 bytes
	var thresholdFactor float64 = 0.8

	t.Logf("poolSize convert float64: %v", float64(poolSize))         // 1.125899906842624e+17
	t.Logf("poolSize convert float64: %v", uint64(float64(poolSize))) // 112589990684262400
	thresholdValue := float64(poolSize) * thresholdFactor
	t.Logf("thresholdValue: %v", uint64(thresholdValue))
	t.Logf("thresholdValue1: %v", math.Float64bits(thresholdValue))
}
