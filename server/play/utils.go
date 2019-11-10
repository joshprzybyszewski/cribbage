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