package model

func (g *Game) IsOver() bool {
	for _, score := range g.CurrentScores {
		if score >= WinningScore {
			return true
		}
	}
	return false
}

func (g *Game) AddPoints(pc PlayerColor, p int, msgs ...string) {
	g.LagScores[pc] = g.CurrentScores[pc]
	g.CurrentScores[pc] = g.CurrentScores[pc] + p
	/* TODO tell all of the players about score updates
	for _, p := range g.Players {
		p.TellAboutScores(g.CurrentScores, g.LagScores, msgs...)
	}
	*/
}
