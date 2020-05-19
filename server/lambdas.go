package server

import (
	"context"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
	"github.com/joshprzybyszewski/cribbage/server/play"
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

func createGame(_ context.Context, db persistence.DB, pIDs []model.PlayerID) (model.Game, error) {
	players := make([]model.Player, len(pIDs))
	for i, id := range pIDs {
		p, err := db.GetPlayer(id)
		if err != nil {
			return model.Game{}, err
		}
		players[i] = p
	}

	pAPIs, err := getPlayerAPIs(db, players)
	if err != nil {
		return model.Game{}, err
	}

	mg, err := play.CreateGame(players, pAPIs)
	if err != nil {
		return model.Game{}, err
	}

	err = db.CreateGame(mg)
	if err != nil {
		return model.Game{}, err
	}

	return mg, nil
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
