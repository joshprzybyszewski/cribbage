package main

import (
	"flag"

	"github.com/joshprzybyszewski/cribbage/play"
	"github.com/joshprzybyszewski/cribbage/server"
)

func main() {
	legacy := flag.Bool(`legacy`, false, `set to true to play the legacy style game`)

	if *legacy {
		err := playLegacy()
		if err != nil {
			panic(err)
		}
		return
	}

	err := playServer()
	if err != nil {
		panic(err)
	}
}

func playServer() error {
	err := server.Setup()
	if err != nil {
		return err
	}

	return nil
}

func playLegacy() error {
	err := play.PlayGame()
	if err != nil {
		panic(err)
	}

	return nil
}
