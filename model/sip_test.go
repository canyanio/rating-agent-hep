package model

import (
	"testing"
	"time"

	"github.com/sipcapture/heplify-server/decoder"
	"github.com/stretchr/testify/assert"
)

func TestParseSIPMessage(t *testing.T) {
	payload := "INVITE sip:001234567890@10.135.0.1:5060;user=phone SIP/2.0\r\n" +
		"Via: SIP/2.0/UDP 10.135.0.12:5060;branch=z9hG4bKhye0bem20x.nx8hnt\r\n" +
		"Max-Forwards: 70\r\n" +
		"From: \"Calling User\" <sip:151@10.135.0.1:5060>;tag=m3l2hbp\r\n" +
		"To: <sip:001234567890@10.135.0.1:5060;user=phone>\r\n" +
		"Call-ID: ud04chatv9q@10.135.0.1\r\n" +
		"CSeq: 10691 INVITE\r\n" +
		"Contact: <sip:151@10.135.0.12;line=12071>;+sip.instance=\"<urn:uuid:0d9a008d-0355-0024-0004-000276f3d664>\"\r\n" +
		"P-Asserted-Identity: <sip:1000@localhost>\r\n" +
		"Allow: INVITE, CANCEL, BYE, ACK, REGISTER, OPTIONS, REFER, SUBSCRIBE, NOTIFY, MESSAGE, INFO, PRACK, UPDATE\r\n" +
		"Content-Disposition: session\r\n" +
		"Supported: replaces,100rel\r\n" +
		"User-Agent: Wildix W-AIR 03.55.00.24 9c7514340722 02:76:f3:d6:64\r\n" +
		"Content-Type: application/sdp\r\n" +
		"Content-Length: 0\r\n"
	hep := &decoder.HEP{Payload: payload}
	msg := SIPMessageFromHEP(hep)

	assert.NotNil(t, msg)

	assert.Equal(t, "151", msg.FromUser)
	assert.Equal(t, "10.135.0.1", msg.FromHost)
}

