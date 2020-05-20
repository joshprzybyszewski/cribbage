package server

import (
	"context"
	"log"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

func commitOrRollback(db persistence.DB, err *error) {
	var err2 error
	if *err != nil {
		err2 = db.Rollback()
	} else {
		err2 = db.Commit()
	}
	if err2 != nil {
		log.Printf("Could not commit/rollback after %+v: %+v\n", err, err2)
	}
}

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

func saveInteraction(ctx context.Context, pm interaction.PlayerMeans) error {
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

	return db.SaveInteraction(pm)
}

func createPlayer(ctx context.Context, p model.Player) error {
	db, err := getDB(ctx)
	if err != nil {
		return err
	}
	defer db.Close()

	if !model.IsValidPlayerID(p.ID) {
		return errInvalidUsername
	}

	err = db.Start()
	if err != nil {
		return err
	}
	defer commitOrRollback(db, &err)

	return db.CreatePlayer(p)
}
