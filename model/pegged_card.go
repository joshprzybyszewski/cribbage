package model

func NewPeggedCard(pID PlayerID, c Card) PeggedCard {
	return PeggedCard{
		Card: c,
		PlayerID: pID,
	}
}