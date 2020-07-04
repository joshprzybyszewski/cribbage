package main

import (
	"github.com/joshprzybyszewski/cribbage/server"
)

func main() {
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
