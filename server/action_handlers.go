package server

import (
	"context"
	"sync"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
)

var _ interaction.ActionHandler = (*npcActionHandler)(nil)

type npcActionHandler struct {
	lock sync.Mutex
}

func (ah *npcActionHandler) Handle(pa model.PlayerAction) error {
	ah.lock.Lock()
	defer ah.lock.Unlock()

	return HandleAction(context.Background(), pa)
}
