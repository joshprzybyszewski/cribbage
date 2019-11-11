package model

func (g *Game) IsOver() bool {
	for _, score := range g.CurrentScores {
		if score >= WinningScore {
			return true
		}
	}
	return false
}
