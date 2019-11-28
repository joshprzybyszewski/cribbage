package model

func (g *Game) IsOver() bool {
	for _, score := range g.CurrentScores {
		if score >= WinningScore {
			return true
		}
	}
	return false
}

func (g *Game) CurrentPeg() int {
	// keep in mind that this logic will need to change if/when we implement stealing points
	if len(g.PeggedCards) == 0 || g.NumActions >= g.PeggedCards[len(g.PeggedCards)-1].Action + len(g.Players) {
		return 0
	}
	cur := 0
	for _, pc := range g.PeggedCards {
		pv := pc.Card.PegValue()
		cur += pv
		if cur > MaxPeggingValue {
			cur = pv
		}
	}
	if cur == MaxPeggingValue {
		return 0
	}
	return cur
}
