package server

import (
	"context"
	"errors"
	"regexp"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
)

var (
	errInvalidUsername error = errors.New(`invalid username`)
)

func (cs *cribbageServer) getGame(gID model.GameID) (model.Game, error) {
	return cs.db.GetGame(gID)
}

func (cs *cribbageServer) getPlayer(pID model.PlayerID) (model.Player, error) {
	return cs.db.GetPlayer(pID)
}

func (cs *cribbageServer) createGame(pIDs []model.PlayerID) (model.Game, error) {
	return createGame(context.Background(), cs.db, pIDs)
}

func (cs *cribbageServer) createPlayer(username, name string) (model.Player, error) {
	// TODO move this into a lambda
	re := regexp.MustCompile("^[a-zA-Z0-9_]*$")
	if !re.MatchString(username) {
		return model.Player{}, errInvalidUsername
	}

	mp := model.Player{
		ID:    model.PlayerID(username),
		Name:  name,
		Games: make(map[model.GameID]model.PlayerColor),
	}
	err := cs.db.CreatePlayer(mp)
	if err != nil {
		return model.Player{}, err
	}
	return mp, nil
}

func (cs *cribbageServer) setInteraction(pID model.PlayerID, im interaction.Means) error {
	// TODO have a way to get the previous interaction and then update with this as the preferred mode
	pm := interaction.New(pID, im)
	return cs.db.SaveInteraction(pm)
}

func (cs *cribbageServer) handleAction(action model.PlayerAction) error {
	return handleAction(context.Background(), cs.db, action)
}
