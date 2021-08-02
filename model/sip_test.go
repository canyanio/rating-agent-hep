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
		payload                   string
		accountTag                string
		destinationAccountTag     string
		sipHeaderCaller           string
		sipHeaderCallee           string
		sipHeaderHistoryInfo      string
		sipHeaderHistoryInfoIndex int
		sipLocalDomains           []string
		accountTagMatchRegexp     string
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
		// No configuration
		{
			payload: "INVITE sip:001234567890@10.135.0.1:5060;user=phone SIP/2.0\r\n" +
				"Via: SIP/2.0/UDP 10.135.0.12:5060;branch=z9hG4bKhye0bem20x.nx8hnt\r\n" +
				"Max-Forwards: 70\r\n" +
				"From: \"Calling User\" <sip:151@10.135.0.1:5060>;tag=m3l2hbp\r\n" +
				"To: <sip:001234567890@10.135.0.1:5060;user=phone>\r\n" +
				"Call-ID: ud04chatv9q@10.135.0.1\r\n" +
				"CSeq: 10691 INVITE\r\n" +
				"Contact: <sip:151@10.135.0.12;line=12071>;+sip.instance=\"<urn:uuid:0d9a008d-0355-0024-0004-000276f3d664>\"\r\n" +
				"Allow: INVITE, CANCEL, BYE, ACK, REGISTER, OPTIONS, REFER, SUBSCRIBE, NOTIFY, MESSAGE, INFO, PRACK, UPDATE\r\n" +
				"Content-Disposition: session\r\n" +
				"Supported: replaces,100rel\r\n" +
				"User-Agent: Wildix W-AIR 03.55.00.24 9c7514340722 02:76:f3:d6:64\r\n" +
				"Content-Type: application/sdp\r\n" +
				"Content-Length: 0\r\n",
			accountTag:            "",
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
		},
		// P-Asserted-Identity with domain filters and regexp
		{
			payload: "INVITE sip:+390759975378@telecomitalia.it;user=phone SIP/2.0\r\n" +
				"Via: SIP/2.0/UDP 80.21.1.119:5060;branch=z9hG4bKvtfill3090fs7a6mvfm0.1\r\n" +
				"To: <sip:+390759975378@operatore.it;user=phone>\r\n" +
				"From: \"+39029501685\" <sip:+39029501685@ictvoip.it;user=phone>;tag=007694f3-0004-0001-0000-0000\r\n" +
				"Call-ID: 007694680076945-0004-0001-0000-0000@172.16.30.36\r\n" +
				"CSeq: 1 INVITE\r\n" +
				"Max-Forwards: 69\r\n" +
				"Contact: <sip:+39029501685@80.21.1.119:5060;transport=udp>\r\n" +
				"Allow: INVITE, ACK, PRACK, CANCEL, BYE, OPTIONS, MESSAGE, NOTIFY, UPDATE, REGISTER, INFO, REFER, SUBSCRIBE, PUBLISH\r\n" +
				"Supported: 100rel\r\n" +
				"P-Asserted-Identity: <sip:+39029501685@ictvoip.it;user=phone>\r\n" +
				"Accept: application/sdp, application/isup, application/xml\r\n" +
				"Content-Type: application/sdp\r\n" +
				"Content-Length: 259\r\n" +
				"Route: <sip:+390759975378@37.10.80.21:5060;user=phone;lr>",
			accountTag:            "029501685",
			accountTagMatchRegexp: "\\+39([0-9]+)",
			sipLocalDomains:       []string{"ictvoip.it", "sip.ictvoip.it"},
		},
		// History-Info with domain filters and regexp
		{
			payload: "INVITE sip:+390759975378@telecomitalia.it;user=phone SIP/2.0\r\n" +
				"Via: SIP/2.0/UDP 80.21.1.119:5060;branch=z9hG4bKvtfill3090fs7a6mvfm0.1\r\n" +
				"To: <sip:+390759975378@operatore.it;user=phone>\r\n" +
				"From: \"+39029501685\" <sip:+39029501685@ictvoip.it;user=phone>;tag=007694f3-0004-0001-0000-0000\r\n" +
				"Call-ID: 007694680076945-0004-0001-0000-0000@172.16.30.36\r\n" +
				"CSeq: 1 INVITE\r\n" +
				"Max-Forwards: 69\r\n" +
				"Contact: <sip:+39029501685@80.21.1.119:5060;transport=udp>\r\n" +
				"Allow: INVITE, ACK, PRACK, CANCEL, BYE, OPTIONS, MESSAGE, NOTIFY, UPDATE, REGISTER, INFO, REFER, SUBSCRIBE, PUBLISH\r\n" +
				"Supported: 100rel\r\n" +
				"P-Asserted-Identity: <sip:+39029501685@ictvoip.it;user=phone>\r\n" +
				"Accept: application/sdp, application/isup, application/xml\r\n" +
				"History-Info: <sip:+3902888804@ictvoip.it;user=phone?Privacy=history>;index=1\r\n" +
				"History-Info: <sip:+390759975378@ictvoip.it;user=phone;cause=302>;index=1.1\r\n" +
				"Content-Type: application/sdp\r\n" +
				"Content-Length: 259\r\n" +
				"Route: <sip:+390759975378@37.10.80.21:5060;user=phone;lr>",
			accountTag:                "02888804",
			accountTagMatchRegexp:     "\\+39([0-9]+)",
			sipHeaderHistoryInfo:      "History-Info",
			sipHeaderHistoryInfoIndex: 1,
			sipLocalDomains:           []string{"ictvoip.it", "sip.ictvoip.it"},
		},
	}

	for _, test := range tests {
		msg := parseSIPMessageWithSettings(test.payload, time.Now(),
			test.sipHeaderCaller, test.sipHeaderCallee, test.sipHeaderHistoryInfo,
			test.sipHeaderHistoryInfoIndex, test.sipLocalDomains, test.accountTagMatchRegexp)
		assert.NotNil(t, msg)

		assert.Equal(t, test.accountTag, msg.AccountTag)
		assert.Equal(t, test.destinationAccountTag, msg.DestinationAccountTag)
	}
}
