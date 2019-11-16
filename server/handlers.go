package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
	"github.com/joshprzybyszewski/cribbage/server/play"
)

func (cs *cribbageServer) createGame(pIDs []model.PlayerID) (model.GameID, error) {
	mg := play.New(pIDs)
	err := cs.db.SaveGame(mg)
	if err != nil {
		return model.GameID(-1), err
	}
	for _, pID := range pIDs {
		err := cs.db.AddPlayerColorToGame(pID, mg.PlayerColors[pID], mg.ID)
		if err != nil {
			return model.GameID(-1), err
		}
	}
	return mg.ID, nil
}

func (cs *cribbageServer) createPlayer(name string) (model.PlayerID, error) {
	pID := model.NewPlayerID()
	mp := model.Player{
		ID:    pID,
		Name:  name,
		Games: make(map[model.GameID]model.PlayerColor),
	}
	err := cs.db.CreatePlayer(mp)
	if err != nil {
		return model.PlayerID(-1), err
	}
	return pID, nil
}

func (cs *cribbageServer) setInteraction(pID model.PlayerID) error {
	// TODO
	var ip interaction.Player
	err := cs.db.SaveInteraction(ip)
	if err != nil {
		return err
	}
	return nil
}

func (cs *cribbageServer) handleAction(action model.PlayerAction) error {
	g, err := cs.db.GetGame(action.GameID)
	if err != nil {
		return err
	}

	pAPIs := make(map[model.PlayerID]interaction.Player, len(g.Players))
	for _, p := range g.Players {
		pAPI, err := cs.db.GetInteraction(p.ID)
		if err != nil {
			return err
		}
		pAPIs[p.ID] = pAPI
	}

	err = play.HandleAction(&g, action, pAPIs)
	if err != nil {
		return err
	}

	return cs.db.SaveGame(g)
}
