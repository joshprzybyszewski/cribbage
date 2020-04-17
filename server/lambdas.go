package server

import (
	"context"
	"errors"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
	"github.com/joshprzybyszewski/cribbage/server/play"
)

var (
	errInvalidUsername error = errors.New(`invalid username`)
)

func HandleAction(ctx context.Context, action model.PlayerAction) error {
	db, err := getDB(ctx)
	if err != nil {
		return err
	}
	return handleAction(ctx, db, action)
}

func handleAction(_ context.Context, db persistence.DB, action model.PlayerAction) error {
	g, err := db.GetGame(action.GameID)
	if err != nil {
		return err
	}

	pAPIs, err := getPlayerAPIs(db, g.Players)
	if err != nil {
		return err
	}
	err = play.HandleAction(&g, action, pAPIs)
	if err != nil {
		return err
	}
	return db.SaveGame(g)
}

func CreateGame(ctx context.Context, pIDs []model.PlayerID) (model.Game, error) {
	db, err := getDB(ctx)
	if err != nil {
		return model.Game{}, err
	}

	return createGame(db, pIDs)
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

func saveInteraction(ctx context.Context, pm interaction.PlayerMeans) error {
	db, err := getDB(ctx)
	if err != nil {
		return err
	}

	return db.SaveInteraction(pm)
}

func createPlayer(ctx context.Context, p model.Player) error {
	db, err := getDB(ctx)
	if err != nil {
		return err
	}

	if !model.IsValidPlayerID(p.ID) {
		return errInvalidUsername
	}

	return db.CreatePlayer(p)
}
