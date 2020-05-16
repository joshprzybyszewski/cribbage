package main

import (
	"flag"

	"github.com/joshprzybyszewski/cribbage/localclient"
	"github.com/joshprzybyszewski/cribbage/server"
)

var (
	client = flag.Bool(`client`, false, `set to true to talk to run as a terminal client against the server`)
)

func main() {

	flag.Parse()

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
	server.Setup()

	return nil
}

func runClient() error {
	return localclient.StartTerminalInteraction()
}
