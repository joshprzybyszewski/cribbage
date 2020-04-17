package network

type CreatePlayerRequest struct {
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
}

type CreateGameRequest struct {
	PlayerIDs []string `json:"playerIDs"`
}

type CreateInteractionRequest struct {
	PlayerID string      `json:"playerID"`
	Mode     string      `json:"mode"`
	Info     interface{} `json:"info"`
}
