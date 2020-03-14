package model

// BeginTransactionRequest is the begin transaction request object
type BeginTransactionRequest struct {
	Tenant                string `json:"tenant"`
	TransactionTag        string `json:"transaction_tag"`
	AccountTag            string `json:"account_tag"`
	DestinationAccountTAg string `json:"destination_account_tag"`
	Source                string `json:"source"`
	Destination           string `json:"destination"`
	TimestampBegin        string `json:"timestamp_begin"`
}

// EndTransactionRequest is the begin transaction request object
type EndTransactionRequest struct {
	Tenant                string `json:"tenant"`
	TransactionTag        string `json:"transaction_tag"`
	AccountTag            string `json:"account_tag"`
	DestinationAccountTAg string `json:"destination_account_tag"`
	TimestampEnd          string `json:"timestamp_begin"`
}
