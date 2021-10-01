package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/persistence/dynamo"
)

func main() {
	dbf, err := dynamo.NewFactory(``)
	if err != nil {
		log.Fatalf("dynamo.NewFactory err: %+v", err)
	}
	fmt.Printf("dbf, err := %+v\n", dbf)

	ctx, toFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer toFn()
	dw, err := dbf.New(ctx)
	if err != nil {
		log.Fatalf("dbf.New(context.Background()) err: %+v", err)
	}

	pID := model.PlayerID(`joshysquashy3`)
	p := model.Player{
		ID:    pID,
		Name:  `jesus is king`,
		Games: map[model.GameID]model.PlayerColor{
			// model.GameID(123): model.Blue,
			// model.GameID(456): model.Red,
			// model.GameID(789): model.Green,
		},
	}
	fmt.Printf("calling dw.CreatePlayer(%+v)\n", p)
	err = dw.CreatePlayer(p)
	if err != nil {
		log.Printf("dw.CreatePlayer err: %+v", err)
	}

	fmt.Printf("calling dw.GetPlayer(%+v)\n", pID)
	p2, err := dw.GetPlayer(pID)
	if err != nil {
		log.Fatalf("dw.GetPlayer err: %+v", err)
	}
	fmt.Printf("player p2 := %+v\n", p2)

	g := model.Game{
		ID: 4,
	}

	fmt.Printf("calling dw.CreateGame(%+v)\n", g)
	err = dw.CreateGame(g)
	if err != nil {
		log.Fatalf("dw.CreateGame err: %+v", err)
	}

	fmt.Printf("calling dw.CreateGame(%+v)\n", g)
	gg, err := dw.GetGame(g.ID)
	if err != nil {
		log.Fatalf("dw.GetGame err: %+v", err)
	}
	fmt.Printf("called dw.GetGame(%+v)\n", gg)

}
