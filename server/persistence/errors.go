package persistence

import (
	"errors"
)

var (
	ErrInvalidPlayerID     error = errors.New(`player id not valid`)
	ErrInvalidPlayerName   error = errors.New(`player name not valid`)
	ErrPlayerNotFound      error = errors.New(`player not found`)
	ErrPlayerAlreadyExists error = errors.New(`username already exists`)

	ErrInvalidGameID         error = errors.New(`game id invalid`)
	ErrGameNotFound          error = errors.New(`game not found`)
	ErrGameInitialSave       error = errors.New(`game must be saved with no actions`)
	ErrGameActionsOutOfOrder error = errors.New(`game actions out of order`)
	ErrGameActionSave        error = errors.New(`game actions save race`)
	ErrGameActionDecode      error = errors.New(`game actions get decode`)
	ErrGameActionWrongGame   error = errors.New(`game action for wrong game`)
	ErrGameActionWrongPlayer error = errors.New(`game action found for wrong player`)

	ErrInteractionNotFound      error = errors.New(`interaction not found`)
	ErrInteractionAlreadyExists error = errors.New(`interaction already exists`)
	ErrInteractionUnexpected    error = errors.New(`unexpected interaction`)
)
