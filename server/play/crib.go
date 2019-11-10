package play

import (
	"github.com/joshprzybyszewski/cribbage/model"
)

func buildCrib(g *model.Game, pAPIs map[model.PlayerID]interaction.Player) error {
	pIDs := playersToDealTo(g)
	desired := numDesiredCribCards(g)

	for _, pId := range pIDs {
		pAPIs[pID].NotifyBlocking(model.CribCard, model.CribBlocker{
			Desired: desired,
			Dealer: g.CurrentDealer,
			PlayerColors: g.PlayerColors,
		})

		go func(pcopy Player) {
			
			err = g.round.AcceptCribCards(pcopy.AddToCrib(g.Dealer().Color(), desired)...)
		}(p)
	}

	return err
}

func handleCribBuild(g *model.Game, buildCribAction PlayerAction, pAPIs map[model.PlayerID]interaction.Player) error {
	if buildCribAction.Overcomes != model.CribCard {
		return errors.New(`Does not attempt to build crib`)
	}
	if err := isWaitingForPlayer(g, buildCribAction); err != nil {
		return err
	}

	bca, ok := buildCribAction.(model.BuildCribAction)
	if !ok {
		return errors.New(`tried building crib with a different action`)
	}

	if len(bca.Cards) != numDesiredCribCards(g) {
		pAPIs[buildCribAction.ID].NotifyBlocking(model.CribCard, `Need to submit all required cards at once`)
		return nil
	}

	removePlayerFromBlockers(g, buildCribAction)

	// TODO start acting on this now
}

func numDesiredCribCards(g *model.Game) int {
	if len(g.Players) > 2 {
		return 1
	}
	return 2
}