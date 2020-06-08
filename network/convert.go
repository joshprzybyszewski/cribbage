package network

import "github.com/joshprzybyszewski/cribbage/model"

func convertToColors(cs map[model.PlayerID]model.PlayerColor) map[model.PlayerID]string {
	colors := make(map[model.PlayerID]string, len(cs))
	for p, c := range cs {
		colors[p] = c.String()
	}
	return colors
}

func convertToScores(mCurrentScores, mLagScores map[model.PlayerColor]int) (map[string]int, map[string]int) {
	current := make(map[string]int, len(mCurrentScores))
	for c, s := range mCurrentScores {
		current[c.String()] = s
	}
	lag := make(map[string]int, len(mLagScores))
	for c, s := range mLagScores {
		lag[c.String()] = s
	}
	return current, lag
}

func convertToBlockingPlayers(bs map[model.PlayerID]model.Blocker) map[model.PlayerID]string {
	blockers := make(map[model.PlayerID]string, len(bs))
	for p, b := range bs {
		blockers[p] = b.String()
	}
	return blockers
}