func TestParseSIPMessageWithSettings(t *testing.T) {
	var tests = []struct {
		payload               string
		accountTag            string
		destinationAccountTag string
		sipHeaderCaller       string
		sipHeaderCallee       string
		sipLocalDomains       []string
	}{
		// P-Asserted-Identity without domain filters
		{
			payload: "INVITE sip:001234567890@10.135.0.1:5060;user=phone SIP/2.0\r\n" +
				"Via: SIP/2.0/UDP 10.135.0.12:5060;branch=z9hG4bKhye0bem20x.nx8hnt\r\n" +
				"Max-Forwards: 70\r\n" +
				"From: \"Calling User\" <sip:151@10.135.0.1:5060>;tag=m3l2hbp\r\n" +
				"To: <sip:001234567890@10.135.0.1:5060;user=phone>\r\n" +
				"Call-ID: ud04chatv9q@10.135.0.1\r\n" +
				"CSeq: 10691 INVITE\r\n" +
				"Contact: <sip:151@10.135.0.12;line=12071>;+sip.instance=\"<urn:uuid:0d9a008d-0355-0024-0004-000276f3d664>\"\r\n" +
				"P-Asserted-Identity: <sip:1000@localhost>\r\n" +
				"Allow: INVITE, CANCEL, BYE, ACK, REGISTER, OPTIONS, REFER, SUBSCRIBE, NOTIFY, MESSAGE, INFO, PRACK, UPDATE\r\n" +
				"Content-Disposition: session\r\n" +
				"Supported: replaces,100rel\r\n" +
				"User-Agent: Wildix W-AIR 03.55.00.24 9c7514340722 02:76:f3:d6:64\r\n" +
				"Content-Type: application/sdp\r\n" +
				"Content-Length: 0\r\n",
			accountTag:            "1000",
			destinationAccountTag: "",
		},
		// Custom SIP header for caller
		{
			payload: "INVITE sip:001234567890@10.135.0.1:5060;user=phone SIP/2.0\r\n" +
				"Via: SIP/2.0/UDP 10.135.0.12:5060;branch=z9hG4bKhye0bem20x.nx8hnt\r\n" +
				"Max-Forwards: 70\r\n" +
				"From: \"Calling User\" <sip:151@10.135.0.1:5060>;tag=m3l2hbp\r\n" +
				"To: <sip:001234567890@10.135.0.1:5060;user=phone>\r\n" +
				"Call-ID: ud04chatv9q@10.135.0.1\r\n" +
				"CSeq: 10691 INVITE\r\n" +
				"Contact: <sip:151@10.135.0.12;line=12071>;+sip.instance=\"<urn:uuid:0d9a008d-0355-0024-0004-000276f3d664>\"\r\n" +
				"P-Asserted-Identity: <sip:1000@localhost>\r\n" +
				"X-Account-Tag: 2000\r\n" +
				"Allow: INVITE, CANCEL, BYE, ACK, REGISTER, OPTIONS, REFER, SUBSCRIBE, NOTIFY, MESSAGE, INFO, PRACK, UPDATE\r\n" +
				"Content-Disposition: session\r\n" +
				"Supported: replaces,100rel\r\n" +
				"User-Agent: Wildix W-AIR 03.55.00.24 9c7514340722 02:76:f3:d6:64\r\n" +
				"Content-Type: application/sdp\r\n" +
				"Content-Length: 0\r\n",
			accountTag:            "2000",
			destinationAccountTag: "",
			sipHeaderCaller:       "X-Account-Tag",
		},
		// Custom SIP header for callee
		{
			payload: "INVITE sip:001234567890@10.135.0.1:5060;user=phone SIP/2.0\r\n" +
				"Via: SIP/2.0/UDP 10.135.0.12:5060;branch=z9hG4bKhye0bem20x.nx8hnt\r\n" +
				"Max-Forwards: 70\r\n" +
				"From: \"Calling User\" <sip:151@10.135.0.1:5060>;tag=m3l2hbp\r\n" +
				"To: <sip:001234567890@10.135.0.1:5060;user=phone>\r\n" +
				"Call-ID: ud04chatv9q@10.135.0.1\r\n" +
				"CSeq: 10691 INVITE\r\n" +
				"Contact: <sip:151@10.135.0.12;line=12071>;+sip.instance=\"<urn:uuid:0d9a008d-0355-0024-0004-000276f3d664>\"\r\n" +
				"P-Asserted-Identity: <sip:1000@localhost>\r\n" +
				"X-Destination-Account-Tag: 2000\r\n" +
				"Allow: INVITE, CANCEL, BYE, ACK, REGISTER, OPTIONS, REFER, SUBSCRIBE, NOTIFY, MESSAGE, INFO, PRACK, UPDATE\r\n" +
				"Content-Disposition: session\r\n" +
				"Supported: replaces,100rel\r\n" +
				"User-Agent: Wildix W-AIR 03.55.00.24 9c7514340722 02:76:f3:d6:64\r\n" +
				"Content-Type: application/sdp\r\n" +
				"Content-Length: 0\r\n",
			accountTag:            "1000",
			destinationAccountTag: "2000",
			sipHeaderCallee:       "X-Destination-Account-Tag",
		},
		// Domain filter, tags from From and To headers
		{
			payload: "INVITE sip:001234567890@10.135.0.1:5060;user=phone SIP/2.0\r\n" +
				"Via: SIP/2.0/UDP 10.135.0.12:5060;branch=z9hG4bKhye0bem20x.nx8hnt\r\n" +
				"Max-Forwards: 70\r\n" +
				"From: \"Calling User\" <sip:151@10.135.0.1:5060>;tag=m3l2hbp\r\n" +
				"To: <sip:001234567890@10.135.0.1:5060;user=phone>\r\n" +
				"Call-ID: ud04chatv9q@10.135.0.1\r\n" +
				"CSeq: 10691 INVITE\r\n" +
				"Contact: <sip:151@10.135.0.12;line=12071>;+sip.instance=\"<urn:uuid:0d9a008d-0355-0024-0004-000276f3d664>\"\r\n" +
				"P-Asserted-Identity: <sip:1000@localhost>\r\n" +
				"X-Destination-Account-Tag: 2000\r\n" +
				"Allow: INVITE, CANCEL, BYE, ACK, REGISTER, OPTIONS, REFER, SUBSCRIBE, NOTIFY, MESSAGE, INFO, PRACK, UPDATE\r\n" +
				"Content-Disposition: session\r\n" +
				"Supported: replaces,100rel\r\n" +
				"User-Agent: Wildix W-AIR 03.55.00.24 9c7514340722 02:76:f3:d6:64\r\n" +
				"Content-Type: application/sdp\r\n" +
				"Content-Length: 0\r\n",
			accountTag:            "151",
			destinationAccountTag: "001234567890",
			sipLocalDomains:       []string{"10.135.0.1"},
		},
		// Domain filter, tag from P-Asserted-Identity
		{
			payload: "INVITE sip:001234567890@10.135.0.1:5060;user=phone SIP/2.0\r\n" +
				"Via: SIP/2.0/UDP 10.135.0.12:5060;branch=z9hG4bKhye0bem20x.nx8hnt\r\n" +
				"Max-Forwards: 70\r\n" +
				"From: \"Calling User\" <sip:151@10.135.0.1:5060>;tag=m3l2hbp\r\n" +
				"To: <sip:001234567890@10.135.0.1:5060;user=phone>\r\n" +
				"Call-ID: ud04chatv9q@10.135.0.1\r\n" +
				"CSeq: 10691 INVITE\r\n" +
				"Contact: <sip:151@10.135.0.12;line=12071>;+sip.instance=\"<urn:uuid:0d9a008d-0355-0024-0004-000276f3d664>\"\r\n" +
				"P-Asserted-Identity: <sip:1000@localhost>\r\n" +
				"X-Destination-Account-Tag: 2000\r\n" +
				"Allow: INVITE, CANCEL, BYE, ACK, REGISTER, OPTIONS, REFER, SUBSCRIBE, NOTIFY, MESSAGE, INFO, PRACK, UPDATE\r\n" +
				"Content-Disposition: session\r\n" +
				"Supported: replaces,100rel\r\n" +
				"User-Agent: Wildix W-AIR 03.55.00.24 9c7514340722 02:76:f3:d6:64\r\n" +
				"Content-Type: application/sdp\r\n" +
				"Content-Length: 0\r\n",
			accountTag:            "1000",
			destinationAccountTag: "",
			sipLocalDomains:       []string{"localhost"},
		},
		// No match
		{
			payload: "INVITE sip:001234567890@10.135.0.1:5060;user=phone SIP/2.0\r\n" +
				"Via: SIP/2.0/UDP 10.135.0.12:5060;branch=z9hG4bKhye0bem20x.nx8hnt\r\n" +
				"Max-Forwards: 70\r\n" +
				"From: \"Calling User\" <sip:151@10.135.0.1:5060>;tag=m3l2hbp\r\n" +
				"To: <sip:001234567890@10.135.0.1:5060;user=phone>\r\n" +
				"Call-ID: ud04chatv9q@10.135.0.1\r\n" +
				"CSeq: 10691 INVITE\r\n" +
				"Contact: <sip:151@10.135.0.12;line=12071>;+sip.instance=\"<urn:uuid:0d9a008d-0355-0024-0004-000276f3d664>\"\r\n" +
				"P-Asserted-Identity: <sip:1000@localhost>\r\n" +
				"X-Destination-Account-Tag: 2000\r\n" +
				"Allow: INVITE, CANCEL, BYE, ACK, REGISTER, OPTIONS, REFER, SUBSCRIBE, NOTIFY, MESSAGE, INFO, PRACK, UPDATE\r\n" +
				"Content-Disposition: session\r\n" +
				"Supported: replaces,100rel\r\n" +
				"User-Agent: Wildix W-AIR 03.55.00.24 9c7514340722 02:76:f3:d6:64\r\n" +
				"Content-Type: application/sdp\r\n" +
				"Content-Length: 0\r\n",
			accountTag:            "",
			destinationAccountTag: "",
			sipLocalDomains:       []string{"dummy"},
		}}

	for _, test := range tests {
		msg := parseSIPMessageWithSettings(test.payload, time.Now(),
			test.sipHeaderCaller, test.sipHeaderCallee, test.sipLocalDomains)
		assert.NotNil(t, msg)

		assert.Equal(t, test.accountTag, msg.AccountTag)
		assert.Equal(t, test.destinationAccountTag, msg.DestinationAccountTag)
	}
}
