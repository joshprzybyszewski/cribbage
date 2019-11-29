package play

import (
	"errors"
	"log"

	"github.com/joshprzybyszewski/cribbage/logic/pegging"
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
)

var _ PhaseHandler = (*peggingHandler)(nil)

type peggingHandler struct{}

func (*peggingHandler) Start(g *model.Game, pAPIs map[model.PlayerID]interaction.Player) error {
	g.PeggedCards = g.PeggedCards[:0]

	// put the player after the dealer as the blocking player
	pIDs := playersToDealTo(g)
	pID := pIDs[0]
	addPlayerToBlocker(g, pID, model.PegCard, pAPIs, `please peg a card`)

	return nil
}

func (*peggingHandler) HandleAction(g *model.Game,
	action model.PlayerAction,
	pAPIs map[model.PlayerID]interaction.Player,
) error {

	// VALIDATE: check the action, then check the peg/go
	if err := validateAction(g, action, model.PegCard); err != nil {
		return err
	}

	pa, ok := action.Action.(model.PegAction)
	if !ok {
		return errors.New(`tried pegging with a different action`)
	}

	pID := action.ID
	if err := validatePegAction(g, pID, pa); err != nil {
		addPlayerToBlocker(g, pID, model.PegCard, pAPIs, err.Error())
		return nil
	}

	// CLEAN: remove this player from the blockers
	if len(g.BlockingPlayers) != 1 {
		log.Printf("Expected one blocker for pegging, but had: %+v\n", g.BlockingPlayers)
	}
	removePlayerFromBlockers(g, action)

	// ACT: do the "say go" or peg
	if pa.SayGo {
		doSayGo(g, action, pAPIs)
	} else if err := doPeg(g, action, pa, pAPIs); err != nil {
		return err
	}

	// PROGRESS: move the game state forward appropriately
	progressAfterPeg(g, action, pAPIs)

	return nil
}

func validatePegAction(g *model.Game, pID model.PlayerID, pa model.PegAction) error {
	if pa.SayGo {
		curPeg := g.CurrentPeg()
		minCardVal := minUnpeggedValue(g.Hands[pID], g.PeggedCards)
		if curPeg+minCardVal <= model.MaxPeggingValue {
			return errors.New(`Cannot say go when has unpegged playable card`)
		}
		return nil
	}

	if !handContains(g.Hands[pID], pa.Card) {
		return errors.New(`Cannot peg card you don't have`)
	}

	if hasBeenPegged(g.PeggedCards, pa.Card) {
		return errors.New(`Cannot peg same card twice`)
	}

	if g.CurrentPeg()+pa.Card.PegValue() > model.MaxPeggingValue {
		return errors.New(`Cannot peg card with this value`)
	}

	return nil
}

func doPeg(
	g *model.Game,
	action model.PlayerAction,
	pa model.PegAction,
	pAPIs map[model.PlayerID]interaction.Player,
) error {

	pts, err := pegging.PointsForCard(g.PeggedCards, pa.Card)
	if err != nil {
		return err
	}

	addPoints(g, action.ID, pts, pAPIs, `pegging`)

	g.PeggedCards = append(g.PeggedCards, model.PeggedCard{
		Card:     pa.Card,
		PlayerID: action.ID,
		Action:   g.NumActions() + 1,
	})

	return nil
}

func doSayGo(g *model.Game,
	action model.PlayerAction,
	pAPIs map[model.PlayerID]interaction.Player,
) {

	if len(g.PeggedCards) == 0 {
		return
	}

	lastPeggerID := g.PeggedCards[len(g.PeggedCards)-1].PlayerID
	if lastPeggerID == action.ID {
		// The go's went all the way around. Take a point
		addPoints(g, action.ID, 1, pAPIs, `the go`)
	}
}

func progressAfterPeg(
	g *model.Game,
	action model.PlayerAction,
	pAPIs map[model.PlayerID]interaction.Player,
) {

	if g.IsOver() {
		// we shouldn't do anything if the game is over
		return
	}

	if len(g.PeggedCards) == 4*len(g.Players) {
		// This was the last card: give one point to this player.
		addPoints(g, action.ID, 1, pAPIs, `last card`)
		return
	}

	// Set the next player to peg as the blocker
	// If we don't require everyone to say go, then we'll need to change the logic
	// in game.CurrentPeg
	nextPlayerIndex := -1
	for i, p := range g.Players {
		if p.ID == action.ID {
			nextPlayerIndex = (i + 1) % len(g.Players)
			break
		}
	}

	bp := g.Players[nextPlayerIndex]
	addPlayerToBlocker(g, bp.ID, model.PegCard, pAPIs, ``)
}
