package main

import (
	"log"

	"github.com/joshprzybyszewski/cribbage/server"
)

func main() {
	log.Fatal(server.Setup())
}
