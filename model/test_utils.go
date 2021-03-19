// +build !prod

package model

var _ TossStats = (*TestingTossStats)(nil)

type TestingTossStats struct {
	min    int
	avg    float64
	median float64
	max    int
}

func NewTestingTossStats(min int, avg, median float64, max int) *TestingTossStats {
	return &TestingTossStats{
		min:    min,
		avg:    avg,
		median: median,
		max:    max,
	}
}

func (ts *TestingTossStats) Min() int {
	return ts.min
}

func (ts *TestingTossStats) Median() float64 {
	return ts.median
}

func (ts *TestingTossStats) Avg() float64 {
	return ts.avg
}

func (ts *TestingTossStats) Max() int {
	return ts.max
}
