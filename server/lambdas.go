package server

import (
	"context"

	"github.com/joshprzybyszewski/cribbage/model"
)

func HandleAction(ctx context.Context, action model.PlayerAction) error {
	db, err := getDB(ctx)
	if err != nil {
		return err
	}
	return handleAction(db, action)
}

func CreateGame(ctx context.Context, pIDs []model.PlayerID) (model.Game, error) {
	db, err := getDB(ctx)
	if err != nil {
		return model.Game{}, err
	}

	return createGame(ctx, db, pIDs)
}

func GetGame(ctx context.Context, gID model.GameID) (model.Game, error) {
	db, err := getDB(ctx)
	if err != nil {
		return model.Game{}, err
	}

	return db.GetGame(gID)
}

func GetPlayer(ctx context.Context, pID model.PlayerID) (model.Player, error) {
	db, err := getDB(ctx)
	if err != nil {
		return model.Player{}, err
	}

	return db.GetPlayer(pID)
}
