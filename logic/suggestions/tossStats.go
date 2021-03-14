package suggestions

import (
	"sort"

	"github.com/joshprzybyszewski/cribbage/model"
)

var _ model.TossStats = (*tossStats)(nil)

type tossStats struct {
	allPts []int

	min    int
	avg    float64
	median float64
	max    int
}

func (ts *tossStats) add(pts int) {
	ts.allPts = append(ts.allPts, pts)
}

func (ts *tossStats) calculate() {
	if len(ts.allPts) == 0 {
		return
	}

	sort.Ints(ts.allPts)

	ts.min = ts.allPts[0]
	ts.max = ts.allPts[len(ts.allPts)-1]

	ts.avg = ts.getAvg()
	ts.median = ts.getMedian()
}

func (ts *tossStats) getAvg() float64 {
	sum := 0
	for _, pt := range ts.allPts {
		sum += pt
	}
	return float64(sum) / float64(len(ts.allPts))
}

func (ts *tossStats) getMedian() float64 {
	midIndex := len(ts.allPts) / 2

	if len(ts.allPts)%2 == 1 {
		return float64(ts.allPts[midIndex])
	}

	return float64(ts.allPts[midIndex-1]+ts.allPts[midIndex]) / 2
}

func (ts *tossStats) Min() int {
	return ts.min
}

func (ts *tossStats) Median() float64 {
	return ts.median
}

func (ts *tossStats) Avg() float64 {
	return ts.avg
}

func (ts *tossStats) Max() int {
	return ts.max
}
