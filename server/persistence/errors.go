package persistence

import (
	"errors"
)

var (
	ErrPlayerNotFound      error = errors.New(`player not found`)
	ErrPlayerAlreadyExists error = errors.New(`username already exists`)
	ErrPlayerColorMismatch error = errors.New(`mismatched player-games color`)

	ErrGameNotFound          error = errors.New(`game not found`)
	ErrGameInitialSave       error = errors.New(`game must be saved with no actions`)
	ErrGameActionsOutOfOrder error = errors.New(`game actions out of order`)

	ErrInteractionNotFound      error = errors.New(`interaction not found`)
	ErrInteractionAlreadyExists error = errors.New(`interaction already exists`)
)
