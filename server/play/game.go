package play

import (
	"errors"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
)

var (
handlers = map[model.Phase]PhaseHandler{
	model.Deal: &dealingHandler{},
	model.BuildCribReady: &cribBuildingHandler{},
	model.BuildCrib: &cribBuildingHandler{},
	model.CutReady: &cuttingHandler{},
	model.Cut: &cuttingHandler{},
	model.PeggingReady: &peggingHandler{},
	model.Pegging: &peggingHandler{},
	model.CountingReady: &handCountingHandler{},
	model.Counting: &handCountingHandler{},
	model.CribCountingReady: &cribCountingHandler{},
	model.CribCounting: &cribCountingHandler{},
	model.DealingReady: &dealingHandler{},
	}
)

func HandleAction(g *model.Game, action model.PlayerAction, pAPIs map[model.PlayerID]interaction.Player) (error) {
	if g.ID != action.GameID {
		return errors.New(`action not for this game`)
	}

	switch p := g.Phase; p {
	case model.Deal,
	model.BuildCrib,
	model.Cut,
	model.Pegging,
	model.Counting,
	model.CribCounting:
		err := handlers[p].HandleAction(g, action, pAPIs)
		if err != nil {
			return err
		}
	}
	

	switch p := g.Phase; p {
	case model.BuildCribReady,
	model.CutReady,
	model.PeggingReady,
	model.CountingReady,
	model.CribCountingReady,
	model.DealingReady:
		err := handlers[p].Start(g, pAPIs)
		if err != nil {
			return err
		}
		g.Phase++
		if g.Phase > model.DealingReady {
			g.Phase = model.Deal
		}
	}
	

	return  nil
}
