package server

import (
	"context"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
	"github.com/joshprzybyszewski/cribbage/server/play"
)

func HandleAction(ctx context.Context, action model.PlayerAction) error {
	db, err := getDB()
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
	db, err := getDB()
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

	err = db.SaveGame(mg)
	if err != nil {
		return model.Game{}, err
	}

	for _, pID := range pIDs {
		err := db.AddPlayerColorToGame(pID, mg.PlayerColors[pID], mg.ID)
		if err != nil {
			return model.Game{}, err
		}
	}

	return mg, nil
}
