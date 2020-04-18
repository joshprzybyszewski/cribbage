package network

type CreateInteractionRequest struct {
	PlayerID string      `json:"playerID"`
	Mode     string      `json:"mode"`
	Info     interface{} `json:"info"`
}
