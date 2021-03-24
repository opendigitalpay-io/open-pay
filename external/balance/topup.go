package balance

type TryTopUpRequest struct {
	UserID   uint64                 `json:"userId"`
	Amount   int64                  `json:"amount"`
	Currency string                 `json:"currency"`
	Metadata map[string]interface{} `json:"metadata"`
}

type TryTopUpResponse struct {
	ID uint64 `json:"id"`
}

type CommitTopUpRequest struct {
	ParentID uint64                 `json:"parentId"`
	Metadata map[string]interface{} `json:"metadata"`
}

type CancelTopUpRequest struct {
	ParentID uint64                 `json:"parentId"`
	Metadata map[string]interface{} `json:"metadata"`
}
