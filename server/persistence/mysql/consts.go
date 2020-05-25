package mysql

import (
	"math"

	"github.com/joshprzybyszewski/cribbage/model"
)

const (
	maxPlayerUUIDLen    = 255
	maxPlayerUUIDLenStr = `255`

	maxPlayerNameLen    = 255
	maxPlayerNameLenStr = `255`

	maxGameID = model.GameID(math.MaxUint32)
)
