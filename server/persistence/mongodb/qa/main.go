package main

import (
	"fmt"
	"os"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
	"github.com/joshprzybyszewski/cribbage/server/persistence/mongodb"
	"github.com/joshprzybyszewski/cribbage/server/play"
)

func main() {
	mdb, err := mongodb.New(``)
	if err != nil {
		fmt.Printf("Error on New: %+v\n", err)
		os.Exit(1)
	}

	josh, ellen, pAPIs := createPlayers(1)
	fmt.Printf("Creating josh: %+v\n", josh)
	err = mdb.CreatePlayer(josh)
	if err != nil {
		fmt.Printf("Error on CreatePlayer: %+v\n", err)
		os.Exit(1)
	}

	g, err := play.CreateGame([]model.Player{josh, ellen}, pAPIs)
	if err != nil {
		fmt.Printf("Error on CreateGame: %+v\n", err)
		os.Exit(1)
	}

	josh.Games[g.ID] = g.PlayerColors[josh.ID]
	ellen.Games[g.ID] = g.PlayerColors[ellen.ID]

	p, err := mdb.GetPlayer(josh.ID)
	if err != nil {
		fmt.Printf("Error on GetPlayer: %+v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Got Player: %+v\n", p)

	fmt.Printf("Made Game: %+v\n", g)
	err = mdb.SaveGame(g)
	if err != nil {
		fmt.Printf("Error on SaveGame: %+v\n", err)
		os.Exit(1)
	}

	gDB, err := mdb.GetGame(g.ID)
	if err != nil {
		fmt.Printf("Error on GetGame: %+v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Got Game:  %+v\n", gDB)

}

func createPlayers(id int) (josh, ellen model.Player, pAPIs map[model.PlayerID]interaction.Player) {
	idStr := fmt.Sprintf("%d", id)
	otherIDStr := fmt.Sprintf("%d", id*2)
	josh = model.Player{
		ID:   model.PlayerID(idStr),
		Name: idStr,
		Games: map[model.GameID]model.PlayerColor{
			model.GameID(5555): model.Red,
			model.GameID(9876): model.Green,
		},
	}
	ellen = model.Player{
		ID:   model.PlayerID(otherIDStr),
		Name: otherIDStr,
		Games: map[model.GameID]model.PlayerColor{
			model.GameID(5555): model.Blue,
			model.GameID(9876): model.Red,
		},
	}
	pAPIs = map[model.PlayerID]interaction.Player{
		josh.ID:  &interaction.Empty{PID: josh.ID},
		ellen.ID: &interaction.Empty{PID: ellen.ID},
	}
	return josh, ellen, pAPIs
}
