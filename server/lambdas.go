package server

import (
	"context"

	"github.com/joshprzybyszewski/cribbage/model"
)

func HandleAction(ctx context.Context, action model.PlayerAction) (err error) {
	db, err := getDB(ctx)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Start()
	if err != nil {
		return err
	}
	defer commitOrRollback(db, &err)

	return handleAction(ctx, db, action)
}

func CreateGame(ctx context.Context, pIDs []model.PlayerID) (model.Game, error) {
	db, err := getDB(ctx)
	if err != nil {
		return model.Game{}, err
	}
	defer db.Close()

	err = db.Start()
	if err != nil {
		return model.Game{}, err
	}
	defer commitOrRollback(db, &err)

	return createGame(ctx, db, pIDs)
}

func GetGame(ctx context.Context, gID model.GameID) (model.Game, error) {
	db, err := getDB(ctx)
	if err != nil {
		return model.Game{}, err
	}
	defer db.Close()

	err = db.Start()
	if err != nil {
		return model.Game{}, err
	}
	defer commitOrRollback(db, &err)

	return db.GetGame(gID)
}

func GetPlayer(ctx context.Context, pID model.PlayerID) (model.Player, error) {
	db, err := getDB(ctx)
	if err != nil {
		return model.Player{}, err
	}
	defer db.Close()

	err = db.Start()
	if err != nil {
		return model.Player{}, err
	}
	defer commitOrRollback(db, &err)

	return db.GetPlayer(pID)
}
