package model

import (
	"time"

	"github.com/mendersoftware/go-lib-micro/config"

	dconfig "github.com/canyanio/rating-agent-hep/config"
	"github.com/sipcapture/heplify-server/decoder"
	"github.com/sipcapture/heplify-server/sipparser"
)

// SIPMessage represents a SIP message extracted from an HEP message
type SIPMessage struct {
	*sipparser.SipMsg
	AccountTag            string
	DestinationAccountTag string
	Timestamp             time.Time
}

// SIPMessageFromHEP returns a HEPMessage from a decoded HEP packet
func SIPMessageFromHEP(hep *decoder.HEP) *SIPMessage {
	return parseSIPMessage(hep.Payload, hep.Timestamp)
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func parseSIPMessage(payload string, timestamp time.Time) *SIPMessage {
	sipHeaderCaller := config.Config.GetString(dconfig.SettingSIPHeaderCaller)
	sipHeaderCallee := config.Config.GetString(dconfig.SettingSIPHeaderCallee)
	sipLocalDomains := config.Config.GetStringSlice(dconfig.SettingSIPLocalDomains)
	return parseSIPMessageWithSettings(payload, timestamp, sipHeaderCaller, sipHeaderCallee, sipLocalDomains)
}

func parseSIPMessageWithSettings(payload string, timestamp time.Time, sipHeaderCaller string, sipHeaderCallee string, sipLocalDomains []string) *SIPMessage {
	customHeaders := []string{}
	if sipHeaderCaller != "" {
		customHeaders = append(customHeaders, sipHeaderCaller)
	}
	if sipHeaderCallee != "" {
		customHeaders = append(customHeaders, sipHeaderCallee)
	}

	msg := sipparser.ParseMsg(payload, []string{}, customHeaders)

	accountTag := ""
	if sipHeaderCaller != "" && msg.CustomHeader[sipHeaderCaller] != "" {
		accountTag = msg.CustomHeader[sipHeaderCaller]
	} else if msg.PAssertedId != nil && (sipLocalDomains == nil || stringInSlice(msg.PaiHost, sipLocalDomains)) {
		accountTag = msg.PaiUser
	} else if sipLocalDomains != nil && stringInSlice(msg.FromHost, sipLocalDomains) {
		accountTag = msg.FromUser
	}

	destinationAccountTag := ""
	if sipHeaderCallee != "" && msg.CustomHeader[sipHeaderCallee] != "" {
		destinationAccountTag = msg.CustomHeader[sipHeaderCallee]
	} else if sipLocalDomains != nil && stringInSlice(msg.ToHost, sipLocalDomains) {
		destinationAccountTag = msg.ToUser
	}

	sipMessage := SIPMessage{
		SipMsg:                msg,
		AccountTag:            accountTag,
		DestinationAccountTag: destinationAccountTag,
		Timestamp:             timestamp,
	}
	return &sipMessage
}
