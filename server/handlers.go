package server

import (
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
	"github.com/joshprzybyszewski/cribbage/server/play"
)

func (cs *cribbageServer) getGame(gID model.GameID) (model.Game, error) {
	return cs.db.GetGame(gID)
}

func (cs *cribbageServer) getPlayer(pID model.PlayerID) (model.Player, error) {
	return cs.db.GetPlayer(pID)
}

func (cs *cribbageServer) createGame(pIDs []model.PlayerID) (model.Game, error) {
	players := make([]model.Player, len(pIDs))
	for i, id := range pIDs {
		p, err := cs.db.GetPlayer(id)
		if err != nil {
			return model.Game{}, err
		}
		players[i] = p
	}
	mg := play.New(players)
	err := cs.db.SaveGame(mg)
	if err != nil {
		return model.Game{}, err
	}
	for _, pID := range pIDs {
		err := cs.db.AddPlayerColorToGame(pID, mg.PlayerColors[pID], mg.ID)
		if err != nil {
			return model.Game{}, err
		}
	}
	return mg, nil
}

func (cs *cribbageServer) createPlayer(name string) (model.Player, error) {
	mp := model.Player{
		ID:    model.NewPlayerID(),
		Name:  name,
		Games: make(map[model.GameID]model.PlayerColor),
	}
	err := cs.db.CreatePlayer(mp)
	if err != nil {
		return model.Player{}, err
	}
	return mp, nil
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
