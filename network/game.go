package network

import "github.com/joshprzybyszewski/cribbage/model"

func NewGetGameResponse(g model.Game) GetGameResponse {
	currentScores, lagScores := convertScores(g.CurrentScores, g.LagScores)
	return GetGameResponse{
		ID:              g.ID,
		Players:         newPlayersFromModels(g.Players),
		PlayerColors:    convertColors(g.PlayerColors),
		CurrentScores:   currentScores,
		LagScores:       lagScores,
		Phase:           g.Phase.String(),
		BlockingPlayers: convertBlockingPlayers(g.BlockingPlayers),
		CurrentDealer:   g.CurrentDealer,
		CutCard:         newCardFromModel(g.CutCard),
		PeggedCards:     convertPeggedCards(g.PeggedCards),
	}
}

func NewGetGameResponseForPlayer(g model.Game, pID model.PlayerID) GetGameResponse {
	resp := NewGetGameResponse(g)
	resp.Hands = convertHands(g.Hands)
	if g.Phase < model.Counting {
		resp.Hands = map[model.PlayerID][]Card{
			pID: resp.Hands[pID],
		}
	}
	if g.Phase >= model.CribCounting {
		resp.Crib = convertCards(g.Crib)
	}
	return resp
}

func NewCreateGameResponse(g model.Game) CreateGameResponse {
	return CreateGameResponse{
		ID:              g.ID,
		Players:         newPlayersFromModels(g.Players),
		PlayerColors:    convertColors(g.PlayerColors),
		BlockingPlayers: convertBlockingPlayers(g.BlockingPlayers),
		CurrentDealer:   g.CurrentDealer,
	}
}

func convertColors(cs map[model.PlayerID]model.PlayerColor) map[model.PlayerID]string {
	colors := make(map[model.PlayerID]string, len(cs))
	for p, c := range cs {
		colors[p] = c.String()
	}
	return colors
}

func convertScores(modelCurrentScores, modelLagScores map[model.PlayerColor]int) (map[string]int, map[string]int) {
	current := make(map[string]int, len(modelCurrentScores))
	for c, s := range modelCurrentScores {
		current[c.String()] = s
	}
	lag := make(map[string]int, len(modelLagScores))
	for c, s := range modelLagScores {
		current[c.String()] = s
	}
	return current, lag
}

func convertBlockingPlayers(bs map[model.PlayerID]model.Blocker) map[model.PlayerID]string {
	blockers := make(map[model.PlayerID]string, len(bs))
	for p, b := range bs {
		blockers[p] = b.String()
	}
	return blockers
}

func convertHands(hands map[model.PlayerID][]model.Card) map[model.PlayerID][]Card {
	networkHands := make(map[model.PlayerID][]Card, len(hands))
	for p, h := range hands {
		networkHands[p] = convertCards(h)
	}
	return networkHands
}

func convertCards(modelCards []model.Card) []Card {
	cards := make([]Card, len(modelCards))
	for i, c := range modelCards {
		cards[i] = newCardFromModel(c)
	}
	return cards
}

func convertPeggedCards(modelPeggedCards []model.PeggedCard) []Card {
	cards := make([]Card, len(modelPeggedCards))
	for i, pc := range modelPeggedCards {
		cards[i] = newCardFromModel(pc.Card)
	}
	return cards
}
