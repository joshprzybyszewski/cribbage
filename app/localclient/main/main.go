package main

import (
	"flag"

	"github.com/joshprzybyszewski/cribbage/localclient"
)

func main() {
	flag.Parse()

	err := runClient()
	if err != nil {
		panic(err)
	}
}

func runClient() error {
	return localclient.StartTerminalInteraction()
}
