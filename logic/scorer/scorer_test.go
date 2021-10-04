package scorer

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshprzybyszewski/cribbage/model"
)

func BenchmarkHandPoints(b *testing.B) {
	hand := randomHand(b, 5)
	b.Run(`scoring random hand`, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s := HandPoints(hand[0], hand[1:])
			require.Less(b, s, 30)
			require.GreaterOrEqual(b, s, 0)
		}
	})
}

func TestPoints(t *testing.T) {
	tests := []struct {
		desc      string
		hand      string
		leadCard  string
		expPoints int
		scoreType scoreType
	}{{
		desc:      `highest scoring hand`,
		leadCard:  `5H`,
		hand:      `5S,5C,5D,JH`,
		expPoints: 29,
		scoreType: quad | fifteen8 | nobs,
	}, {
		desc:      `quad`,
		leadCard:  `8H`,
		hand:      `8S,8C,8D,10H`,
		expPoints: 12,
		scoreType: quad,
	}, {
		desc:      `triplet`,
		leadCard:  `KH`,
		hand:      `8S,8C,8D,10H`,
		expPoints: 6,
		scoreType: triplet,
	}, {
		desc:      `one pair`,
		leadCard:  `KH`,
		hand:      `8S,8C,2D,10H`,
		expPoints: 2,
		scoreType: onepair,
	}, {
		desc:      `two pair`,
		leadCard:  `KH`,
		hand:      `KS,8C,2D,8H`,
		expPoints: 4,
		scoreType: twopair,
	}, {
		desc:      `triple run of 3`,
		leadCard:  `8H`,
		hand:      `9S,8C,10D,8H`,
		expPoints: 15,
		scoreType: tripleRunOfThree,
	}, {
		desc:      `double run of 4`,
		leadCard:  `JH`,
		hand:      `8S,8C,9D,10H`,
		expPoints: 10,
		scoreType: doubleRunOfFour,
	}, {
		desc:      `double double run of 3`,
		leadCard:  `9H`,
		hand:      `8S,8C,9D,10H`,
		expPoints: 16,
		scoreType: doubleDoubleRunOfThree,
	}, {
		desc:      `double run of 3`,
		leadCard:  `KH`,
		hand:      `8S,8C,9D,10H`,
		expPoints: 8,
		scoreType: doubleRunOfThree,
	}, {
		desc:      `run of 3`,
		leadCard:  `KH`,
		hand:      `8S,2C,9D,10H`,
		expPoints: 3,
		scoreType: run3,
	}, {
		desc:      `run of 4`,
		leadCard:  `KH`,
		hand:      `8S,JC,9D,10H`,
		expPoints: 4,
		scoreType: run4,
	}, {
		desc:      `run of 5`,
		leadCard:  `QH`,
		hand:      `8S,JC,9D,10H`,
		expPoints: 5,
		scoreType: run5,
	}, {
		desc:      `only nobs`,
		leadCard:  `6H`,
		hand:      `JH,KC,10D,8H`,
		expPoints: 1,
		scoreType: nobs,
	}, {
		desc:      `flush`,
		leadCard:  `3H`,
		hand:      `8D,4D,10D,6D`,
		expPoints: 6,
		scoreType: flush4 | fifteen1,
	}, {
		desc:      `random hand I got while playing`,
		leadCard:  `4C`,
		hand:      `5H,3D,7D,7S`,
		expPoints: 9,
		scoreType: run3 | fifteen2 | onepair,
	}, {
		desc:      `just looking for ways to break it`,
		leadCard:  `9C`,
		hand:      `1H,5D,7D,7S`,
		expPoints: 6,
		scoreType: fifteen2 | onepair,
	}, {
		desc:      `another hand to break it`,
		leadCard:  `7H`,
		hand:      `6D,6S,10H,9C`,
		expPoints: 6,
		scoreType: fifteen2 | onepair,
	}, {
		desc:      `actual double run of three with a fifteen`,
		leadCard:  `7H`,
		hand:      `6D,6S,10H,8C`,
		expPoints: 10,
		scoreType: doubleRunOfThree | fifteen1,
	}, {
		desc:      `flush, double run of 4, and fifteens`,
		leadCard:  `7D`,
		hand:      `6S,7S,8S,9S`,
		expPoints: 20,
		scoreType: flush4 | doubleRunOfFour | fifteen3,
	}, {
		desc:      `triplet across lead`,
		leadCard:  `8H`,
		hand:      `8S,8C,QD,10H`,
		expPoints: 6,
		scoreType: triplet,
	}, {
		desc:      `run of 5 and flush`,
		leadCard:  `AS`,
		hand:      `3S,2S,5S,4S`,
		expPoints: 12,
		scoreType: flush5 | run5 | fifteen1,
	}, {
		desc:      `15 for 2`,
		leadCard:  `8H`,
		hand:      `7S,AC,2D,KH`,
		expPoints: 2,
		scoreType: fifteen1,
	}, {
		desc:      `run of 5 that adds up to 15`,
		leadCard:  `5H`,
		hand:      `AS,2C,3D,4H`,
		expPoints: 7,
		scoreType: run5 | fifteen1,
	}, {
		desc:      `cards that up to under 15`,
		leadCard:  `4H`,
		hand:      `AS,2C,3D,4H`,
		expPoints: 10,
		scoreType: doubleRunOfFour,
	}, {
		desc:      `cards that up to over 46`,
		leadCard:  `6H`,
		hand:      `KS,10C,QD,KH`,
		expPoints: 2,
		scoreType: onepair,
	}}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			h := parseHand(t, tc.hand)
			cut := parseHand(t, tc.leadCard)
			require.Len(t, cut, 1)

			actPoints := HandPoints(cut[0], h)
			assert.Equal(t, tc.expPoints, actPoints)
			desc, _ := pointsWithDesc(cut[0], h, false)
			assert.Equal(t, tc.scoreType, desc)
		})
	}
}

func TestScoringPoorlySizedHands(t *testing.T) {
	// Asserting zero also checks that the func doesn't panic
	assert.Zero(t, CribPoints(model.Card{}, make([]model.Card, 5)))
	assert.Zero(t, CribPoints(model.Card{}, make([]model.Card, 6)))
	assert.Zero(t, HandPoints(model.Card{}, make([]model.Card, 5)))
	assert.Zero(t, HandPoints(model.Card{}, make([]model.Card, 6)))
}

func parseHand(tb testing.TB, handStr string) []model.Card {
	strs := strings.Split(handStr, `,`)
	hand := make([]model.Card, len(strs))
	for i, s := range strs {
		c, err := model.NewCardFromExternalString(strings.TrimSpace(s))
		require.NoError(tb, err)
		hand[i] = c
	}
	rand.Shuffle(len(hand), func(i, j int) {
		hand[i], hand[j] = hand[j], hand[i]
	})
	return hand
}

func randomHand(t testing.TB, n int) []model.Card {
	if n > 6 {
		t.Fatal(`you really don't need a hand with more than 6 cards, trust me`)
	}
	hand := make([]model.Card, 0, n)
	seen := make(map[model.Card]struct{}, n)
	for len(hand) < n {
		c := model.NewCardFromNumber(rand.Intn(52))
		if _, ok := seen[c]; ok {
			continue
		}
		hand = append(hand, c)
		seen[c] = struct{}{}
	}
	return hand
}
