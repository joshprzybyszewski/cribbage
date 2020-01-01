package interaction

import (
	"github.com/joshprzybyszewski/cribbage/logic/strategy"
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/utils/rand"
)

var _ npc = (*simpleNPC)(nil)

type simpleNPC struct{}

func (npc *simpleNPC) getBuildCribAction(hand []model.Card, isDealer bool) (model.BuildCribAction, error) {
	return cribActionHelper(hand, Simple, isDealer)
}

func (npc *simpleNPC) getPegAction(hand []model.Card, prevPegs []model.PeggedCard, curPeg int) model.PegAction {
	// try random strategies until we either have to say go or have a valid peg card
	card, sayGo := randomPegStrategy(hand, prevPegs, curPeg)
	for card.PegValue()+curPeg > 31 && !sayGo {
		card, sayGo = randomPegStrategy(hand, prevPegs, curPeg)
	}
	return model.PegAction{
		Card:  card,
		SayGo: sayGo,
	}
}

func randomPegStrategy(hand []model.Card, prevPegs []model.PeggedCard, curPeg int) (model.Card, bool) {
	var card model.Card
	var sayGo bool
	switch rand.Int() % 4 {
	case 0:
		card, sayGo = strategy.PegToFifteen(hand, prevPegs, curPeg)
	case 1:
		card, sayGo = strategy.PegToThirtyOne(hand, prevPegs, curPeg)
	case 2:
		card, sayGo = strategy.PegToPair(hand, prevPegs, curPeg)
	default:
		card, sayGo = strategy.PegToRun(hand, prevPegs, curPeg)
	}
	return card, sayGo
}
