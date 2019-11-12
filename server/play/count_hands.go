package play

import (
	"errors"

	"github.com/joshprzybyszewski/cribbage/server/interaction"
	"github.com/joshprzybyszewski/cribbage/logic/scorer"
	"github.com/joshprzybyszewski/cribbage/model"
)

var _ PhaseHandler = (*handCountingHandler)(nil)
type handCountingHandler struct {}

func (*handCountingHandler) Start(g *model.Game, pAPIs map[model.PlayerID]interaction.Player) error {
	pIDs := playersToDealTo(g)
	firstPlayerID := pIDs[0]
	addPlayerToBlocker(g, firstPlayerID, model.CountHand, pAPIs)

	return nil
}

func (*handCountingHandler) HandleAction(g *model.Game, action model.PlayerAction, pAPIs map[model.PlayerID]interaction.Player) error {
	if err := validateAction(g, action, model.CountHand); err != nil {
		return err
	}

	cha, ok := action.Action.(model.CountHandAction)
	if !ok {
		return errors.New(`tried counting hand with a different action`)
	}

	pID := action.ID
	hand := g.Hands[pID]
	leadCard := g.CutCard
	pts := scorer.HandPoints(leadCard, hand)

	if cha.Pts != pts {
		addPlayerToBlocker(g, pID, model.CountHand, pAPIs, `you did not submit the correct number of points for your hand`)
		return nil
	}

	addPoints(g, pID, pts, pAPIs, `hand (`+leadCard.String()+`: `+handString(hand) +`)`)

	if g.IsOver() {
		return nil
	}

	pIDs := playersToDealTo(g)
	nextScorerIndex := len(pIDs) // invalid index
	for i, id := range pIDs {
		if id == pID {
			nextScorerIndex = i +1
			break
		}
	}
	
	if nextScorerIndex <= len(pIDs)-1 {
		nextID := pIDs[nextScorerIndex]
		addPlayerToBlocker(g, nextID, model.CountHand, pAPIs)
	} else {
		g.Phase++
	}

	return nil
}
