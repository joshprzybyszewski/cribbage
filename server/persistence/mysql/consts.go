package mysql

import (
	"math"

	"github.com/joshprzybyszewski/cribbage/model"
)

const (
	maxPlayerUUIDLen    = 191
	maxPlayerUUIDLenStr = `191`

	maxPlayerNameLen    = 191
	maxPlayerNameLenStr = `191`

	maxGameID = model.GameID(math.MaxUint32)
)
