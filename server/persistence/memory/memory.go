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
	lock sync.Mutex

	games        map[model.GameID][]model.Game
	players      map[model.PlayerID]model.Player
	interactions map[model.PlayerID]interaction.Player
}

func New() persistence.DB {
	return &memory{}
}

func (m *memory) GetGame(id model.GameID) (model.Game, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if games, ok := m.games[id]; ok {
		return games[len(games)-1], nil
	}
	return model.Game{}, errors.New(`does not have player`)
}

func (m *memory) GetPlayer(id model.PlayerID) (model.Player, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if p, ok := m.players[id]; ok {
		return p, nil
	}
	return model.Player{}, errors.New(`does not have player`)
}

func (m *memory) GetInteraction(id model.PlayerID) (interaction.Player, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if i, ok := m.interactions[id]; ok {
		return i, nil
	}
	return nil, errors.New(`does not have player`)
}

func (m *memory) SaveGame(g model.Game) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	id := g.ID

	if len(m.games[id]) != g.NumActions {
		return errors.New(`game does not know about recent actions`)
	}

	m.games[id] = append(m.games[id], g)

	return nil
}

func (m *memory) CreatePlayer(p model.Player) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	id := p.ID
	if _, ok := m.players[id]; ok {
		return errors.New(`player already exists`)
	}

	m.players[id] = p
	return nil
}

func (m *memory) AddPlayerColorToGame(id model.PlayerID, color model.PlayerColor, gID model.GameID) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.players[id]; !ok {
		return errors.New(`player does not exist`)
	}

	m.players[id].Games[gID] = color
	return nil

}

func (m *memory) SaveInteraction(i interaction.Player) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.interactions[i.ID()] = i
	return nil
}
