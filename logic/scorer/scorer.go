package scorer

import (
	"fmt"
	"sort"

	"github.com/joshprzybyszewski/cribbage/model"
)

func HandPoints(lead model.Card, hand []model.Card) int {
	return points(lead, hand, false)
}

func CribPoints(lead model.Card, crib []model.Card) int {
	return points(lead, crib, true)
}

func points(lead model.Card, hand []model.Card, isCrib bool) int {
	if len(hand) != 4 {
		if LOG {
			fmt.Printf("Expected hand size 4, got %d\n", len(hand))
		}
		return 0
	}

	allScoreTypes, totalPoints := pointsWithDesc(lead, hand, isCrib)

	numDescribedPoints := describePoints(lead, hand, allScoreTypes)
	if numDescribedPoints != totalPoints && LOG {
		fmt.Println(`error!`)
		fmt.Printf("calced:    %d\n", totalPoints)
		fmt.Printf("described: %d\n", numDescribedPoints)
	}

	return totalPoints
}

const numCardsToScore = 5

func pointsWithDesc(lead model.Card, hand []model.Card, isCrib bool) (scoreType, int) {
	iterationResult := iterateHand(lead, hand, isCrib)

	fSt, fPts := scoreFifteens(iterationResult.ptValues)
	iterationResult.scoreType |= fSt
	iterationResult.totalPoints += uint8(fPts)

	rSt, rPts := scoreRuns(iterationResult.values, iterationResult.valuesToCounts[:])
	iterationResult.scoreType |= rSt
	iterationResult.totalPoints += uint8(rPts)

	// resolve pairs/runs score types, since e.g. tripleRunOfThree includes the triplet
	// and the triplet bit shouldn't be set
	iterationResult.scoreType = resolveScoreType(iterationResult.scoreType)

	return iterationResult.scoreType, int(iterationResult.totalPoints)
}

type iterateHandResult struct {
	valuesToCounts [15]uint8
	values         [numCardsToScore]int
	ptValues       [numCardsToScore]int
	scoreType      scoreType
	totalPoints    uint8

	pairPoints uint8
}

// iterateHand does all the scoring possible with a single iteration through the hand
// we can't do fifteens or runs here directly, but the information returned will aid in doing those efficiently
func iterateHand(cut model.Card, hand []model.Card, isCrib bool) iterateHandResult {
	res := iterateHandResult{}
	numSuitsMatching := uint8(0)
	hasNobs := false
	for i, c := range [5]model.Card{hand[0], hand[1], hand[2], hand[3], cut} {
		// pairs
		res.valuesToCounts[c.Value]++
		switch res.valuesToCounts[c.Value] {
		case 2:
			res.pairPoints += 2
		case 3:
			// we've already added 2 for this value
			res.pairPoints += 4
		case 4:
			// we've already added 6 for this value
			res.pairPoints += 6
		}

		// flushes
		if c.Suit == hand[0].Suit && c != cut {
			numSuitsMatching++
		}

		// nobs
		if c != cut && c.Value == model.JackValue && c.Suit == cut.Suit {
			hasNobs = true
		}

		// value set
		res.values[i] = c.Value
		res.ptValues[i] = c.PegValue()
	}

	// accumulate results
	// nobs
	if hasNobs {
		res.totalPoints++
		res.scoreType |= nobs
	}
	// flushes
	if numSuitsMatching == 4 {
		if hand[0].Suit == cut.Suit {
			res.totalPoints += 5
			res.scoreType |= flush5
		} else if !isCrib {
			res.totalPoints += 4
			res.scoreType |= flush4
		}
	}
	// pairs
	switch res.pairPoints {
	case 2:
		res.scoreType |= onepair
	case 4:
		res.scoreType |= twopair
	case 6:
		res.scoreType |= triplet
	case 8:
		res.scoreType |= triplet | onepair
	case 12:
		res.scoreType |= quad
	}
	res.totalPoints += res.pairPoints

	return res
}

func scoreFifteens(ptVals [numCardsToScore]int) (scoreType, int) {
	if (ptVals[0]|ptVals[1]|ptVals[2]|ptVals[3]|ptVals[4])&1 == 0 {
		// all even numbered cards => no fifteens possible
		return 0, 0
	}

	sum := ptVals[0] + ptVals[1] + ptVals[2] + ptVals[3] + ptVals[4]
	if sum == 15 {
		// only one fifteen possible
		return fifteen1, 2
	} else if sum < 15 || sum > 46 {
		return 0, 0
	}

	var numFifteens int

	for i := 0; i < len(ptVals); i++ {
		many := howManyAddUpTo(15-ptVals[i], ptVals[i+1:])
		numFifteens += many
	}

	if numFifteens > 0 {
		return fifteen0 << numFifteens, numFifteens * 2
	}
	return 0, 0
}

