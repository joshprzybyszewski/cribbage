package suggestions

import (
	"testing"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/network"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ model.TossStats = (*testingTossStats)(nil)

type testingTossStats struct {
	min    int
	avg    float64
	median float64
	max    int
}

func (ts *testingTossStats) Min() int {
	return ts.min
}

func (ts *testingTossStats) Median() float64 {
	return ts.median
}

func (ts *testingTossStats) Avg() float64 {
	return ts.avg
}

func (ts *testingTossStats) Max() int {
	return ts.max
}

func TestGetStatsForHand(t *testing.T) {
	tests := []struct {
		desc         string
		hand         []string
		tossed       []string
		expHandStats *testingTossStats
		expCribStats *testingTossStats
	}{{
		desc:   `low pointer`,
		hand:   []string{`AH`, `QH`, `JC`, `9D`},
		tossed: []string{`10S`, `KS`},
		expHandStats: &testingTossStats{
			avg:    2.282608695652174,
			median: 2,
			max:    7,
		},
		expCribStats: &testingTossStats{
			avg:    3.527975406236276,
			median: 2,
			max:    20,
		},
	}, {
		desc:   `max pointer`,
		hand:   []string{`5h`, `jd`, `5C`, `5s`},
		tossed: []string{`4d`, `6d`},
		expHandStats: &testingTossStats{
			min:    14,
			avg:    16.608695652173914,
			median: 14,
			max:    29,
		},
		expCribStats: &testingTossStats{
			avg:    3.9596179183135707,
			median: 4,
			max:    24,
		},
	}}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			actHand, actCrib := getStatsForHand(
				network.ModelCardsFromStrings(tc.hand...),
				network.ModelCardsFromStrings(tc.tossed...),
			)
			assert.Equal(t, tc.expHandStats.Min(), actHand.Min())
			assert.Equal(t, tc.expHandStats.Avg(), actHand.Avg())
			assert.Equal(t, tc.expHandStats.Median(), actHand.Median())
			assert.Equal(t, tc.expHandStats.Max(), actHand.Max())
			assert.Equal(t, tc.expCribStats.Min(), actCrib.Min())
			assert.Equal(t, tc.expCribStats.Avg(), actCrib.Avg())
			assert.Equal(t, tc.expCribStats.Median(), actCrib.Median())
			assert.Equal(t, tc.expCribStats.Max(), actCrib.Max())
		})
	}
}

func TestGetAllTosses(t *testing.T) {
	tests := []struct {
		desc         string
		hand         []string
		expSummaries []model.TossSummary
	}{{
		desc: `max pointer`,
		hand: []string{`5h`, `jd`, `5C`, `5s`, `6d`},
		expSummaries: []model.TossSummary{{
			Kept:   network.ModelCardsFromStrings(`5h`, `jd`, `5C`, `5s`),
			Tossed: network.ModelCardsFromStrings(`6d`),
			HandStats: &testingTossStats{
				min:    14,
				avg:    16.574468085106382,
				median: 14,
				max:    29,
			},
			CribStats: &testingTossStats{
				avg:    4.228551004961735,
				median: 4,
				max:    24,
			},
		}, {
			Kept:   network.ModelCardsFromStrings(`5h`, `jd`, `5C`, `6d`),
			Tossed: network.ModelCardsFromStrings(`5s`),
			HandStats: &testingTossStats{
				min:    6,
				avg:    9.46808510638298,
				median: 10,
				max:    17,
			},
			CribStats: &testingTossStats{
				min:    2,
				avg:    6.144672441342191,
				median: 6,
				max:    24,
			},
		}, {
			Kept:   network.ModelCardsFromStrings(`5h`, `jd`, `5s`, `6d`),
			Tossed: network.ModelCardsFromStrings(`5c`),
			HandStats: &testingTossStats{
				min:    6,
				avg:    9.46808510638298,
				median: 10,
				max:    17,
			},
			CribStats: &testingTossStats{
				min:    2,
				avg:    6.144672441342191,
				median: 6,
				max:    24,
			},
		}, {
			Kept:   network.ModelCardsFromStrings(`5h`, `5C`, `5s`, `6d`),
			Tossed: network.ModelCardsFromStrings(`jd`),
			HandStats: &testingTossStats{
				min:    8,
				avg:    12.51063829787234,
				median: 14,
				max:    23,
			},
			CribStats: &testingTossStats{
				avg:    4.066820844896701,
				median: 4,
				max:    21,
			},
		}, {
			Kept:   network.ModelCardsFromStrings(`jd`, `5C`, `5s`, `6d`),
			Tossed: network.ModelCardsFromStrings(`5h`),
			HandStats: &testingTossStats{
				min:    6,
				avg:    9.46808510638298,
				median: 10,
				max:    17,
			},
			CribStats: &testingTossStats{
				min:    2,
				avg:    6.144672441342191,
				median: 6,
				max:    24,
			},
		}},
	}}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			actSums, err := GetAllTosses(
				network.ModelCardsFromStrings(tc.hand...),
			)
			require.NoError(t, err)
			for i, := range actSums {
				actSum := actSums[i]
				expSum := tc.expSummaries[i]
				assert.Equal(t, expSum.Kept, actSum.Kept)
				assert.Equal(t, expSum.Tossed, actSum.Tossed)
				assert.Equal(t, expSum.HandStats.Min(), actSum.HandStats.Min())
				assert.Equal(t, expSum.HandStats.Avg(), actSum.HandStats.Avg())
				assert.Equal(t, expSum.HandStats.Median(), actSum.HandStats.Median())
				assert.Equal(t, expSum.HandStats.Max(), actSum.HandStats.Max())
				assert.Equal(t, expSum.CribStats.Min(), actSum.CribStats.Min())
				assert.Equal(t, expSum.CribStats.Avg(), actSum.CribStats.Avg())
				assert.Equal(t, expSum.CribStats.Median(), actSum.CribStats.Median())
				assert.Equal(t, expSum.CribStats.Max(), actSum.CribStats.Max())
			}

		})
	}
}
