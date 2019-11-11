package interaction

import (
	"github.com/joshprzybyszewski/cribbage/model"
)

type Player interface {
	NotifyBlocking(model.Blocker, interface{})
	NotifyMessage(interface{})
	NotifyScoreUpdate(CurrentScores, LagScores   map[model.PlayerColor]int, msgs ...string)
}