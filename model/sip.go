package model

import (
	"github.com/marv2097/siprocket"
	"github.com/sipcapture/heplify-server/decoder"
)

// SIPMessage represents a SIP message extracted from an HEP message
type SIPMessage struct {
	siprocket.SipMsg
}

// SIPMessageFromHEP returns a HEPMessage from a decoded HEP packet
func SIPMessageFromHEP(hep *decoder.HEP) *SIPMessage {
	msg := siprocket.Parse([]byte(hep.Payload))
	return &SIPMessage{msg}
}
