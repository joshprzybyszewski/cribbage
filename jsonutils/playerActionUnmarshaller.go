package jsonutils

import (
	"encoding/json"
	"errors"

	"github.com/joshprzybyszewski/cribbage/model"
)

// UnmarshalPlayerAction takes json-marshalled bytes and returns the model.PlayerAction
// The advantage is that we can unmarshal the Action which is typed based on the Blocker
func UnmarshalPlayerAction(b []byte) (model.PlayerAction, error) {
	// We can store the RawMessage and then switch on the Overcomes type later
	// otherwise Action becomes a map[string]interface{}
	var raw json.RawMessage
	action := model.PlayerAction{
		Action: &raw,
	}
	err := json.Unmarshal(b, &action)
	if err != nil {
		return model.PlayerAction{}, err
	}

	err = unmarshalActionIntoPlayerAction(&action, raw)
	if err != nil {
		return model.PlayerAction{}, err
	}

	return action, nil
}

func unmarshalActionIntoPlayerAction(
	action *model.PlayerAction,
	raw json.RawMessage,
) error {

	blockerActions := map[model.Blocker]func() interface{}{
		model.DealCards: func() interface{} { return &model.DealAction{} },
		model.CribCard:  func() interface{} { return &model.BuildCribAction{} },
		model.CutCard:   func() interface{} { return &model.CutDeckAction{} },
		model.PegCard:   func() interface{} { return &model.PegAction{} },
		model.CountHand: func() interface{} { return &model.CountHandAction{} },
		model.CountCrib: func() interface{} { return &model.CountCribAction{} },
	}

	subActionFn, ok := blockerActions[action.Overcomes]
	if !ok {
		return errors.New(`unknown action type`)
	}
	subAction := subActionFn()

	if err := json.Unmarshal(raw, subAction); err != nil {
		return err
	}

	switch t := subAction.(type) {
	case *model.DealAction:
		action.Action = *t
	case *model.BuildCribAction:
		action.Action = *t
	case *model.CutDeckAction:
		action.Action = *t
	case *model.PegAction:
		action.Action = *t
	case *model.CountHandAction:
		action.Action = *t
	case *model.CountCribAction:
		action.Action = *t
	}

	return nil
}
