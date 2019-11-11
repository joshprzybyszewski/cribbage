package play

import (
	"errors"

	"github.com/joshprzybyszewski/cribbage/model"
)

func playersToDealTo(g *model.Game) []model.PlayerID {
	pIDs := make([]PlayerID, len(g.Players))
	dealerIndex := 0
	for i, p := range g.Players {
		pIDs[i] = p.ID
		if p.ID == g.CurrentDealer {
			dealerIndex = i
		}
	}

	if dealerIndex == len(pIDs)-1 {
		// Return the slice of player IDs if the dealer is last in line
		return pIDs
	}

	// put the dealer last in line for cards
	return append(pIDs[dealerIndex+1:], pIDs[:dealerIndex+1]...)
}

func isWaitingForPlayer(g *model.Game, action PlayerAction) error {
	isWaitingForThisPlayer := false
	for _, bp := range g.BlockingPlayers {
		if action.ID == bp.ID {
			isWaitingForThisPlayer = true
		}
	}
	if !isWaitingForThisPlayer {
		return errors.New(`Game is not blocked by this player`)
	}

	return nil
}

func removePlayerFromBlockers(g *model.Game, action PlayerAction) {
	if len(g.BlockingPlayers) == 1 && g.BlockingPlayers[0].ID == action.ID {
		g.BlockingPlayers = g.BlockingPlayers[:0]
		return
	} else {
		blockIndex := -1
		for i, bp := range g.BlockingPlayers {
			if bp.ID == action.ID {
				blockIndex = i
				break
			}
		}
		if blockIndex == len(g.BlockingPlayers) - 1 {
			g.BlockingPlayers = g.BlockingPlayers[:i]
		} else if blockIndex >= 0{
			g.BlockingPlayers = append(g.BlockingPlayers[:i], g.BlockingPlayers[i+1:]...)
		} else {
			log.Errorf(`Did not find player "%+v" in list of blocking players (%+v)`, action.ID, g.BlockingPlayers)
		}
	}
}

func roundCutter(g *model.Game) model.PlayerID {
	pIDs := playersToDealTo(g)
	return pIDs[len(pIDs)-2]
}

func addPoints(g *model.Game, pID model.PlayerID, pts int, pAPIs map[model.PlayerID]interaction.Player, msgs ...string) {
	pc := g.PlayerColors[pID]
	g.LagScores[pc] = g.CurrentScores[pc]
	g.CurrentScores[pc] = g.CurrentScores[pc] + pts
	
	for _, pAPI := range pAPIs {
		pAPI.NotifyScoreUpdate(g.CurrentScores, g.LagScores, msgs...)
	}
}

// isSuperSet returns true if all of the cards in sub exist in super
func isSuperSet(super, sub []model.Card) bool {
	superMap := make(map[model.Card]struct{}, len(super))
	for _, c := range super {
		superMap[c] = struct{}{}
	}
	for _, c := range sub {
		if _, ok := superMap[c]; !ok {
			return false
		}
	}
	return true
}

// removeSubset returns a new slice which is made from super and does not have
// and cards from sub in it. It does not check if sub is a subset of super, but
// assumes you have already done so
func removeSubset(super, sub []model.Card) []model.Card {
	subMap := make(map[model.Card]struct{}, len(sub))
	for _, c := range sub {
		subMap[c] = struct{}{}
	}

	ret := make([]model.Card,0,len(super))
	for _, c := range super {
		if _, ok := subMap[c]; ok {
			continue
		}
		ret = append(ret, c)
	}
	return ret
}