package main

import (
	"flag"

	"github.com/joshprzybyszewski/cribbage/server"
)

func main() {
	flag.Parse()

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
