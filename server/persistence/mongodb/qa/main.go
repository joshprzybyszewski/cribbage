package main

import (
	"context"
	"fmt"

	"github.com/joshprzybyszewski/cribbage/server/persistence/mongodb"
)

func main() {
	db, err := mongodb.New(context.Background(), ``)
	fmt.Printf("db, err := %+v, %+v", db, err)
}
