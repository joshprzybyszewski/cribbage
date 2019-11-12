package play

import (
	"errors"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
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

func isWaitingForPlayer(g *model.Game, action model.PlayerAction) error {
	isWaitingForThisPlayer := false
	for bpID := range g.BlockingPlayers {
		if action.ID == bpID {
			isWaitingForThisPlayer = true
		}
	}
	if !isWaitingForThisPlayer {
		return errors.New(`Game is not blocked by this player`)
	}

	return nil
}

func addPlayerToBlocker(g *model.Game, pID model.PlayerID, reason model.Blocker, pAPIs map[model.PlayerID]interaction.Player, msgs ...interface{}) {
	if br, ok := g.BlockingPlayers[pID]; ok && br != reason {
			log.Printf("Same player (%s) blocking for new reason (%v vs. %v)", pID, br, reason)
	}
	g.BlockingPlayers[pID] = reason
	pAPI := pAPIs[pID]
	pAPI.NotifyBlocking(reason, msgs...)
}

func removePlayerFromBlockers(g *model.Game, action model.PlayerAction) {
	if br, ok := g.BlockingPlayers[action.ID]; ok && br == action.Overcomes {
		delete(g.BlockingPlayers, action.ID)
	} else if !ok {
		log.Printf(`Did not find player "%+v" in list of blocking players (%+v)`, action.ID, g.BlockingPlayers)
	} else {
		log.Printf(`Player was not blocked by "%v", but rather "%v"`, action.Overcomes, br)
	}
}

func roundCutter(g *model.Game) model.PlayerID {
	pIDs := playersToDealTo(g)
	return pIDs[len(pIDs)-2]
}

func addPoints(g *model.Game, pID model.PlayerID, pts int, pAPIs map[model.PlayerID]interaction.Player, msgs ...string) {
	if pts == 0 {
		return
	} else if pts < 0 {
		log.Printf(`Attempted to score %d points for player %v`, pts, pID)
		return
	}

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

func handContains(hand []model.Card, c model.Card) bool {
	for _, hc := range hand {
		if hc == c {
			return true
		}
	}
	return false
}

func hasBeenPegged(pegged []model.PeggedCard, c model.Card) bool {
	for _, pc := range pegged {
		if pc.Card == c {
			return true
		}
	}
	return true
}

func minUnpeggedValue(hand []model.Card, pegged []model.PeggedCard) int {
	peggedMap := make(map[model.Card]struct{}, len(sub))
	for _, pc := range pegged {
		peggedMap[pc.Card] = struct{}{}
	}

	min := model.MaxPeggingValue
	for _, hc := range hand {
		if _, ok := peggedMap[hc]; ok {
			continue
		}
		if pv := hc.PegValue(); pv < min {
			min = pv
		}
	}
	return min
}