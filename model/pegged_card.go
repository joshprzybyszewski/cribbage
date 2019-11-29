package model

func NewPeggedCard(pID PlayerID, c Card, numActions int) PeggedCard {
	return PeggedCard{
		Card:     c,
		PlayerID: pID,
		Action:   numActions,
	}
}

func NewPeggedCardFromString(pID PlayerID, cStr string, numActions int) PeggedCard {
	return PeggedCard{
		Card:     NewCardFromString(cStr),
		PlayerID: pID,
		Action:   numActions,
	}
}
