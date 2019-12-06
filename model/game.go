package model

func (g *Game) IsOver() bool {
	for _, score := range g.CurrentScores {
		if score >= WinningScore {
			return true
		}
	}
	return false
}

func (g *Game) NumActions() int {
	return len(g.Actions)
}

func (g *Game) AddAction(a PlayerAction) {
	g.Actions = append(g.Actions, a)
}

func (g *Game) CurrentPeg() int {
	if len(g.PeggedCards) == 0 {
		return 0
	}
	if g.goesAround() {
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

func (g *Game) goesAround() bool {
	if len(g.PeggedCards) == 0 {
		return false
	}

	lastPeggedCard := g.PeggedCards[len(g.PeggedCards)-1]
	lastPlayerWhoPlayed := lastPeggedCard.PlayerID
	for actIndex := g.NumActions() - 1; actIndex >= lastPeggedCard.Action; actIndex-- {
		act := g.Actions[actIndex]
		if pa, ok := act.Action.(PegAction); ok {
			if !pa.SayGo {
				// if anybody else has played a card, the goes have not gone around
				return false
			} else if act.ID == lastPlayerWhoPlayed {
				// if the last player who played has also said go, then the goes have gone around
				return true
			}
		}
	}

	return false
}
