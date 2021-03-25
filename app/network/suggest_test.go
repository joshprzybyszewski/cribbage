package network

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/joshprzybyszewski/cribbage/model"
)

func TestConvertToGetSuggestHandResponse(t *testing.T) {
	testCases := []struct {
		input     []model.TossSummary
		expOutput []GetSuggestHandResponse
	}{{
		input:     nil,
		expOutput: nil,
	}, {
		input: []model.TossSummary{{
			Kept:   ModelCardsFromStrings(`AH`, `AD`, `AS`, `AC`),
			Tossed: ModelCardsFromStrings(`2H`, `2D`),
			HandStats: model.NewTestingTossStats(
				0, 1, 2, 3,
			),
			CribStats: model.NewTestingTossStats(
				4, 5, 6, 7,
			),
		}},
		expOutput: []GetSuggestHandResponse{{
			Hand: []string{`AH`, `AD`, `AS`, `AC`},
			Toss: []string{`2H`, `2D`},
			HandPts: PointStats{
				Min:    0,
				Avg:    1,
				Median: 2,
				Max:    3,
			},
			CribPts: PointStats{
				Min:    4,
				Avg:    5,
				Median: 6,
				Max:    7,
			},
		}},
	}}

	for _, tc := range testCases {
		act := ConvertToGetSuggestHandResponse(tc.input)
		assert.Equal(t, tc.expOutput, act)
	}
}
