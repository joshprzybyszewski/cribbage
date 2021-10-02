package jsonutils

import (
	"encoding/json"

	"github.com/joshprzybyszewski/cribbage/model"
)

// UnmarshalGame takes in json marshaled bytes of a model.Game
// The main advantage is that the list of actions can be deserialized
// into the interface{} type.
func UnmarshalGame(b []byte) (model.Game, error) {
	game := model.Game{}

	err := json.Unmarshal(b, &game)
	if err != nil {
		return model.Game{}, err
	}

	for i := range game.Actions {
		a := game.Actions[i]
		b, err := json.Marshal(a)
		if err != nil {
			return model.Game{}, err
		}

		pa, err := UnmarshalPlayerAction(b)
		if err != nil {
			return model.Game{}, err
		}
		game.Actions[i] = pa
	}

	if game.Hands == nil {
		game.Hands = make(map[model.PlayerID][]model.Card, len(game.Players))
	}

	if game.BlockingPlayers == nil {
		game.BlockingPlayers = make(map[model.PlayerID]model.Blocker, len(game.Players))
	}

	if game.PlayerColors == nil {
		game.PlayerColors = make(map[model.PlayerID]model.PlayerColor, len(game.Players))
	}

	return game, nil
}
