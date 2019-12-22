package server

import (
	"context"
	"time"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
)

const (
	defaultTimeout time.Duration = 10 * time.Second
)

func getGame(gID model.GameID) (model.Game, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	return GetGame(ctx, gID)
}

func getPlayer(pID model.PlayerID) (model.Player, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	return GetPlayer(ctx, pID)
}

func (cs *cribbageServer) createGame(pIDs []model.PlayerID) (model.Game, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	return CreateGame(ctx, pIDs)
}

func createPlayerFromNames(username, name string) (model.Player, error) {
	mp := model.Player{
		ID:    model.PlayerID(username),
		Name:  name,
		Games: make(map[model.GameID]model.PlayerColor),
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	err := createPlayer(ctx, mp)
	if err != nil {
		return model.Player{}, err
	}
	return mp, nil
}

func (cs *cribbageServer) setInteraction(pID model.PlayerID, im interaction.Means) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	// TODO have a way to get the previous interaction and then update with this as the preferred mode
	pm := interaction.New(pID, im)
	return saveInteraction(ctx, pm)
}

func (cs *cribbageServer) handleAction(action model.PlayerAction) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	return HandleAction(ctx, action)
}
