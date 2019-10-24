package message

// Item one record of todo list
type Item struct {
	ID   uint64 `json:"id"`
	Text string `json:"text"`
}

// RequestID is an important field for deduplication requests.
// It helps to make an idempotent request. Useful for retries of request.
type AddTodoRequest struct {
	RequestID string `json:"requestID"`
	Text      string `json:"text"`
}

type AddTodoResponse struct {
	ID uint64 `json:"id"`
}

// RequestID is an important field for deduplication requests.
// It helps to make an idempotent request. Useful for retries of request.
type ResolveTodoRequest struct {
	RequestID string `json:"requestID"`
	ID        uint64 `json:"id"`
}

type ResolveTodoResponse struct {
	Status string `json:"status"`
}
