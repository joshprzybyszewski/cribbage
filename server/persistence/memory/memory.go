package memory

import (
	"errors"
	"sync"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

var _ persistence.DB = (*memory)(nil)

type memory struct {
	games        map[model.GameID]model.Game
	gameLocks    map[model.GameID]*sync.Mutex

	players      map[model.PlayerID]model.Player
	playerLocks  map[model.PlayerID]*sync.Mutex

	interactions map[model.PlayerID]interaction.Player
	interactionLocks map[model.PlayerID]*sync.Mutex
}

func New() persistence.DB {
	return &memory{
		
	}
}

func (m *memory) GetGame(id model.GameID) (model.Game, error) {
	gl, ok := m.gameLocks[id]
	if !ok {
		m.gameLocks[id] = &sync.Mutex{}
	}
	gl.Lock()
	if g, ok := m.games[id]; ok {
		return g, nil
	}
	return model.Game{}, errors.New(`does not have player`)
}

func (m *memory) GetPlayer(id model.PlayerID) (model.Player, error) {
	pl, ok := m.playerLocks[id]
	if !ok {
		m.playerLocks[id] = &sync.Mutex{}
	}
	pl.Lock()
	if p, ok := m.players[id]; ok {
		return p, nil
	}
	return model.Player{}, errors.New(`does not have player`)
}

func (m *memory) GetInteraction(id model.PlayerID) (interaction.Player, error) {
	il, ok := m.interactionLocks[id]
	if !ok {
		m.interactionLocks[id] = &sync.Mutex{}
	}
	il.Lock()
	if i, ok := m.interactions[id]; ok {
		return i, nil
	}
	return nil, errors.New(`does not have player`)
}

func (m *memory) SaveGame(g model.Game) (error) {
	m.games[g.ID] = g
	m.gameLocks[g.ID].Unlock()
	return nil
}

func (m *memory) SavePlayer(p model.Player) (error) {
	m.players[p.ID] = p
	m.playerLocks[p.ID].Unlock()
	return nil
}

func (m *memory) SaveInteraction(i interaction.Player) (error) {
	m.interactions[i.ID()] = i
	m.interactionLocks[i.ID()].Unlock()
	return nil
}
