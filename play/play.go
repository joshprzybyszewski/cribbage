package play

import (
	"fmt"

	survey "github.com/AlecAivazis/survey/v2"

	"github.com/joshprzybyszewski/cribbage/game"
)

func PlayGame() error {
	human := game.NewHumanPlayer(game.Blue)
	npc := getOpponent(game.Red)

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

func getOpponent(color game.PegColor) game.Player {
	const dumb = `dumb`
	const simple = `simple`
	const calculated = `calculated`
	opponentChoice := ``
	prompt := &survey.Select{
		Message: "Who would you like to play?",
		Options: []string{dumb, simple, calculated},
		Filter:  func(filter string, value string, index int) bool { return true },
	}
	survey.AskOne(prompt, &opponentChoice)

	switch opponentChoice {
	case dumb:
		return game.NewDumbNPC(color)
	case simple:
		return game.NewSimpleNPC(color)
	case calculated:
		return game.NewSimpleNPC(color)
	}

	return game.NewDumbNPC(color)
}
