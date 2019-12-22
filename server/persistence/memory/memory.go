package memory

import (
	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

func New() persistence.DB {
	return persistence.New(
		getGameService(),
		getPlayerService(),
		getInteractionService(),
	)
}
