package gateway

type CardRequest struct {
	ID            uint64
	RequestID     uint64
	ParentID      uint64
	TransferTxnID uint64
	GatewayTxnID  string
	GatewayToken  string
	Gateway       SUPPORTED_GATEWAY
	Amount        int64
	Currency      string
	RequestType   REQUEST_TYPE
	AutoCapture   bool
	Status        STATUS
	Metadata      map[string]interface{}
	CreatedAt     int64
	UpdatedAt     int64
}
