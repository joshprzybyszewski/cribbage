package rand

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntn(t *testing.T) {
	max := 1000
	numRuns := 100 * max
	results := make(map[int]int, max)
	for i := 0; i < numRuns; i++ {
		results[Intn(max)]++
	}
	for _, num := range results {
		assert.Less(t, num, (numRuns/max)*2)
	}
}

func TestInt64n(t *testing.T) {
	max := int64(1000)
	numRuns := 100 * max
	results := make(map[int64]int, max)
	for i := int64(0); i < numRuns; i++ {
		results[Int64n(max)]++
	}
	for _, num := range results {
		assert.Less(t, num, int((numRuns/max)*2))
	}
}

func TestInt(t *testing.T) {
	for i := 0; i < 100; i++ {
		r := Int()
		assert.Less(t, r, 10000)
		assert.Greater(t, r, -1)
	}
}

func TestFloat64(t *testing.T) {
	for i := 0; i < 100; i++ {
		r := Float64()
		assert.Less(t, r, 1.00001)
		assert.Greater(t, r, -0.00001)
	}
}
