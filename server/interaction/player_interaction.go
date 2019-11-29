package interaction

import (
	"github.com/joshprzybyszewski/cribbage/model"
)

type Player interface {
	ID() model.PlayerID

	NotifyBlocking(model.Blocker, model.Game, string) error
	NotifyMessage(model.Game, string) error
	NotifyScoreUpdate(g model.Game, msgs ...string) error
}

func New(pID model.PlayerID, im model.InteractionMeans) Player {
	switch im.Means {
	case `localhost`:
		return newLocalhostPlayer(pID, im.Info)
	}

	return nil
}