func howManyAddUpTo(goal int, ptVals []int) int {
	if len(ptVals) == 0 {
		return 0
	}

	var many int
	for i, o := range ptVals {
		if o == goal {
			many++
		} else if o < goal {
			// o is less than the goal. See what we can find with it
			subWith := howManyAddUpTo(goal-o, ptVals[i+1:])
			many += subWith
		}
	}

	return many
}

func scoreRuns(values [numCardsToScore]int, valuesToCounts []uint8) (scoreType, int) {
	for _, v := range values {
		if valuesToCounts[v-1] > 0 {
			// this is already part of a run; skip calculation
			continue
		}
		runLen := uint8(1)
		mult := valuesToCounts[v]
		// we're at the potential beginning of a run
		for next := v + 1; valuesToCounts[next] > 0; next++ {
			mult *= valuesToCounts[next]
			runLen++
		}
		if runLen >= 3 {
			return calculateTypeAndPoints(runLen, mult)
		}
	}
	return 0, 0
}

func calculateTypeAndPoints(longest, mult uint8) (scoreType, int) {
	if longest < 3 {
		return 0, 0
	}
	if longest == 5 {
		return run5, 5
	}
	// typeMap maps the run length and multiplier to a scoring type
	typeMap := [5][5]scoreType{
		3: {1: run3, 2: doubleRunOfThree, 3: tripleRunOfThree, 4: doubleDoubleRunOfThree},
		4: {1: run4, 2: doubleRunOfFour},
	}
	return typeMap[longest][mult], int(longest) * int(mult)
}

func resolveScoreType(st scoreType) scoreType {
	// these masks are generated as following:
	// 00011111
	//    ^ doubleRunOfThree bit
	// 00000011
	//      ^ tripleRunOfThree bit (minus one to get the two bits below)
	// they include the bits set for th
	runsWithPairsMask := generateBitMask(tripleRunOfThree, doubleRunOfThree)
	if runsWithPairs := runsWithPairsMask & st; runsWithPairs > 0 {
		// we have a run that has a pair baked in, like tripleRunOfThree
		// let's mask off quad, triplet, onepair, and twopair
		pairsMask := generateBitMask(quad, twopair)
		return st &^ pairsMask
	}
	return st
}

// generateBitMask generates a bit mask with the bits from min to max (inclusive) set and the rest unset
func generateBitMask(min, max scoreType) scoreType {
	// these masks are generated as following:
	// 00011111
	//    ^ max bit
	// 00000011
	//      ^ min bit (minus one to get the two bits below)
	return ^(min - 1) & ((max << 1) - 1)
}

