package message

type AddUserRequest struct {
	// RequestID is an important field for deduplication requests.
	// It helps to make an idempotent request. Useful for retries of request.
	RequestID string `json:"requestID"`
}

type AddUserResponse struct {
	UserID uint64 `json:"userID"`
	Token  string `json:"token"`
}
