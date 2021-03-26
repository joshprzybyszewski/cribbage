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

	values := make([]int, 5)
	ptValues := make([]int, 5)
	for i, c := range hand {
		values[i] = c.Value
		ptValues[i] = c.PegValue()
	}
	values[4] = lead.Value
	ptValues[4] = lead.PegValue()

	sort.Ints(values)
	sort.Ints(ptValues)

	totalPoints := 0
	var allScoreTypes scoreType

	st, pts := scoreFifteens(ptValues)
	allScoreTypes = allScoreTypes | st
	totalPoints += pts

	st, pts = scoreRunsAndPairs(values)
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
func scoreFifteens(ptVals []int) (scoreType, int) {
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

	for i := 0; i < len(ptVals) && ptVals[i] < 8; i++ {
		many := howManyAddUpTo(15-ptVals[i], ptVals[i+1:])
		numFifteens += many
	}

	st := fifteen0 << numFifteens

	return st, int(numFifteens * 2)
}

func howManyAddUpTo(goal int, ptVals []int) uint {
	if len(ptVals) == 0 {
		return 0
	}

	var many uint
	for i, o := range ptVals {
		if o > goal {
			break
		} else if o == goal {
			many++
		} else {
			// o is less than the goal. See what we can find with it
			subWith := howManyAddUpTo(goal-o, ptVals[i+1:])
			many += subWith
		}
	}

	return many
}

// Assumes input is sorted and has len 5
func scoreRunsAndPairs(values []int) (scoreType, int) { //nolint:gocyclo
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
