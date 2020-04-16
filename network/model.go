package network

type CreatePlayerRequest struct {
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
}
