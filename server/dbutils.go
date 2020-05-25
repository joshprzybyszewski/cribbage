package server

import (
	"context"
	"log"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
	"github.com/joshprzybyszewski/cribbage/server/play"
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
	err := db.Start()
	if err != nil {
		return err
	}
	defer commitOrRollback(db, &err)

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
	err := db.Start()
	if err != nil {
		return model.Game{}, err
	}
	defer commitOrRollback(db, &err)

	var p model.Player
	players := make([]model.Player, len(pIDs))
	for i, id := range pIDs {
		p, err = db.GetPlayer(id)
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

func getGame(_ context.Context, db persistence.DB, gID model.GameID) (model.Game, error) {
	err := db.Start()
	if err != nil {
		return model.Game{}, err
	}
	defer commitOrRollback(db, &err)

	return db.GetGame(gID)
}

func getPlayer(_ context.Context, db persistence.DB, pID model.PlayerID) (model.Player, error) {
	err := db.Start()
	if err != nil {
		return model.Player{}, err
	}
	defer commitOrRollback(db, &err)

	return db.GetPlayer(pID)
}

func saveInteraction(_ context.Context, db persistence.DB, pm interaction.PlayerMeans) error {
	err := db.Start()
	if err != nil {
		return err
	}
	defer commitOrRollback(db, &err)

	return db.SaveInteraction(pm)
}

func createPlayer(_ context.Context, db persistence.DB, p model.Player) error {
	err := db.Start()
	if err != nil {
		return err
	}
	defer commitOrRollback(db, &err)

	return db.CreatePlayer(p)
}