func scoreRunsAndPairs(values []int) (scoreType, int) { //nolint:gocyclo,go-staticcheck
	// this logic assumes we're sorted
	sort.Ints(values)

	min := values[0]
	max := values[4]

	// check quad for all hands is the same
	if values[3] == min || values[1] == max {
		// this is a quad of the either form
		// Y Y Y Y|X (wlog cuz Y|X X X X)
		return quad, 12 /* 4 of a kind is worth 12 */
	}

	diffOf5Cards := max - min

	if diffOf5Cards == 1 {
		// this hand is either
		// A A A|B B (wlog cuz A A|B B B)
		// A A A A|B (wlog cuz A|B B B B)
		// but we've already checked the quad
		// therefore this is a triple and a pair
		return triplet | onepair, 8 /* 6 for a triple, 2 for a pair */
	} else if diffOf5Cards == 2 {
		// check hands with triples
		if values[1] == values[3] {
			// we know that values[0] = values[1] - 1 and values[3] = values[4] + 1
			// because we already checked quads
			// A|B B B|C
			return tripleRunOfThree, 15 /* triple run of 3 */
		} else if values[2] == min {
			// A A A|B|C
			// A A A|C C
			if values[3] == max {
				// triple/pair combo
				return triplet | onepair, 8 /* 6 for a triple, 2 for a pair */
			}
			// triple run of 3 is worth 15
			// because 3 runs of three (worth 9) and a triple for 6
			return tripleRunOfThree, 15 /* triple run of 3 */
		} else if values[2] == max {
			// A|B|C C C
			// A A|C C C
			if values[1] == min {
				return triplet | onepair, 8 /* 6 for a triple, 2 for a pair */
			}
			return tripleRunOfThree, 15 /* triple run of 3 */
		}

		// If the hand doesn't have a quad or a triple, then it must be a double-double run of three
		return doubleDoubleRunOfThree, 16 /* double double run of 3 */
	}

	// checking triplets for remaining hands are all the same
	// look for (triple/pair, triple/odds)
	if values[1] == values[3] {
		// this is definitely a triple without a pair
		// A|B B B|D
		return triplet, 6 /* triplet */
	} else if values[2] == min {
		// A A A|B|D
		// A A A|D D
		if values[3] == max {
			// triple/pair combo
			return triplet | onepair, 8 /* 6 for a triple, 2 for a pair */
		}
		return triplet, 6 /* triplet */
	} else if values[2] == max {
		// A|B|D D D
		// A A|D D D
		if values[1] == min {
			return triplet | onepair, 8 /* 6 for a triple, 2 for a pair */
		}
		return triplet, 6 /* triplet */
	}

	switch diffOf5Cards {
	case 3:
		// If we have a diff of 3, then we either have a two pair or a double run of 4
		// check double run of four: A A|B|C|D
		// check two-pair          : A A|B B|D
		numPairs := 0
		for i := 1; i < len(values); i++ {
			if values[i] == values[i-1] {
				numPairs++
			}
		}

		if numPairs == 2 {
			// we found two pairs, which mean it cannot be a double run of 4
			return twopair, 4 /* two pair */
		}

		return doubleRunOfFour, 10 /* double run of four */

	case 4:
		// So here's the options:
		// Run of 5
		// 1,2,3,4,5
		// double run of three
		// 1,2,x,3,5
		// two pair
		// 1,1,3,3,5
		// 1,3,3,5,5
		// pair + run of three
		// 1,2,3,5,5
		// 1,1,3,4,5
		// pair alone (regardless of fifteens)
		// 6,6,7,9,10
		// 6,7,9,10,10
		var numPairs, numIncs int
		for i := 0; i < len(values)-1; i++ {
			if iv, nv := values[i], values[i+1]; iv == nv {
				numPairs++
			} else if iv+1 == nv {
				numIncs++
			}
		}
		if numPairs == 2 {
			return twopair, 4 /* two pair */
		} else if numIncs == 4 {
			return run5, 5 /* run of 5 */
		} else if values[2] == max-2 && (values[3] == max || values[1] == min) {
			return run3 | onepair, 5 /* run of three and a pair*/
		} else if (values[2] == min+1 && values[3] != values[2]+1) ||
			(values[2] == max-1 && values[1] != values[2]-1) {
			return onepair, 2 /* a pair*/
		}

		return doubleRunOfThree, 8 /* double run of 3 */
	}

	// check run of 4 (and middle set run of 3)
	if (values[1]+1 == values[2]) && (values[2]+1 == values[3]) {
		if (min+1 == values[1]) || (values[3]+1 == max) {
			// there's no way all 5 can be, so we just check for 4
			return run4, 4
		}
		// since we were checking for 4, but we found the middle three to be a run,
		// then we can check the two ends for a double run
		if min == values[1] || values[3] == max {
			return doubleRunOfThree, 8 /* double run of 3 */
		}
		return run3, 3 /* run of 3 */
	}

	numPairs := 0
	for i := 1; i < len(values); i++ {
		if values[i] == values[i-1] {
			numPairs++
		}
	}
	if numPairs == 2 {
		return twopair, 4 /* two points per pair*/
	}

	if values[3]-min == 2 {
		// either two pair or a double run of three
		// but we already checked two pair
		return doubleRunOfThree, 8 /* double run of 3 */
	} else if max-values[1] == 2 {
		// either two pair or a double run of three
		// but we already checked two pair
		return doubleRunOfThree, 8 /* double run of 3 */
	}

	var st scoreType
	pts := 0

	if ((min+1 == values[1]) && (values[1]+1 == values[2])) ||
		((values[2]+1 == values[3]) && (values[3]+1 == max)) {
		st = st | run3 /* run of 3 (and maybe a pair) */
		pts += 3
	}

	if numPairs == 1 {
		// we already checked for two above
		st = st | onepair /* one pair */
		pts += 2
	}

	return st, pts
}
