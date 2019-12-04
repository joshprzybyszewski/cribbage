package main

import (
	"fmt"
	"os"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/persistence/mongodb"
)

func main() {
	mdb, err := mongodb.New(`mongodb://localhost:27017`)
	if err != nil {
		fmt.Printf("Error on New: %+v\n", err)
		os.Exit(1)
	}

	id := 49
	idStr := fmt.Sprintf("%d", id)
	otherIDStr := fmt.Sprintf("%d", id*2)
	g1ID := model.GameID(id)
	pID := model.PlayerID(idStr)
	josh := model.Player{
		ID:   pID,
		Name: idStr,
		Games: map[model.GameID]model.PlayerColor{
			g1ID:               model.Blue,
			model.GameID(5555): model.Red,
			model.GameID(9876): model.Green,
		},
	}
	ellen := model.Player{
		ID:   model.PlayerID(otherIDStr),
		Name: otherIDStr,
		Games: map[model.GameID]model.PlayerColor{
			g1ID:               model.Red,
			model.GameID(5555): model.Blue,
			model.GameID(9876): model.Red,
		},
	}
	fmt.Printf("Creating josh: %+v\n", josh)
	err = mdb.CreatePlayer(josh)
	if err != nil {
		fmt.Printf("Error on CreatePlayer: %+v\n", err)
		os.Exit(1)
	}

	p, err := mdb.GetPlayer(pID)
	if err != nil {
		fmt.Printf("Error on GetPlayer: %+v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Got Player: %+v\n", p)

	g := model.Game{
		ID:              g1ID,
		Players:         []model.Player{josh, ellen},
		Deck:            model.NewDeck(),
		BlockingPlayers: map[model.PlayerID]model.Blocker{ellen.ID: model.CountHand},
		CurrentDealer:   josh.ID,
		PlayerColors:    map[model.PlayerID]model.PlayerColor{josh.ID: model.Blue, ellen.ID: model.Red},
		CurrentScores:   map[model.PlayerColor]int{model.Blue: 0, model.Red: 0},
		LagScores:       map[model.PlayerColor]int{model.Blue: 0, model.Red: 0},
		Phase:           model.Counting,
		Hands: map[model.PlayerID][]model.Card{
			josh.ID: {
				model.NewCardFromString(`7s`),
				model.NewCardFromString(`8s`),
				model.NewCardFromString(`9s`),
				model.NewCardFromString(`10s`),
			},
			ellen.ID: {
				model.NewCardFromString(`7c`),
				model.NewCardFromString(`8c`),
				model.NewCardFromString(`9c`),
				model.NewCardFromString(`10c`),
			},
		},
		CutCard:     model.NewCardFromString(`7h`),
		Crib:        make([]model.Card, 4),
		PeggedCards: make([]model.PeggedCard, 0, 8),
	}

	fmt.Printf("Made Game: %+v\n", g)
	err = mdb.SaveGame(g)
	if err != nil {
		fmt.Printf("Error on SaveGame: %+v\n", err)
		os.Exit(1)
	}

	gDB, err := mdb.GetGame(g1ID)
	if err != nil {
		fmt.Printf("Error on GetGame: %+v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Got Game:  %+v\n", gDB)

}
