package suggestions

import (
	"testing"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/network"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetStatsForHand(t *testing.T) {
	tests := []struct {
		desc         string
		hand         []string
		tossed       []string
		expHandStats *model.TestingTossStats
		expCribStats *model.TestingTossStats
	}{{
		desc:   `low pointer`,
		hand:   []string{`AH`, `QH`, `JC`, `9D`},
		tossed: []string{`10S`, `KS`},
		expHandStats: model.NewTestingTossStats(
			0,
			2.282608695652174,
			2,
			7,
		),
		expCribStats: model.NewTestingTossStats(
			0,
			3.527975406236276,
			2,
			20,
		),
	}, {
		desc:   `max pointer`,
		hand:   []string{`5h`, `jd`, `5C`, `5s`},
		tossed: []string{`4d`, `6d`},
		expHandStats: model.NewTestingTossStats(
			14,
			16.608695652173914,
			14,
			29,
		),
		expCribStats: model.NewTestingTossStats(
			0,
			3.9596179183135707,
			4,
			24,
		),
	}}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			actHand, actCrib := getStatsForHand(
				network.ModelCardsFromStrings(tc.hand...),
				network.ModelCardsFromStrings(tc.tossed...),
			)
			assert.Equal(t, tc.expHandStats.Min(), actHand.Min())
			assert.InEpsilon(t, tc.expHandStats.Avg(), actHand.Avg(), 0.001)
			assert.InEpsilon(t, tc.expHandStats.Median(), actHand.Median(), 0.001)
			assert.Equal(t, tc.expHandStats.Max(), actHand.Max())
			assert.Equal(t, tc.expCribStats.Min(), actCrib.Min())
			assert.InEpsilon(t, tc.expCribStats.Avg(), actCrib.Avg(), 0.001)
			assert.InEpsilon(t, tc.expCribStats.Median(), actCrib.Median(), 0.001)
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
			HandStats: model.NewTestingTossStats(
				14,
				16.574468085106382,
				14,
				29,
			),
			CribStats: model.NewTestingTossStats(
				0,
				4.228551004961735,
				4,
				24,
			),
		}, {
			Kept:   network.ModelCardsFromStrings(`5h`, `jd`, `5C`, `6d`),
			Tossed: network.ModelCardsFromStrings(`5s`),
			HandStats: model.NewTestingTossStats(
				6,
				9.46808510638298,
				10,
				17,
			),
			CribStats: model.NewTestingTossStats(
				2,
				6.144672441342191,
				6,
				24,
			),
		}, {
			Kept:   network.ModelCardsFromStrings(`5h`, `jd`, `5s`, `6d`),
			Tossed: network.ModelCardsFromStrings(`5c`),
			HandStats: model.NewTestingTossStats(
				6,
				9.46808510638298,
				10,
				17,
			),
			CribStats: model.NewTestingTossStats(
				2,
				6.144672441342191,
				6,
				24,
			),
		}, {
			Kept:   network.ModelCardsFromStrings(`5h`, `5C`, `5s`, `6d`),
			Tossed: network.ModelCardsFromStrings(`jd`),
			HandStats: model.NewTestingTossStats(
				8,
				12.51063829787234,
				14,
				23,
			),
			CribStats: model.NewTestingTossStats(
				0,
				4.066820844896701,
				4,
				21,
			),
		}, {
			Kept:   network.ModelCardsFromStrings(`jd`, `5C`, `5s`, `6d`),
			Tossed: network.ModelCardsFromStrings(`5h`),
			HandStats: model.NewTestingTossStats(
				6,
				9.46808510638298,
				10,
				17,
			),
			CribStats: model.NewTestingTossStats(
				2,
				6.144672441342191,
				6,
				24,
			),
		}},
	}}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			actSumms, err := GetAllTosses(
				network.ModelCardsFromStrings(tc.hand...),
			)
			require.NoError(t, err)
			for i := range actSumms {
				actSum := actSumms[i]
				expSum := tc.expSummaries[i]
				assert.Equal(t, expSum.Kept, actSum.Kept)
				assert.Equal(t, expSum.Tossed, actSum.Tossed)
				assert.Equal(t, expSum.HandStats.Min(), actSum.HandStats.Min())
				assert.InEpsilon(t, expSum.HandStats.Avg(), actSum.HandStats.Avg(), 0.001)
				assert.InEpsilon(t, expSum.HandStats.Median(), actSum.HandStats.Median(), 0.001)
				assert.Equal(t, expSum.HandStats.Max(), actSum.HandStats.Max())
				assert.Equal(t, expSum.CribStats.Min(), actSum.CribStats.Min())
				assert.InEpsilon(t, expSum.CribStats.Avg(), actSum.CribStats.Avg(), 0.001)
				assert.InEpsilon(t, expSum.CribStats.Median(), actSum.CribStats.Median(), 0.001)
				assert.Equal(t, expSum.CribStats.Max(), actSum.CribStats.Max())
			}

		})
	}
}
