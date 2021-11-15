package main

import (
	"log"

	"github.com/joshprzybyszewski/cribbage/server"
)

func main() {
	if err := server.Setup(); err != nil {
		log.Fatal(err)
	}
}
