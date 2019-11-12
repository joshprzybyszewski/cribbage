package model

import (
	"log"

	"github.com/google/uuid"
)

func NewGameID() GameID {
	r, err := uuid.NewRandom()
	if err != nil {
		log.Printf("NewRandom failed\n")
		return GameID(-1)
	}

	return GameID(r.ID())
}

func NewPlayerID() PlayerID {
	r, err := uuid.NewRandom()
	if err != nil {
		log.Printf("NewRandom failed\n")
		return PlayerID(-1)
	}

	return PlayerID(r.ID())
}
