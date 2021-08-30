package main

import (
	"fmt"

	"github.com/joshprzybyszewski/cribbage/server/persistence/dynamo"
)

func main() {
	dbf, err := dynamo.NewFactory(``)
	fmt.Printf("dbf, err := %+v, %+v", dbf, err)
}
