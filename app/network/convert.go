package network

import (
	"sort"

	"github.com/joshprzybyszewski/cribbage/model"
)

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
	teams := make([]GetGameResponseTeam, 0, 3)
	for i, p := range g.Players {
		ps := make([]Player, 0, 2)
		ps = append(ps, convertToPlayer(p))
		if len(g.Players) == 4 {
			if i > 1 {
				break
			}
			ps = append(ps, convertToPlayer(g.Players[i+2]))
		}

		color, ok := p.Games[g.ID]
		if !ok {
			color = model.UnsetColor
		}

		colorStr := color.String()
		if color == model.UnsetColor {
			colorStr = ``
		}

		t := GetGameResponseTeam{
			Color:        colorStr,
			CurrentScore: g.CurrentScores[color],
			LagScore:     g.LagScores[color],
			Players:      ps,
		}

		teams = append(teams, t)
	}

	sort.Slice(teams, func(i, j int) bool {
		return teams[i].Color == `` && teams[j].Color != `` && teams[i].Color > teams[j].Color
	})
	return teams
}
