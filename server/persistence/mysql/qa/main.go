package main

import (
	"context"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joshprzybyszewski/cribbage/model"
	cribsql "github.com/joshprzybyszewski/cribbage/server/persistence/mysql"
)

func main() {
	per, err := cribsql.New(context.Background())
	if err != nil {
		panic(err)
	}
	err = per.CreatePlayer(model.Player{
		ID:   `c`,
		Name: `connor`,
	})
	if err != nil {
		panic(err)
	}

	player, err := per.GetPlayer(`c`)
	if err != nil {
		panic(err)
	}
	fmt.Printf("got player: %+v\n", player)
}
