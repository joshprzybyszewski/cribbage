package main

import (
	"context"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joshprzybyszewski/cribbage/model"
	cribsql "github.com/joshprzybyszewski/cribbage/server/persistence/mysql"
)

func main() {
	mysql, err := cribsql.New(context.Background())
	if err != nil {
		panic(err)
	}
	err = mysql.CreatePlayer(model.Player{
		ID:   `c`,
		Name: `connor`,
	})
	if err != nil {
		panic(err)
	}

}
