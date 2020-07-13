package network

import "github.com/joshprzybyszewski/cribbage/model"

func convertToColor(c model.PlayerColor) string {
	return c.String()
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

func convertFromScores(teams []GetGameResponseTeam) (mCur, mLag map[model.PlayerColor]int) {
	mCur = make(map[model.PlayerColor]int, len(teams))
	mLag = make(map[model.PlayerColor]int, len(teams))
	for _, t := range teams {
		clr := model.NewPlayerColorFromString(t.Color)
		mCur[clr] = t.CurrentScore
		mLag[clr] = t.LagScore
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

func convertToTeams(g model.Game) []GetGameResponseTeam {
	teams := make([]GetGameResponseTeam, 0, len(g.CurrentScores))
	for c, s := range g.CurrentScores {
		t := GetGameResponseTeam{
			Color:        c.String(),
			CurrentScore: s,
			LagScore:     g.LagScores[c],
		}
		// a team will only ever have 1 or 2 players on it
		ps := make([]Player, 0, 2)
		for _, p := range g.Players {
			if g.PlayerColors[p.ID] == c {
				ps = append(ps, convertToPlayer(p))
			}
		}
		t.Players = ps
		teams = append(teams, t)
	}
	return teams
}
