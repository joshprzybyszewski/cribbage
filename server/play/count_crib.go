package play

import (
	"errors"

	"github.com/joshprzybyszewski/cribbage/logic/scorer"
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
)

var _ PhaseHandler = (*cribCountingHandler)(nil)

type cribCountingHandler struct{}

func (*cribCountingHandler) Start(g *model.Game, pAPIs map[model.PlayerID]interaction.Player) error {
	addPlayerToBlocker(g, g.CurrentDealer, model.CountCrib, pAPIs, ``)

	return nil
}

func (*cribCountingHandler) HandleAction(
	g *model.Game,
	action model.PlayerAction,
	pAPIs map[model.PlayerID]interaction.Player,
) error {

	if err := validateAction(g, action, model.CountCrib); err != nil {
		return err
	}

	cca, ok := action.Action.(model.CountCribAction)
	if !ok {
		return errors.New(`tried counting crib with a different action`)
	}

	pID := action.ID
	if pID != g.CurrentDealer {
		return errors.New(`not the dealer tried counting the crib`)
	}

	crib := g.Crib
	leadCard := g.CutCard
	pts := scorer.CribPoints(leadCard, crib)

	if cca.Pts == 19 {
		cca.Pts = 0
	}

	if cca.Pts != pts {
		addPlayerToBlocker(g, pID, model.CountCrib, pAPIs, `you did not submit the correct number of points for the crib`)
		return errors.New(`wrong number of points`)
	}

	addPoints(g, pID, pts, pAPIs, `crib (`+leadCard.String()+`: `+handString(crib)+`)`)

	if g.IsOver() {
		return nil
	}
	removePlayerFromBlockers(g, action)

	// Move forward to dealing
	pIDs := playersToDealTo(g)
	for i, id := range pIDs {
		if id == pID {
			g.CurrentDealer = pIDs[(i+1)%len(pIDs)]
			break
		}
	}

	return nil
}
