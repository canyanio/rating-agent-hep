package model

// BeginTransaction is the begin transaction message
type BeginTransaction struct {
	Request BeginTransactionRequest `json:"request"`
}

// BeginTransactionRequest is the begin transaction request object
type BeginTransactionRequest struct {
	Tenant                string   `json:"tenant"`
	TransactionTag        string   `json:"transaction_tag"`
	AccountTag            string   `json:"account_tag"`
	DestinationAccountTag string   `json:"destination_account_tag"`
	Source                string   `json:"source"`
	Destination           string   `json:"destination"`
	ProductTag            string   `json:"product_tag,omitempty"`
	Tags                  []string `json:"tags,omitempty"`
	TimestampBegin        string   `json:"timestamp_begin"`
}

// EndTransaction is the begin transaction message
type EndTransaction struct {
	Request EndTransactionRequest `json:"request"`
}

// EndTransactionRequest is the begin transaction request object
type EndTransactionRequest struct {
	Tenant                string `json:"tenant"`
	TransactionTag        string `json:"transaction_tag"`
	AccountTag            string `json:"account_tag"`
	DestinationAccountTag string `json:"destination_account_tag"`
	TimestampEnd          string `json:"timestamp_end"`
}
