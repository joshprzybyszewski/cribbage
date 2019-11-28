package model

import (
	"log"

	"github.com/google/uuid"
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
