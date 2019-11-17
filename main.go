package main

import (
	"flag"

	"github.com/joshprzybyszewski/cribbage/local_client"
	"github.com/joshprzybyszewski/cribbage/play"
	"github.com/joshprzybyszewski/cribbage/server"
)

func main() {
	legacy := flag.Bool(`legacy`, false, `set to true to play the legacy style game`)
	client := flag.Bool(`client`, false, `set to true to talk to run as a terminal client against the server`)

	flag.Parse()

	if *legacy {
		err := playLegacy()
		if err != nil {
			panic(err)
		}
		return
	}

	if *client {
		err := runClient()
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

func runClient() error {
	return local_client.StartTerminalInteraction()
}

func playLegacy() error {
	err := play.PlayGame()
	if err != nil {
		panic(err)
	}

	return nil
}
