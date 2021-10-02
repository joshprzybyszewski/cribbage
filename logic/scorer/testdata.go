//+build !prod

package scorer

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/stretchr/testify/require"
)

type pointsType int

const (
	pointsTypeRuns pointsType = iota
	pointsTypePairs
	pointsTypeFifteens
	pointsTypeFlushes
	pointsTypeNobs
)

type testHand struct {
	desc           string
	hand           []model.Card
	cut            model.Card
	points         int
	scoreType      scoreType
	pointsByType   map[pointsType]int
	values         [numCardsToScore]int
	valuesToCounts [15]uint8
}

func scoringTestCases(tb testing.TB) []testHand {
	rawHands := []struct {
		desc         string
		hand         string
		cut          string
		points       int
		scoreType    scoreType
		pointsByType map[pointsType]int
	}{{
		desc:      `highest scoring hand`,
		cut:       `5H`,
		hand:      `5S,5C,5D,JH`,
		points:    29,
		scoreType: quad | fifteen8 | nobs,
		pointsByType: map[pointsType]int{
			pointsTypePairs:    12,
			pointsTypeFifteens: 16,
			pointsTypeNobs:     1,
		},
	}, {
		desc:      `quad`,
		cut:       `8H`,
		hand:      `8S,8C,8D,10H`,
		points:    12,
		scoreType: quad,
		pointsByType: map[pointsType]int{
			pointsTypePairs: 12,
		},
	}, {
		desc:      `triplet`,
		cut:       `KH`,
		hand:      `8S,8C,8D,10H`,
		points:    6,
		scoreType: triplet,
		pointsByType: map[pointsType]int{
			pointsTypePairs: 6,
		},
	}, {
		desc:      `one pair`,
		cut:       `KH`,
		hand:      `8S,8C,2D,10H`,
		points:    2,
		scoreType: onepair,
		pointsByType: map[pointsType]int{
			pointsTypePairs: 2,
		},
	}, {
		desc:      `two pair`,
		cut:       `KH`,
		hand:      `KS,8C,2D,8H`,
		points:    4,
		scoreType: twopair,
		pointsByType: map[pointsType]int{
			pointsTypePairs: 4,
		},
	}, {
		desc:      `triple run of 3`,
		cut:       `8H`,
		hand:      `9S,8C,10D,8H`,
		points:    15,
		scoreType: tripleRunOfThree,
		pointsByType: map[pointsType]int{
			pointsTypePairs: 6,
			pointsTypeRuns:  9,
		},
	}, {
		desc:      `double run of 4`,
		cut:       `JH`,
		hand:      `8S,8C,9D,10H`,
		points:    10,
		scoreType: doubleRunOfFour,
		pointsByType: map[pointsType]int{
			pointsTypePairs: 2,
			pointsTypeRuns:  8,
		},
	}, {
		desc:      `double double run of 3`,
		cut:       `9H`,
		hand:      `8S,8C,9D,10H`,
		points:    16,
		scoreType: doubleDoubleRunOfThree,
		pointsByType: map[pointsType]int{
			pointsTypePairs: 4,
			pointsTypeRuns:  12,
		},
	}, {
		desc:      `double run of 3`,
		cut:       `KH`,
		hand:      `8S,8C,9D,10H`,
		points:    8,
		scoreType: doubleRunOfThree,
		pointsByType: map[pointsType]int{
			pointsTypePairs: 2,
			pointsTypeRuns:  6,
		},
	}, {
		desc:      `run of 3`,
		cut:       `KH`,
		hand:      `8S,2C,9D,10H`,
		points:    3,
		scoreType: run3,
		pointsByType: map[pointsType]int{
			pointsTypeRuns: 3,
		},
	}, {
		desc:      `run of 4`,
		cut:       `KH`,
		hand:      `8S,JC,9D,10H`,
		points:    4,
		scoreType: run4,
		pointsByType: map[pointsType]int{
			pointsTypeRuns: 4,
		},
	}, {
		desc:      `run of 5`,
		cut:       `QH`,
		hand:      `8S,JC,9D,10H`,
		points:    5,
		scoreType: run5,
		pointsByType: map[pointsType]int{
			pointsTypeRuns: 5,
		},
	}, {
		desc:      `only nobs`,
		cut:       `6H`,
		hand:      `JH,KC,10D,8H`,
		points:    1,
		scoreType: nobs,
		pointsByType: map[pointsType]int{
			pointsTypeNobs: 1,
		},
	}, {
		desc:      `flush`,
		cut:       `3H`,
		hand:      `8D,4D,10D,6D`,
		points:    6,
		scoreType: flush4 | fifteen1,
		pointsByType: map[pointsType]int{
			pointsTypeFifteens: 2,
			pointsTypeFlushes:  4,
		},
	}, {
		desc:      `random hand I got while playing`,
		cut:       `4C`,
		hand:      `5H,3D,7D,7S`,
		points:    9,
		scoreType: run3 | fifteen2 | onepair,
		pointsByType: map[pointsType]int{
			pointsTypePairs:    2,
			pointsTypeRuns:     3,
			pointsTypeFifteens: 4,
		},
	}, {
		desc:      `just looking for ways to break it`,
		cut:       `9C`,
		hand:      `1H,5D,7D,7S`,
		points:    6,
		scoreType: fifteen2 | onepair,
		pointsByType: map[pointsType]int{
			pointsTypePairs:    2,
			pointsTypeFifteens: 4,
		},
	}, {
		desc:      `another hand to break it`,
		cut:       `7H`,
		hand:      `6D,6S,10H,9C`,
		points:    6,
		scoreType: fifteen2 | onepair,
		pointsByType: map[pointsType]int{
			pointsTypePairs:    2,
			pointsTypeFifteens: 4,
		},
	}, {
		desc:      `actual double run of three with a fifteen`,
		cut:       `7H`,
		hand:      `6D,6S,10H,8C`,
		points:    10,
		scoreType: doubleRunOfThree | fifteen1,
		pointsByType: map[pointsType]int{
			pointsTypeRuns:     6,
			pointsTypePairs:    2,
			pointsTypeFifteens: 2,
		},
	}, {
		desc:      `flush, double run of 4, and fifteens`,
		cut:       `7D`,
		hand:      `6S,7S,8S,9S`,
		points:    20,
		scoreType: flush4 | doubleRunOfFour | fifteen3,
		pointsByType: map[pointsType]int{
			pointsTypeRuns:     8,
			pointsTypePairs:    2,
			pointsTypeFifteens: 6,
			pointsTypeFlushes:  4,
		},
	}, {
		desc:      `triplet across lead`,
		cut:       `8H`,
		hand:      `8S,8C,QD,10H`,
		points:    6,
		scoreType: triplet,
		pointsByType: map[pointsType]int{
			pointsTypePairs: 6,
		},
	}, {
		desc:      `run of 5 and flush`,
		cut:       `AS`,
		hand:      `3S,2S,5S,4S`,
		points:    12,
		scoreType: flush5 | run5 | fifteen1,
		pointsByType: map[pointsType]int{
			pointsTypeRuns:     5,
			pointsTypeFifteens: 2,
			pointsTypeFlushes:  5,
		},
	}, {
		desc:      `15 for 2`,
		cut:       `8H`,
		hand:      `7S,AC,2D,KH`,
		points:    2,
		scoreType: fifteen1,
		pointsByType: map[pointsType]int{
			pointsTypeFifteens: 2,
		},
	}, {
		desc:      `run of 5 that adds up to 15`,
		cut:       `5H`,
		hand:      `AS,2C,3D,4H`,
		points:    7,
		scoreType: run5 | fifteen1,
		pointsByType: map[pointsType]int{
			pointsTypeRuns:     5,
			pointsTypeFifteens: 2,
		},
	}, {
		desc:      `cards that up to under 15`,
		cut:       `4H`,
		hand:      `AS,2C,3D,4H`,
		points:    10,
		scoreType: doubleRunOfFour,
		pointsByType: map[pointsType]int{
			pointsTypeRuns:  8,
			pointsTypePairs: 2,
		},
	}, {
		desc:      `cards that up to over 46`,
		cut:       `6H`,
		hand:      `KS,10C,QD,KH`,
		points:    2,
		scoreType: onepair,
		pointsByType: map[pointsType]int{
			pointsTypePairs: 2,
		},
	}}

	hands := make([]testHand, len(rawHands))
	for i, h := range rawHands {
		cutCards := parseHand(tb, h.cut)
		require.LessOrEqual(tb, len(cutCards), 1)
		var cut model.Card
		if len(cutCards) == 1 {
			cut = cutCards[0]
		}

		total := 0
		for _, s := range h.pointsByType {
			total += s
		}
		require.Equal(tb, h.points, total, `test case points mismatch`, h.desc)

		hands[i] = testHand{
			desc:         h.desc,
			hand:         parseHand(tb, h.hand),
			cut:          cut,
			points:       h.points,
			scoreType:    h.scoreType,
			pointsByType: h.pointsByType,
		}

		for j, c := range hands[i].hand {
			hands[i].values[j] = c.Value
			hands[i].valuesToCounts[c.Value]++
		}
		hands[i].values[4] = hands[i].cut.Value
		hands[i].valuesToCounts[hands[i].cut.Value]++
	}
	return hands
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
