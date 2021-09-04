package scorer

import (
	"fmt"

	"github.com/joshprzybyszewski/cribbage/model"
)

func HandPoints(lead model.Card, hand []model.Card) int {
	return points(lead, hand, false)
}

func CribPoints(lead model.Card, crib []model.Card) int {
	return points(lead, crib, true)
}

const cardsPerHand = 5

// map card value (1 to 13) to number of occurrences in the hand
// we have to include index 0 and index 14 here because some calculations pass those values in
type valueToCount [15]uint8

func (vtc valueToCount) get(value int) (uint8, bool) {
	return vtc[value], vtc[value] > 0
}

// values is used to store a card's point value or rank value without allocating
type values [cardsPerHand]int

func points(lead model.Card, hand []model.Card, isCrib bool) int {
	if len(hand) != 4 {
		if LOG {
			fmt.Printf("Expected hand size 4, got %d\n", len(hand))
		}
		return 0
	}

	var (
		vals   values
		ptVals values
	)
	allCards := [5]model.Card{hand[0], hand[1], hand[2], hand[3], lead}
	for i, c := range allCards {
		// building up info for later
		vals[i] = c.Value
		ptVals[i] = c.PegValue()
	}

	totalPoints := 0
	var allScoreTypes scoreType

	st, pts := scoreFifteens(ptVals)
	allScoreTypes = allScoreTypes | st
	totalPoints += pts

	st, pts = scoreRunsAndPairs(vals)
	allScoreTypes = allScoreTypes | st
	totalPoints += pts

	st, pts = scoreFlushesAndNobs(lead, hand, isCrib)
	allScoreTypes = allScoreTypes | st
	totalPoints += pts

	numDescribedPoints := describePoints(lead, hand, allScoreTypes)
	if numDescribedPoints != totalPoints && LOG {
		fmt.Println(`error!`)
		fmt.Printf("calced:    %d\n", totalPoints)
		fmt.Printf("described: %d\n", numDescribedPoints)
	}

	return totalPoints
}

// Assumes input is sorted and has len 5
func scoreFifteens(ptVals values) (scoreType, int) {
	if (ptVals[0]|ptVals[1]|ptVals[2]|ptVals[3]|ptVals[4])&1 == 0 {
		// all even numbered cards => no fifteens possible
		return none, 0
	}

	sum := ptVals[0] + ptVals[1] + ptVals[2] + ptVals[3] + ptVals[4]
	if sum == 15 {
		// only one fifteen possible
		return fifteen1, 2
	} else if sum < 15 || sum > 46 {
		return none, 0
	}

	var numFifteens uint

	for i, v := range ptVals {
		many := howManyAddUpTo(15-v, ptVals, i+1)
		numFifteens += many
	}

	st := fifteen0 << numFifteens

	return st, int(numFifteens * 2)
}

func howManyAddUpTo(goal int, ptVals values, start int) uint {
	if start == cardsPerHand {
		return 0
	}

	var many uint
	for i := start; i < cardsPerHand; i++ {
		o := ptVals[i]
		if o > goal {
			continue
		} else if o == goal {
			many++
		} else {
			// o is less than the goal. See what we can find with it
			subWith := howManyAddUpTo(goal-o, ptVals, i+1)
			many += subWith
		}
	}

	return many
}

func scorePairs(values values, valuesToCounts valueToCount) (scoreType, int) {
	pairPoints := 0
	pairType := none
	alreadyHasPair := false
	for n, ct := range valuesToCounts {
		if n == 0 {
			continue
		}
		switch ct {
		case 4:
			return quad, 12
		case 3:
			return triplet, 6
		case 2:
			pairPoints = 2
			pairType = onepair
			if alreadyHasPair {
				return twopair, 4
			}
			alreadyHasPair = true
		}
	}
	return pairType, pairPoints
}

func scoreRuns(values values, valuesToCounts valueToCount) (scoreType, int) {
	usedLongest := uint8(0)
	usedMult := uint8(1)
	var mult uint8
	var longest uint8
	for _, v := range values {
		if _, ok := valuesToCounts.get(v - 1); ok {
			// this is already part of a run; skip calculation
			continue
		}
		longest = 1
		// we're at the potential beginning of a run
		nextUp := v + 1
		mult, _ = valuesToCounts.get(v)
		for ct, ok := valuesToCounts.get(nextUp); ok; ct, ok = valuesToCounts.get(nextUp) {
			mult *= ct
			longest++
			nextUp++
		}
		if longest >= 3 && longest > usedLongest {
			// we have a valid run!
			usedMult = mult
			usedLongest = longest
		}
		// not a valid run, :sadge:
		mult = 1
		longest = 1
	}
	return calculateTypeAndPoints(usedLongest, usedMult)
}

