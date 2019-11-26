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
	cur := 0
	for _, pc := range g.PeggedCards {
		pv := pc.Card.PegValue()
		cur += pv
		if cur > MaxPeggingValue {
			cur = pv
		}
	}
	return cur
}
