package model

import (
	"time"
)

// Call stores the status of a call in the state manager
type Call struct {
	Tenant                string    `json:"tenant"`
	TransactionTag        string    `json:"transaction_tag"`
	AccountTag            string    `json:"account_tag"`
	DestinationAccountTag string    `json:"destination_account_tag"`
	Source                string    `json:"source"`
	Destination           string    `json:"destination"`
	CSeq                  string    `json:"cseq"`
	TimestampInvite       time.Time `json:"timestamp_begin"`
	TimestampAck          time.Time `json:"timestamp_ack"`
	TimestampBye          time.Time `json:"timestamp_bye"`
}
