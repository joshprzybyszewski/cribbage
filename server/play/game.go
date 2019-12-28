package play

import (
	"errors"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

var (
	ErrActionNotForGame error = errors.New(`action not for game`)
	ErrPlayerNotInGame  error = errors.New(`player is not in this game`)
	ErrGameAlreadyOver  error = errors.New(`game is already over`)
)

func CreateGame(players []model.Player, pAPIs map[model.PlayerID]interaction.Player, db persistence.DB) (model.Game, error) {
	playersCopy := make([]model.Player, len(players))
	colorsByID := make(map[model.PlayerID]model.PlayerColor, len(players))
	curScores := make(map[model.PlayerColor]int, len(players))
	lagScores := make(map[model.PlayerColor]int, len(players))

	playerColors := []model.PlayerColor{
		model.Blue,
		model.Red,
	}
	if len(players) == 3 {
		playerColors = append(playerColors, model.Green)
	} else if len(players) == 4 {
		playerColors = append(playerColors, model.Blue, model.Red)
	}

	for i, p := range players {
		playersCopy[i] = p
		color := playerColors[i]
		colorsByID[p.ID] = color
		curScores[color] = 0
		lagScores[color] = 0
	}

	g := model.Game{
		ID:              model.NewGameID(),
		Players:         playersCopy,
		BlockingPlayers: make(map[model.PlayerID]model.Blocker, len(players)),
		CurrentDealer:   players[0].ID,
		PlayerColors:    colorsByID,
		CurrentScores:   curScores,
		LagScores:       lagScores,
		Phase:           model.DealingReady,
		Hands:           make(map[model.PlayerID][]model.Card, len(players)),
		CutCard:         model.Card{},
		Crib:            make([]model.Card, 0, 4),
		PeggedCards:     make([]model.PeggedCard, 0, 4*len(players)),
	}

	err := db.SaveGame(g)
	if err != nil {
		return model.Game{}, err
	}

	// TODO should we actually run start handlers upon game creation? This can lead to
	// trying to get this game from the DB before it's saved to the DB
	err = runStartHandlers(&g, pAPIs)
	if err != nil {
		return model.Game{}, err
	}

	return g, nil
}

var (
	handlers = map[model.Phase]PhaseHandler{
		model.Deal:              &dealingHandler{},
		model.BuildCribReady:    &cribBuildingHandler{},
		model.BuildCrib:         &cribBuildingHandler{},
		model.CutReady:          &cuttingHandler{},
		model.Cut:               &cuttingHandler{},
		model.PeggingReady:      &peggingHandler{},
		model.Pegging:           &peggingHandler{},
		model.CountingReady:     &handCountingHandler{},
		model.Counting:          &handCountingHandler{},
		model.CribCountingReady: &cribCountingHandler{},
		model.CribCounting:      &cribCountingHandler{},
		model.DealingReady:      &dealingHandler{},
	}
)

func HandleAction(g *model.Game,
	action model.PlayerAction,
	pAPIs map[model.PlayerID]interaction.Player,
) error {

	if g.ID != action.GameID {
		return ErrActionNotForGame
	}
	playerIsInGame := false
	for i := range g.Players {
		if g.Players[i].ID == action.ID {
			playerIsInGame = true
			break
		}
	}
	if !playerIsInGame {
		return ErrPlayerNotInGame
	}
	if g.IsOver() {
		return ErrGameAlreadyOver
	}
	switch p := g.Phase; p {
	case model.Deal,
		model.BuildCrib,
		model.Cut,
		model.Pegging,
		model.Counting,
		model.CribCounting:
		err := handlers[p].HandleAction(g, action, pAPIs)
		if err != nil {
			return err
		}
		g.AddAction(action)
	}

	if g.IsOver() {
		return nil
	}

	if len(g.BlockingPlayers) == 0 {
		g.Phase++
	}

	return runStartHandlers(g, pAPIs)
}

func runStartHandlers(g *model.Game, pAPIs map[model.PlayerID]interaction.Player) error {
	switch p := g.Phase; p {
	case model.BuildCribReady,
		model.CutReady,
		model.PeggingReady,
		model.CountingReady,
		model.CribCountingReady,
		model.DealingReady:
		err := handlers[p].Start(g, pAPIs)
		if err != nil {
			return err
		}
		g.Phase++
		if g.Phase > model.DealingReady {
			g.Phase = model.Deal
		}
	}

	return nil
}
