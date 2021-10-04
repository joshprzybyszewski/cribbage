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

const (
	numCardsToScore = 5
	// we only need 15, but due to memory alignment, doing this actually
	// saves us 8 bytes in iterateHandResult
	valuesToCountsLength = 16
)

type pairScorer struct {
	valuesToCounts [valuesToCountsLength]uint8
	points         uint8
}

func (ps *pairScorer) atEach(c model.Card) {
	ps.valuesToCounts[c.Value]++
	switch ps.valuesToCounts[c.Value] {
	case 2:
		ps.points += 2
	case 3:
		// we've already added 2 for this value
		ps.points += 4
	case 4:
		// we've already added 6 for this value
		ps.points += 6
	}
}

func (ps *pairScorer) accumulate() (scoreType, uint8) {
	var st scoreType
	switch ps.points {
	case 2:
		st = onepair
	case 4:
		st = twopair
	case 6:
		st = triplet
	case 8:
		st = triplet | onepair
	case 12:
		st = quad
	}
	return st, ps.points
}

type flushScorer struct {
	isHandFlush bool
	isCrib      bool
	firstCard   model.Card
	leadCard    model.Card
}

func newFlushScorer(lead model.Card, hand []model.Card, isCrib bool) flushScorer {
	return flushScorer{
		isHandFlush: true, // innocent until proven guilty
		isCrib:      isCrib,
		firstCard:   hand[0],
		leadCard:    lead,
	}
}

func (fs *flushScorer) atEach(c model.Card) {
	if c == fs.leadCard {
		return
	}
	if c.Suit != fs.firstCard.Suit {
		fs.isHandFlush = false
	}
}

func (fs *flushScorer) accumulate() (scoreType, uint8) {
	if !fs.isHandFlush {
		return 0, 0
	}
	if fs.firstCard.Suit == fs.leadCard.Suit {
		return flush5, 5
	}
	if fs.isCrib {
		return 0, 0
	}
	return flush4, 4
}

type iterateHandResult struct {
	valuesToCounts [valuesToCountsLength]uint8
	values         [numCardsToScore]int
	ptValues       [numCardsToScore]int
	scoreType      scoreType
	totalPoints    uint8
}

// iterateHand does all the scoring possible with a single iteration through the hand
// we can't do fifteens or runs here directly, but the information returned will aid in doing those efficiently
func iterateHand(lead model.Card, hand []model.Card, isCrib bool) iterateHandResult {
	res := iterateHandResult{}
	hasNobs := false
	ps := &pairScorer{}
	fs := newFlushScorer(lead, hand, isCrib)
	for i, c := range [5]model.Card{hand[0], hand[1], hand[2], hand[3], lead} {
		ps.atEach(c)
		fs.atEach(c)
		if c != lead && c.Value == model.JackValue && c.Suit == lead.Suit {
			hasNobs = true
		}
		res.values[i] = c.Value
		res.ptValues[i] = c.PegValue()
	}

	if hasNobs {
		res.totalPoints++
		res.scoreType |= nobs
	}
	fSt, fPts := fs.accumulate()
	res.totalPoints += fPts
	res.scoreType |= fSt

	pSt, pPts := ps.accumulate()
	res.totalPoints += pPts
	res.scoreType |= pSt
	res.valuesToCounts = ps.valuesToCounts

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
