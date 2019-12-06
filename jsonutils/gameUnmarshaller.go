package jsonutils

import (
	"encoding/json"

	"github.com/joshprzybyszewski/cribbage/model"
)

func UnmarshalGame(b []byte) (model.Game, error) {
	game := model.Game{}

	err := json.Unmarshal(b, &game)
	if err != nil {
		return model.Game{}, err
	}

	for i, a := range game.Actions {
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

	return game, nil
}
