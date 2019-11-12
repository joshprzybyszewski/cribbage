package memory

import (
	"errors"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
)

type Memory struct {
	games        map[model.GameID]model.Game
	players      map[model.PlayerID]model.Player
	interactions map[model.PlayerID]interaction.Player
}

func (m *Memory) GetGame(id model.GameID) (model.Game, error) {
	if g, ok := m.games[id]; ok {
		return g, nil
	}
	return model.Game{}, errors.New(`does not have player`)
}

func (m *Memory) GetPlayer(id model.PlayerID) (model.Player, error) {
	if p, ok := m.players[id]; ok {
		return p, nil
	}
	return model.Player{}, errors.New(`does not have player`)
}

func (m *Memory) GetInteraction(id model.PlayerID) (interaction.Player, error) {
	if i, ok := m.interactions[id]; ok {
		return i, nil
	}
	return nil, errors.New(`does not have player`)
}
