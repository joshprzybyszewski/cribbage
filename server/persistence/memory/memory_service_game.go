package memory

import (
	"errors"
	"sync"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

var gservice *gameService
var _ persistence.GameService = (*gameService)(nil)

type gameService struct {
	lock sync.Mutex

	games map[model.GameID][]model.Game
}

func getGameService() persistence.GameService {
	if gservice == nil {
		gservice = &gameService{
			games: map[model.GameID][]model.Game{},
		}
	}
	return gservice
}

func (gs *gameService) Get(id model.GameID) (model.Game, error) {
	gs.lock.Lock()
	defer gs.lock.Unlock()

	if games, ok := gs.games[id]; ok {
		g := games[len(games)-1]
		return g, nil
	}
	return model.Game{}, persistence.ErrGameNotFound
}

func (gs *gameService) GetAt(id model.GameID, numActions uint) (model.Game, error) {
	gs.lock.Lock()
	defer gs.lock.Unlock()

	if games, ok := gs.games[id]; ok {
		if int(numActions) >= len(games) {
			return model.Game{}, persistence.ErrGameNotFound
		}
		g := games[numActions]
		return g, nil
	}
	return model.Game{}, persistence.ErrGameNotFound
}

func (gs *gameService) UpdatePlayerColor(gID model.GameID, pID model.PlayerID, color model.PlayerColor) error {
	gs.lock.Lock()
	defer gs.lock.Unlock()

	gameList := gs.games[gID]
	mostRecent := gameList[len(gameList)-1]
	if c, ok := mostRecent.PlayerColors[pID]; !ok {
		if mostRecent.PlayerColors == nil {
			mostRecent.PlayerColors = make(map[model.PlayerID]model.PlayerColor, 1)
		}
		mostRecent.PlayerColors[pID] = color
		gameList[len(gameList)-1] = mostRecent
	} else if c != color {
		return errors.New(`mismatched game-player colors`)
	}

	return nil
}

func (gs *gameService) Begin(g model.Game) error {
	return gs.Save(g)
}

func (gs *gameService) Save(g model.Game) error {
	gs.lock.Lock()
	defer gs.lock.Unlock()

	id := g.ID

	savedGames := gs.games[id]
	err := validateGameState(savedGames, g)
	if err != nil {
		return err
	}

	gs.games[id] = append(gs.games[id], g)

	return nil
}

func validateGameState(savedGames []model.Game, newGameState model.Game) error {
	if len(savedGames) != newGameState.NumActions() {
		return persistence.ErrGameActionsOutOfOrder
	}
	for i := range savedGames {
		savedActions := savedGames[i].Actions
		myKnownActions := newGameState.Actions[:i]
		if len(savedActions) != len(myKnownActions) {
			return persistence.ErrGameActionsOutOfOrder
		}
		for ai, a := range savedActions {
			if a.ID != myKnownActions[ai].ID || a.Overcomes != myKnownActions[ai].Overcomes {
				return persistence.ErrGameActionsOutOfOrder
			}
		}
	}
	return nil
}
