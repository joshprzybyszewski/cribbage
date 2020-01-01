package interaction

import (
	"errors"

	"github.com/joshprzybyszewski/cribbage/model"
)

type Player interface {
	ID() model.PlayerID

	NotifyBlocking(model.Blocker, model.Game, string) error
	NotifyMessage(model.Game, string) error
	NotifyScoreUpdate(g model.Game, msgs ...string) error
}

func New(pID model.PlayerID, m Means) PlayerMeans {
	return PlayerMeans{
		PlayerID:      pID,
		PreferredMode: m.Mode,
		Interactions:  []Means{m},
	}
}

func FromPlayerMeans(pm PlayerMeans) (Player, error) {
	pID := pm.PlayerID
	means := pm.getMeans(pm.PreferredMode)

	switch means.Mode {
	case Localhost:
		return newLocalhostPlayer(pID, means.Info), nil
	case NPC:
		ah, ok := means.Info.(ActionHandler)
		if !ok {
			return nil, errors.New(`player means info should contain an action handler, but it doesn't`)
		}
		return NewNPCPlayer(pID, ah)
	}

	return nil, errors.New(`mode unsupported`)
}

type ActionHandler interface {
	Handle(action model.PlayerAction) error
}
