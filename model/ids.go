package model

import (
	"log"
	"regexp"

	"github.com/google/uuid"
)

var (
	validPIDRegex = regexp.MustCompile("^[\\w]*$")
)

func NewGameID() GameID {
	gID := InvalidGameID
	for gID == InvalidGameID {
		r, err := uuid.NewRandom()
		if err != nil {
			log.Printf("NewGameID.NewRandom failed\n")
			return InvalidGameID
		}

		gID = GameID(r.ID())
	}

	return gID
}

func IsValidPlayerID(pID PlayerID) bool {
	return validPIDRegex.MatchString(string(pID))
}
