package suggestions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTossStatsCalculate(t *testing.T) {
	ts := &tossStats{
		allPts: []int{
			1, 1, 1,
			4, 4, 4,
			10, 10,
		},
	}

	ts.add(10)
	assert.Len(t, ts.allPts, 9)

	// before calling calculate
	assert.Zero(t, ts.Min())
	assert.Zero(t, ts.Median())
	assert.Zero(t, ts.Avg())
	assert.Zero(t, ts.Max())

	ts.calculate()
	assert.Equal(t, 1, ts.Min())
	assert.InEpsilon(t, float64(4), ts.Median(), 0.001)
	assert.InEpsilon(t, float64(5), ts.Avg(), 0.001)
	assert.Equal(t, 10, ts.Max())
}
