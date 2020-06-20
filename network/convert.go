package network

import "github.com/joshprzybyszewski/cribbage/model"

func convertToColor(c model.PlayerColor) string {
	return c.String()
}

func convertFromColor(c string) model.PlayerColor {
	return model.NewPlayerColorFromString(c)
}

func convertToPhase(p model.Phase) string {
	return p.String()
}

func convertFromPhase(p string) model.Phase {
	return model.NewPhaseFromString(p)
}

func convertToBlocker(b model.Blocker) string {
	return b.String()
}

func convertFromBlocker(b string) model.Blocker {
	return model.NewBlockerFromString(b)
}

func convertToColors(cs map[model.PlayerID]model.PlayerColor) map[model.PlayerID]string {
	colors := make(map[model.PlayerID]string, len(cs))
	for p, c := range cs {
		colors[p] = convertToColor(c)
	}
	return colors
}

func convertFromColors(cs map[model.PlayerID]string) map[model.PlayerID]model.PlayerColor {
	colors := make(map[model.PlayerID]model.PlayerColor, len(cs))
	for p, c := range cs {
		colors[p] = convertFromColor(c)
	}
	return colors
}

func convertToScores(mCurrentScores, mLagScores map[model.PlayerColor]int) (map[string]int, map[string]int) {
	current := make(map[string]int, len(mCurrentScores))
	for c, s := range mCurrentScores {
		current[convertToColor(c)] = s
	}
	lag := make(map[string]int, len(mLagScores))
	for c, s := range mLagScores {
		lag[convertToColor(c)] = s
	}
	return current, lag
}

func convertFromScores(nCurScores, nLagScores map[string]int) (mCur, mLag map[model.PlayerColor]int) {
	mCur = make(map[model.PlayerColor]int, len(nCurScores))
	for c, s := range nCurScores {
		mCur[convertFromColor(c)] = s
	}
	mLag = make(map[model.PlayerColor]int, len(nLagScores))
	for c, s := range nLagScores {
		mLag[convertFromColor(c)] = s
	}
	return mCur, mLag
}

func convertToBlockingPlayers(bs map[model.PlayerID]model.Blocker) map[model.PlayerID]string {
	blockers := make(map[model.PlayerID]string, len(bs))
	for p, b := range bs {
		blockers[p] = convertToBlocker(b)
	}
	return blockers
}

func convertFromBlockingPlayers(bs map[model.PlayerID]string) map[model.PlayerID]model.Blocker {
	blockers := make(map[model.PlayerID]model.Blocker, len(bs))
	for p, b := range bs {
		blockers[p] = convertFromBlocker(b)
	}
	return blockers
}