func calculateTypeAndPoints(longest, mult uint8) (scoreType, int) {
	if longest < 3 {
		return none, 0
	}
	st := none
	switch longest {
	case 3:
		switch mult {
		case 1:
			st = run3
		case 2:
			st = doubleRunOfThree
		case 3:
			st = tripleRunOfThree
		case 4:
			st = doubleDoubleRunOfThree
		}
	case 4:
		switch mult {
		case 1:
			st = run4
		case 2:
			st = doubleRunOfFour
		}
	case 5:
		st = run5
	}
	return st, int(longest) * int(mult)
}

func scoreRunsAndPairs(values values) (scoreType, int) { //nolint:gocyclo
	var valuesToCounts valueToCount
	for _, v := range values {
		valuesToCounts[v]++
	}
	pairType, pairPts := scorePairs(values, valuesToCounts)
	runType, runPts := scoreRuns(values, valuesToCounts)
	return pairType | runType, pairPts + runPts

	// min := values[0]
	// max := values[4]

	// // check quad for all hands is the same
	// if values[3] == min || values[1] == max {
	// 	// this is a quad of the either form
	// 	// Y Y Y Y|X (wlog cuz Y|X X X X)
	// 	return quad, 12 /* 4 of a kind is worth 12 */
	// }

	// diffOf5Cards := max - min

	// if diffOf5Cards == 1 {
	// 	// this hand is either
	// 	// A A A|B B (wlog cuz A A|B B B)
	// 	// A A A A|B (wlog cuz A|B B B B)
	// 	// but we've already checked the quad
	// 	// therefore this is a triple and a pair
	// 	return triplet | onepair, 8 /* 6 for a triple, 2 for a pair */
	// } else if diffOf5Cards == 2 {
	// 	// check hands with triples
	// 	if values[1] == values[3] {
	// 		// we know that values[0] = values[1] - 1 and values[3] = values[4] + 1
	// 		// because we already checked quads
	// 		// A|B B B|C
	// 		return tripleRunOfThree, 15 /* triple run of 3 */
	// 	} else if values[2] == min {
	// 		// A A A|B|C
	// 		// A A A|C C
	// 		if values[3] == max {
	// 			// triple/pair combo
	// 			return triplet | onepair, 8 /* 6 for a triple, 2 for a pair */
	// 		}
	// 		// triple run of 3 is worth 15
	// 		// because 3 runs of three (worth 9) and a triple for 6
	// 		return tripleRunOfThree, 15 /* triple run of 3 */
	// 	} else if values[2] == max {
	// 		// A|B|C C C
	// 		// A A|C C C
	// 		if values[1] == min {
	// 			return triplet | onepair, 8 /* 6 for a triple, 2 for a pair */
	// 		}
	// 		return tripleRunOfThree, 15 /* triple run of 3 */
	// 	}

	// 	// If the hand doesn't have a quad or a triple, then it must be a double-double run of three
	// 	return doubleDoubleRunOfThree, 16 /* double double run of 3 */
	// }

	// // checking triplets for remaining hands are all the same
	// // look for (triple/pair, triple/odds)
	// if values[1] == values[3] {
	// 	// this is definitely a triple without a pair
	// 	// A|B B B|D
	// 	return triplet, 6 /* triplet */
	// } else if values[2] == min {
	// 	// A A A|B|D
	// 	// A A A|D D
	// 	if values[3] == max {
	// 		// triple/pair combo
	// 		return triplet | onepair, 8 /* 6 for a triple, 2 for a pair */
	// 	}
	// 	return triplet, 6 /* triplet */
	// } else if values[2] == max {
	// 	// A|B|D D D
	// 	// A A|D D D
	// 	if values[1] == min {
	// 		return triplet | onepair, 8 /* 6 for a triple, 2 for a pair */
	// 	}
	// 	return triplet, 6 /* triplet */
	// }

	// switch diffOf5Cards {
	// case 3:
	// 	// If we have a diff of 3, then we either have a two pair or a double run of 4
	// 	// check double run of four: A A|B|C|D
	// 	// check two-pair          : A A|B B|D
	// 	numPairs := 0
	// 	for i := 1; i < len(values); i++ {
	// 		if values[i] == values[i-1] {
	// 			numPairs++
	// 		}
	// 	}

	// 	if numPairs == 2 {
	// 		// we found two pairs, which mean it cannot be a double run of 4
	// 		return twopair, 4 /* two pair */
	// 	}

	// 	return doubleRunOfFour, 10 /* double run of four */

	// case 4:
	// 	// So here's the options:
	// 	// Run of 5
	// 	// 1,2,3,4,5
	// 	// double run of three
	// 	// 1,2,x,3,5
	// 	// two pair
	// 	// 1,1,3,3,5
	// 	// 1,3,3,5,5
	// 	// pair + run of three
	// 	// 1,2,3,5,5
	// 	// 1,1,3,4,5
	// 	// pair alone (regardless of fifteens)
	// 	// 6,6,7,9,10
	// 	// 6,7,9,10,10
	// 	var numPairs, numIncs int
	// 	for i := 0; i < len(values)-1; i++ {
	// 		if iv, nv := values[i], values[i+1]; iv == nv {
	// 			numPairs++
	// 		} else if iv+1 == nv {
	// 			numIncs++
	// 		}
	// 	}
	// 	if numPairs == 2 {
	// 		return twopair, 4 /* two pair */
	// 	} else if numIncs == 4 {
	// 		return run5, 5 /* run of 5 */
	// 	} else if values[2] == max-2 && (values[3] == max || values[1] == min) {
	// 		return run3 | onepair, 5 /* run of three and a pair*/
	// 	} else if (values[2] == min+1 && values[3] != values[2]+1) ||
	// 		(values[2] == max-1 && values[1] != values[2]-1) {
	// 		return onepair, 2 /* a pair*/
	// 	}

	// 	return doubleRunOfThree, 8 /* double run of 3 */
	// }

	// // check run of 4 (and middle set run of 3)
	// if (values[1]+1 == values[2]) && (values[2]+1 == values[3]) {
	// 	if (min+1 == values[1]) || (values[3]+1 == max) {
	// 		// there's no way all 5 can be, so we just check for 4
	// 		return run4, 4
	// 	}
	// 	// since we were checking for 4, but we found the middle three to be a run,
	// 	// then we can check the two ends for a double run
	// 	if min == values[1] || values[3] == max {
	// 		return doubleRunOfThree, 8 /* double run of 3 */
	// 	}
	// 	return run3, 3 /* run of 3 */
	// }

	// numPairs := 0
	// for i := 1; i < len(values); i++ {
	// 	if values[i] == values[i-1] {
	// 		numPairs++
	// 	}
	// }
	// if numPairs == 2 {
	// 	return twopair, 4 /* two points per pair*/
	// }

	// if values[3]-min == 2 {
	// 	// either two pair or a double run of three
	// 	// but we already checked two pair
	// 	return doubleRunOfThree, 8 /* double run of 3 */
	// } else if max-values[1] == 2 {
	// 	// either two pair or a double run of three
	// 	// but we already checked two pair
	// 	return doubleRunOfThree, 8 /* double run of 3 */
	// }

	// var st scoreType
	// pts := 0

	// if ((min+1 == values[1]) && (values[1]+1 == values[2])) ||
	// 	((values[2]+1 == values[3]) && (values[3]+1 == max)) {
	// 	st = st | run3 /* run of 3 (and maybe a pair) */
	// 	pts += 3
	// }

	// if numPairs == 1 {
	// 	// we already checked for two above
	// 	st = st | onepair /* one pair */
	// 	pts += 2
	// }

	// return st, pts
}

func scoreFlushesAndNobs(lead model.Card, hand []model.Card, isCrib bool) (scoreType, int) {
	var st scoreType
	pts := 0
	handSuit := hand[0].Suit
	isHandFlush := true
	for _, c := range hand {
		if c.Suit != handSuit {
			isHandFlush = false
		}
		if c.Value == 11 /* Jack */ && c.Suit == lead.Suit {
			st = st | nobs
			pts++
		}
	}

	if isHandFlush {
		if lead.Suit == handSuit {
			st = st | flush5
			pts += 5
		} else if !isCrib {
			st = st | flush4
			pts += 4
		}
	}

	return st, pts
}
