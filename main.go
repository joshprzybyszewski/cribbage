package main

import (
	"github.com/joshprzybyszewski/cribbage/play"
)

func main() {
	err := play.PlayGame()
	if err != nil {
		panic(err)
	}
}
