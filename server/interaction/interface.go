package interaction

import (
	"errors"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server"
	"github.com/joshprzybyszewski/cribbage/server/interaction/npc"
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
		return npc.NewNPCPlayer(pID, server.HandleAction)
	}

	return nil, errors.New(`mode unsupported`)
}
