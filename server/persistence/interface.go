package persistence

import (
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
)

type DB interface {
	// Close should be called to close any connections needed on the database
	Close() error

	// Start will start a transaction on the database
	Start() error
	// Commit will commit a transaction, if one exists
	Commit() error
	// Rollback will rollback the transaction of changes on the database
	Rollback() error

	ServicesWrapper
}

type ServicesWrapper interface {
	CreatePlayer(p model.Player) error
	GetPlayer(id model.PlayerID) (model.Player, error)
	AddPlayerColorToGame(id model.PlayerID, color model.PlayerColor, gID model.GameID) error

	GetGame(id model.GameID) (model.Game, error)
	GetGameAction(id model.GameID, numActions uint) (model.Game, error)
	SaveGame(g model.Game) error

	GetInteraction(id model.PlayerID) (interaction.PlayerMeans, error)
	SaveInteraction(pm interaction.PlayerMeans) error
}

type services struct {
	games        GameService
	players      PlayerService
	interactions InteractionService
}

func NewServicesWrapper(gs GameService, ps PlayerService, is InteractionService) ServicesWrapper {
	return &services{
		games:        gs,
		players:      ps,
		interactions: is,
	}
}

func (d *services) CreatePlayer(p model.Player) error {
	return d.players.Create(p)
}

func (d *services) GetPlayer(id model.PlayerID) (model.Player, error) {
	return d.players.Get(id)
}

func (d *services) AddPlayerColorToGame(pID model.PlayerID, color model.PlayerColor, gID model.GameID) error {
	err := d.games.UpdatePlayerColor(gID, pID, color)
	if err != nil {
		return err
	}
	return d.players.UpdateGameColor(pID, gID, color)
}

func (d *services) GetGame(id model.GameID) (model.Game, error) {
	return d.games.Get(id)
}

func (d *services) GetGameAction(id model.GameID, numActions uint) (model.Game, error) {
	return d.games.GetAt(id, numActions)
}

func (d *services) SaveGame(g model.Game) error {
	return d.games.Save(g)
}

func (d *services) GetInteraction(id model.PlayerID) (interaction.PlayerMeans, error) {
	return d.interactions.Get(id)
}

func (d *services) SaveInteraction(pm interaction.PlayerMeans) error {
	return d.interactions.Update(pm)
}
