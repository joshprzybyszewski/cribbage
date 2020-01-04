package main

import (
	"context"
	"fmt"

	"github.com/joshprzybyszewski/cribbage/utils/rand"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joshprzybyszewski/cribbage/model"
	cribsql "github.com/joshprzybyszewski/cribbage/server/persistence/mysql"
)

func main() {
	per, err := cribsql.New(context.Background())
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	id := model.PlayerID(rand.String(5))
	err = per.CreatePlayer(model.Player{
		ID:   id,
		Name: `connor`,
		Games: map[model.GameID]model.PlayerColor{
			1:  model.Blue,
			7:  model.Red,
			11: model.Green,
		},
	})
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	} else {
		fmt.Println(`created player.`)
	}

	player, err := per.GetPlayer(id)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	fmt.Printf("got player: %+v\n", player)
}
