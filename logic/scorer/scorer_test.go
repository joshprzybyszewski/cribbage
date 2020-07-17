package scorer

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshprzybyszewski/cribbage/model"
)

func TestPointsStandardFunThings(t *testing.T) {
	testCases := []struct {
		desc      string
		leadCard  string
		hand      string
		expPoints int
	}{{
		desc:      `highest scoring hand`,
		leadCard:  `5H`,
		hand:      `5S,5C,5D,JH`,
		expPoints: 29,
	}, {
		desc:      `quad`,
		leadCard:  `8H`,
		hand:      `8S,8C,8D,10H`,
		expPoints: 12,
	}, {
		desc:      `triplet`,
		leadCard:  `KH`,
		hand:      `8S,8C,8D,10H`,
		expPoints: 6,
	}, {
		desc:      `one pair`,
		leadCard:  `KH`,
		hand:      `8S,8C,2D,10H`,
		expPoints: 2,
	}, {
		desc:      `two pair`,
		leadCard:  `KH`,
		hand:      `KS,8C,2D,8H`,
		expPoints: 4,
	}, {
		desc:      `triple run of 3`,
		leadCard:  `8H`,
		hand:      `9S,8C,10D,8H`,
		expPoints: 15,
	}, {
		desc:      `double run of 4`,
		leadCard:  `JH`,
		hand:      `8S,8C,9D,10H`,
		expPoints: 10,
	}, {
		desc:      `double double run of 3`,
		leadCard:  `9H`,
		hand:      `8S,8C,9D,10H`,
		expPoints: 16,
	}, {
		desc:      `double run of 3`,
		leadCard:  `KH`,
		hand:      `8S,8C,9D,10H`,
		expPoints: 8,
	}, {
		desc:      `run of 3`,
		leadCard:  `KH`,
		hand:      `8S,2C,9D,10H`,
		expPoints: 3,
	}, {
		desc:      `run of 4`,
		leadCard:  `KH`,
		hand:      `8S,JC,9D,10H`,
		expPoints: 4,
	}, {
		desc:      `run of 5`,
		leadCard:  `QH`,
		hand:      `8S,JC,9D,10H`,
		expPoints: 5,
	}, {
		desc:      `only nobs`,
		leadCard:  `6H`,
		hand:      `JH,KC,10D,8H`,
		expPoints: 1,
	}, {
		desc:      `flush`,
		leadCard:  `3H`,
		hand:      `8D,4D,10D,6D`,
		expPoints: 6,
	}, {
		desc:      `random hand I got while playing`,
		leadCard:  `4C`,
		hand:      `5H,3D,7D,7S`,
		expPoints: 9,
	}, {
		desc:      `just looking for ways to break it`,
		leadCard:  `9C`,
		hand:      `1H,5D,7D,7S`,
		expPoints: 6,
	}}

	for _, tc := range testCases {
		lead := model.NewCardFromString(tc.leadCard)
		hand := make([]model.Card, 4)
		cardStrs := strings.Split(tc.hand, `,`)
		require.Len(t, cardStrs, 4)
		for i, c := range cardStrs {
			hand[i] = model.NewCardFromString(c)
		}

		actPoints := HandPoints(lead, hand)
		assert.Equal(t, tc.expPoints, actPoints, tc.desc)
	}
}

func TestPointsForFifteens(t *testing.T) {
	testCases := []struct {
		desc      string
		leadCard  string
		hand      string
		expPoints int
	}{{
		desc:      `highest scoring hand`,
		leadCard:  `5H`,
		hand:      `5S,5C,5D,JH`,
		expPoints: 29,
	}, {
		desc:      `15 for 2`,
		leadCard:  `8H`,
		hand:      `7S,AC,2D,KH`,
		expPoints: 2,
	}, {
		desc:      `run of 5 that adds up to 15`,
		leadCard:  `5H`,
		hand:      `AS,2C,3D,4H`,
		expPoints: 7,
	}, {
		desc:      `cards that up to under 15`,
		leadCard:  `4H`,
		hand:      `AS,2C,3D,4H`,
		expPoints: 10,
	}, {
		desc:      `cards that up to over 46`,
		leadCard:  `6H`,
		hand:      `KS,10C,QD,KH`,
		expPoints: 2,
	}}

	for _, tc := range testCases {
		lead := model.NewCardFromString(tc.leadCard)
		hand := make([]model.Card, 4)
		cardStrs := strings.Split(tc.hand, `,`)
		require.Len(t, cardStrs, 4)
		for i, c := range cardStrs {
			hand[i] = model.NewCardFromString(c)
		}

		actPoints := HandPoints(lead, hand)
		assert.Equal(t, tc.expPoints, actPoints, tc.desc)
	}
}

func TestPointsOddInteractions(t *testing.T) {
	testCases := []struct {
		desc      string
		leadCard  string
		hand      string
		expPoints int
	}{{
		desc:      `flush, double run of 4, and fifteens`,
		leadCard:  `7D`,
		hand:      `6S,7S,8S,9S`,
		expPoints: 20,
	}, {
		desc:      `triplet across lead`,
		leadCard:  `8H`,
		hand:      `8S,8C,QD,10H`,
		expPoints: 6,
	}, {
		desc:      `run of 5 and flush`,
		leadCard:  `AS`,
		hand:      `3S,2S,5S,4S`,
		expPoints: 12,
	}}

	for _, tc := range testCases {
		lead := model.NewCardFromString(tc.leadCard)
		hand := make([]model.Card, 4)
		cardStrs := strings.Split(tc.hand, `,`)
		require.Len(t, cardStrs, 4)
		for i, c := range cardStrs {
			hand[i] = model.NewCardFromString(c)
		}

		actPoints := HandPoints(lead, hand)
		assert.Equal(t, tc.expPoints, actPoints, tc.desc)
	}
}

func TestScoringPoorlySizedHands(t *testing.T) {
	// Asserting zero also checks that the func doesn't panic
	assert.Zero(t, CribPoints(model.Card{}, make([]model.Card, 5)))
	assert.Zero(t, CribPoints(model.Card{}, make([]model.Card, 6)))
	assert.Zero(t, HandPoints(model.Card{}, make([]model.Card, 5)))
	assert.Zero(t, HandPoints(model.Card{}, make([]model.Card, 6)))
}
