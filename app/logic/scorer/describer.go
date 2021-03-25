package scorer

import (
	"fmt"

	"github.com/joshprzybyszewski/cribbage/model"
)

const LOG = false

type scoreType int

const (
	none scoreType = 1 << iota
	quad
	triplet
	onepair
	twopair
	tripleRunOfThree
	doubleRunOfFour
	doubleDoubleRunOfThree
	doubleRunOfThree
	run3
	run4
	run5
	flush4
	flush5
	nobs
	fifteen0
	fifteen1
	fifteen2
	fifteen3
	fifteen4
	fifteen5
	fifteen6
	fifteen7
	fifteen8

	allFifteens = fifteen1 | fifteen2 | fifteen3 | fifteen4 | fifteen5 | fifteen6 | fifteen7 | fifteen8
)

func describePoints(lead model.Card, hand []model.Card, st scoreType) int { //nolint:gocyclo
	if LOG {
		fmt.Println("---------")
		fmt.Printf("lead: %v\n", lead.String())
		fmt.Printf("hand: %v, %v, %v, %v\n", hand[0].String(), hand[1].String(), hand[2].String(), hand[3].String())
	}

	runningTotal := 0

	for runningFifteen, notter := st&allFifteens, fifteen1; runningFifteen != 0; {
		runningTotal += 2
		logPoints("fifteen", runningTotal)
		runningFifteen = runningFifteen & ^notter
		notter = notter << 1
	}

	if st&quad != 0 {
		runningTotal += 12
		logPoints("quad", runningTotal)
	}

	if st&tripleRunOfThree != 0 {
		runningTotal += 15
		logPoints("triple run of 3", runningTotal)
	}

	if st&doubleRunOfFour != 0 {
		runningTotal += 10
		logPoints("double run of 4", runningTotal)
	}

	if st&doubleDoubleRunOfThree != 0 {
		runningTotal += 16
		logPoints("double double run of 3", runningTotal)
	}

	if st&doubleRunOfThree != 0 {
		runningTotal += 8
		logPoints("double run of 3", runningTotal)
	}

	if st&triplet != 0 {
		runningTotal += 6
		logPoints("triplet", runningTotal)
	}

	if st&onepair != 0 {
		runningTotal += 2
		logPoints("pair", runningTotal)
	}

	if st&twopair != 0 {
		runningTotal += 4
		logPoints("two-pair", runningTotal)
	}

	if st&run5 != 0 {
		runningTotal += 5
		logPoints("run of 5", runningTotal)
	}

	if st&run4 != 0 {
		runningTotal += 4
		logPoints("run of 4", runningTotal)
	}

	if st&run3 != 0 {
		runningTotal += 3
		logPoints("run of 3", runningTotal)
	}

	if st&flush4 != 0 {
		runningTotal += 4
		logPoints("flush of 4", runningTotal)
	}

	if st&flush5 != 0 {
		runningTotal += 5
		logPoints("flush of 5", runningTotal)
	}

	if st&nobs != 0 {
		runningTotal++
		logPoints("nobs", runningTotal)
	}

	if LOG {
		fmt.Println("=========")
	}

	return runningTotal
}

func logPoints(msg string, points int) {
	if LOG {
		fmt.Printf("(%2d) %v\n", points, msg)
	}
}
