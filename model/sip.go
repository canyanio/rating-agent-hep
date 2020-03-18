package model

import (
	"time"

	"github.com/sipcapture/heplify-server/decoder"
	"github.com/sipcapture/heplify-server/sipparser"
)

// SIPMessage represents a SIP message extracted from an HEP message
type SIPMessage struct {
	*sipparser.SipMsg
	Timestamp time.Time
}

// SIPMessageFromHEP returns a HEPMessage from a decoded HEP packet
func SIPMessageFromHEP(hep *decoder.HEP) *SIPMessage {
	msg := sipparser.ParseMsg(hep.Payload, []string{}, []string{})
	timestamp := time.Unix(0, int64(hep.GetTsec())*int64(time.Second)+int64(hep.GetTmsec())*int64(time.Millisecond))
	sipMessage := SIPMessage{
		SipMsg:    msg,
		Timestamp: timestamp,
	}
	return &sipMessage
}
