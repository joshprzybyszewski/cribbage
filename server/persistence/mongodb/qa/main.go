package main

import (
	"fmt"

	"github.com/joshprzybyszewski/cribbage/server/persistence/mongodb"
)

func main() {
	dbf, err := mongodb.NewFactory(``)
	fmt.Printf("dbf, err := %+v, %+v", dbf, err)
}
