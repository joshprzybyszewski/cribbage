package model

func NewPeggedCard(pID PlayerID, c Card, numActions int) PeggedCard {
	return PeggedCard{
		Card:     c,
		PlayerID: pID,
		Action: numActions,
	}
}
