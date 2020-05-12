package persistence

import (
	"errors"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
)

type DB interface {
	CreatePlayer(p model.Player) error
	GetPlayer(id model.PlayerID) (model.Player, error)
	AddPlayerColorToGame(id model.PlayerID, color model.PlayerColor, gID model.GameID) error

	CreateGame(g model.Game) error
	GetGame(id model.GameID) (model.Game, error)
	GetGameAction(id model.GameID, numActions uint) (model.Game, error)
	SaveGame(g model.Game) error

	GetInteraction(id model.PlayerID) (interaction.PlayerMeans, error)
	SaveInteraction(pm interaction.PlayerMeans) error
}

type db struct {
	games        GameService
	players      PlayerService
	interactions InteractionService
}

func New(gs GameService, ps PlayerService, is InteractionService) DB {
	return &db{
		games:        gs,
		players:      ps,
		interactions: is,
	}
}

func (d *db) CreatePlayer(p model.Player) error {
	return d.players.Create(p)
}

func (d *db) GetPlayer(id model.PlayerID) (model.Player, error) {
	return d.players.Get(id)
}

func (d *db) AddPlayerColorToGame(pID model.PlayerID, color model.PlayerColor, gID model.GameID) error {
	err := d.games.UpdatePlayerColor(gID, pID, color)
	if err != nil {
		return err
	}
	return d.players.UpdateGameColor(pID, gID, color)
}

func (d *db) GetGame(id model.GameID) (model.Game, error) {
	g, err := d.games.Get(id)
	if err != nil {
		return model.Game{}, err
	}

	for i, player := range g.Players {
		// overwrite the player that the game service knows
		// about with the player that the players service knows about
		p, err := d.players.Get(player.ID)
		if err != nil {
			return model.Game{}, err
		}
		g.Players[i] = p
	}

	return g, nil
}

func (d *db) GetGameAction(id model.GameID, numActions uint) (model.Game, error) {
	g, err := d.games.GetAt(id, numActions)
	if err != nil {
		return model.Game{}, err
	}

	for i, player := range g.Players {
		// overwrite the player that the game service knows
		// about with the player that the players service knows about
		p, err := d.GetPlayer(player.ID)
		if err != nil {
			return model.Game{}, err
		}
		g.Players[i] = p
	}

	return g, nil
}

func (d *db) CreateGame(g model.Game) error {
	if g.NumActions() != 0 {
		return errors.New(`cannot create game with actions`)
	}

	err := d.players.BeginGame(g.ID, g.Players)
	if err != nil {
		return err
	}

	return d.games.Begin(g)
}

func (d *db) SaveGame(g model.Game) error {
	return d.games.Save(g)
}

func (d *db) GetInteraction(id model.PlayerID) (interaction.PlayerMeans, error) {
	return d.interactions.Get(id)
}

func (d *db) SaveInteraction(pm interaction.PlayerMeans) error {
	return d.interactions.Update(pm)
}
