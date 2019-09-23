package play

import (
	"errors"

	"github.com/joshprzybyszewski/cribbage/game"
)

func PlayGame() {
	g := game.New()
	for !g.IsOver() {
		// shuffle
		// deal
		// build crib
		// cut
		// peg
		// count
		// count crib
	}
}
