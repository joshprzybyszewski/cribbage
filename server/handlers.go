package server

import (
	"errors"
	"regexp"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
	"github.com/joshprzybyszewski/cribbage/server/play"
)

var (
	errInvalidUsername error = errors.New(`invalid username`)
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
	pAPIs, err := cs.getPlayerAPIs(players)
	if err != nil {
		return model.Game{}, err
	}
	mg, err := play.CreateGame(players, pAPIs)
	if err != nil {
		return model.Game{}, err
	}
	err = cs.db.SaveGame(mg)
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

func (cs *cribbageServer) createPlayer(username, name string) (model.Player, error) {
	re := regexp.MustCompile("^[a-zA-Z0-9_]*$")
	if !re.MatchString(username) {
		return model.Player{}, errInvalidUsername
	}

	mp := model.Player{
		ID:    model.PlayerID(username),
		Name:  name,
		Games: make(map[model.GameID]model.PlayerColor),
	}
	err := cs.db.CreatePlayer(mp)
	if err != nil {
		return model.Player{}, err
	}
	return mp, nil
}

func (cs *cribbageServer) setInteraction(pID model.PlayerID, im model.InteractionMeans) error {
	ip := interaction.New(pID, im)
	return cs.db.SaveInteraction(ip)
}

func (cs *cribbageServer) handleAction(action model.PlayerAction) error {
	g, err := cs.db.GetGame(action.GameID)
	if err != nil {
		return err
	}

	pAPIs, err := cs.getPlayerAPIs(g.Players)
	if err != nil {
		return err
	}

	err = play.HandleAction(&g, action, pAPIs)
	if err != nil {
		return err
	}

	return cs.db.SaveGame(g)
}

func (cs *cribbageServer) getPlayerAPIs(players []model.Player) (map[model.PlayerID]interaction.Player, error) {
	pAPIs := make(map[model.PlayerID]interaction.Player, len(players))
	for _, p := range players {
		pAPI, err := cs.db.GetInteraction(p.ID)
		if err != nil {
			return nil, err
		}
		pAPIs[p.ID] = pAPI
	}
	return pAPIs, nil
}
