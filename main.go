package main

import (
	"github.com/joshprzybyszewski/cribbage/startup"
)

func main() {
	if err := startup.PlayServer(); err != nil {
		panic(err)
	}
}
