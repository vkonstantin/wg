package message

// Error that should be returned to the client
type Error struct {
	HttpCode int    `json:"-"`
	Message  string `json:"message"`
}
