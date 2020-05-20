package server

import (
	"context"
	"errors"
	"log"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
	"github.com/joshprzybyszewski/cribbage/server/play"
)

var (
	errInvalidUsername error = errors.New(`invalid username`)
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

// TODO this is unused?
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

// TODO this is unused?
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
