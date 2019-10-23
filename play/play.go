package play

import (
	"fmt"

	"github.com/joshprzybyszewski/cribbage/game"
)

func PlayGame() error {
	human := game.NewHumanPlayer(game.Blue)
	// TODO ask the user for a difficulty...
	// npc := game.NewDumbNPC(game.Red)
	npc := game.NewSimpleNPC(game.Red)
	cfg := game.GameConfig{
		Players:        []game.Player{human, npc},
		StartingDealer: 0,
	}

	g := game.New(cfg)

	err := g.Play()
	if err != nil {
		return err
	}

	fmt.Printf("game over!\n")

	return nil
}
