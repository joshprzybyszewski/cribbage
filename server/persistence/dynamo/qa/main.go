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

	pID := model.PlayerID(`joshysquashy`)
	p := model.Player{
		ID:   pID,
		Name: `jesus is king`,
		Games: map[model.GameID]model.PlayerColor{
			model.GameID(123): model.Blue,
			model.GameID(456): model.Red,
			model.GameID(789): model.Green,
		},
	}
	fmt.Printf("calling dw.CreatePlayer(%+v)\n", p)
	err = dw.CreatePlayer(p)
	if err != nil {
		log.Fatalf("dw.CreatePlayer err: %+v", err)
	}

	fmt.Printf("calling dw.GetPlayer(%+v)\n", pID)
	p2, err := dw.GetPlayer(pID)
	if err != nil {
		log.Fatalf("dw.GetPlayer err: %+v", err)
	}
	fmt.Printf("player p2 := %+v\n", p2)

}
