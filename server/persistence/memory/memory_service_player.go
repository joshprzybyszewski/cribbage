package memory

import (
	"errors"
	"sync"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

var pservice *playerService

func savePlayerService(newPS *playerService) error {
	pservice.lock.Lock()
	defer pservice.lock.Unlock()

	newPS.lock.Lock()
	defer newPS.lock.Unlock()

	for pid, p := range newPS.players {
		pservice.players[pid] = p
	}

	return nil
}

var _ persistence.PlayerService = (*playerService)(nil)

type playerService struct {
	lock sync.Mutex

	players map[model.PlayerID]model.Player
}

func getPlayerService() persistence.PlayerService {
	if pservice == nil {
		pservice = &playerService{
			players: map[model.PlayerID]model.Player{},
		}
	}
	return pservice
}

func (ps *playerService) Copy() *playerService {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	cpy := make(map[model.PlayerID]model.Player, len(ps.players))
	for id, player := range ps.players {
		pCpy := player
		cpy[id] = pCpy
	}

	return &playerService{
		players: cpy,
	}
}

func (ps *playerService) Get(id model.PlayerID) (model.Player, error) {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	if i, ok := ps.players[id]; ok {
		return i, nil
	}

	return model.Player{}, persistence.ErrPlayerNotFound
}

func (ps *playerService) Create(p model.Player) error {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	id := p.ID
	if _, ok := ps.players[id]; ok {
		return persistence.ErrPlayerAlreadyExists
	}

	ps.players[id] = p
	return nil
}

func (ps *playerService) UpdateGameColor(pID model.PlayerID, gID model.GameID, color model.PlayerColor) error {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	// Assign color to player
	if _, ok := ps.players[pID]; !ok {
		return persistence.ErrPlayerNotFound
	}
	pCopy := ps.players[pID]

	if pCopy.Games == nil {
		pCopy.Games = map[model.GameID]model.PlayerColor{}
	}

	if c, ok := pCopy.Games[gID]; !ok {
		pCopy.Games[gID] = color
		ps.players[pID] = pCopy
	} else if c != color {
		return errors.New(`mismatched player colors`)
	}
	return nil
}
