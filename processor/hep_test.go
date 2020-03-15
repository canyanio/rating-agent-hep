package processor

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHEPProcessor(t *testing.T) {
	srv := NewHEPProcessor()
	assert.NotNil(t, srv)
}

func TestProcess(t *testing.T) {
	srv := NewHEPProcessor()

	cwd, _ := os.Getwd()
	path := filepath.Join(cwd, "..", "testdata", "hep-invite.bin")
	packet, _ := ioutil.ReadFile(path)

	_, err := srv.Process(packet)
	assert.Nil(t, err)
}

func TestProcessInvalid(t *testing.T) {
	srv := NewHEPProcessor()

	packet := []byte{}

	_, err := srv.Process(packet)
	assert.Error(t, err)
}

func TestHepFromBytesInvalid(t *testing.T) {
	srv := NewHEPProcessor()

	packet := []byte{}

	hepPacket, err := srv.hepFromBytes(packet)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "proto: required field \"Version\" not set")
	assert.Nil(t, hepPacket)
}

func TestHepFromBytes(t *testing.T) {
	srv := NewHEPProcessor()

	cwd, _ := os.Getwd()
	path := filepath.Join(cwd, "..", "testdata", "hep-invite.bin")
	packet, _ := ioutil.ReadFile(path)

	hepPacket, err := srv.hepFromBytes(packet)
	assert.Nil(t, err)
	assert.NotNil(t, hepPacket)
	assert.Equal(t, uint32(2), hepPacket.Version)
	assert.Equal(t, uint32(0x11), hepPacket.Protocol)
	assert.Equal(t, "192.168.192.2", hepPacket.SrcIP)
	assert.Equal(t, uint32(0x13c4), hepPacket.SrcPort)
	assert.Equal(t, "192.168.192.5", hepPacket.DstIP)
	assert.Equal(t, uint32(0x13c4), hepPacket.DstPort)

	expectedPayload := "INVITE sip:service@192.168.192.5:5060 SIP/2.0\r\n" +
		"Via: SIP/2.0/UDP 192.168.192.2:5060;branch=z9hG4bK-18-1-0\r\n" +
		"From: sipp <sip:1000@192.168.192.2:5060>;tag=1\r\n" +
		"To: sut <sip:39040123456@anotherdomain.com:5060>\r\n" +
		"Call-ID: 1-18@192.168.192.2\r\n" +
		"CSeq: 1 INVITE\r\n" +
		"Contact: sip:1000@192.168.192.2:5060\r\n" +
		"Max-Forwards: 70\r\n" +
		"Subject: Test\r\n" +
		"Content-Type: application/sdp\r\n" +
		"Content-Length:   137\r\n\r\n" +
		"v=0\r\n" +
		"o=user1 53655765 2353687637 IN IP4 192.168.192.2\r\n" +
		"s=-\r\n" +
		"c=IN IP4 192.168.192.2\r\n" +
		"t=0 0\r\n" +
		"m=audio 6000 RTP/AVP 0\r\n" +
		"a=rtpmap:0 PCMU/8000\r\n"
	assert.Equal(t, expectedPayload, hepPacket.Payload)
}
