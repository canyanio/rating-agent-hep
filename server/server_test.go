package server

import (
	"fmt"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/canyanio/rating-agent-hep/model"
	mock_processor "github.com/canyanio/rating-agent-hep/processor/mock"
)

func getFreeUDPPort() (int, error) {
	addr, err := net.ResolveUDPAddr("udp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenUDP("udp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.LocalAddr().(*net.UDPAddr).Port, nil
}

func TestNewUDPServer(t *testing.T) {
	srv := NewUDPServer()
	assert.NotNil(t, srv)
}

func TestUDPServerStart(t *testing.T) {
	buff := []byte("sample string")

	msg := &model.SIPMessage{}

	// mock processor, to check the received bytes
	mockProcessor := &mock_processor.HEPProcessor{}
	mockProcessor.On("Process",
		mock.MatchedBy(func(packet []byte) bool {
			return reflect.DeepEqual(packet, buff)
		}),
	).Return(msg, nil)

	// new UDP server with mocked processor
	srv := NewUDPServer()
	assert.NotNil(t, srv)

	srv.setProcessor(mockProcessor)

	// get a free UDP port
	udpPort, err := getFreeUDPPort()
	assert.Nil(t, err)
	listen := fmt.Sprintf("localhost:%d", udpPort)
	srv.setListen(listen)

	// start the server
	go srv.Start()
	defer srv.Stop()

	// connect to the server
	raddr, err := net.ResolveUDPAddr("udp", listen)
	assert.Nil(t, err)
	conn, err := net.DialUDP("udp", nil, raddr)
	assert.Nil(t, err)
	defer conn.Close()

	// write the buff
	bytes, err := conn.Write(buff)
	assert.Nil(t, err)
	assert.Equal(t, 13, bytes)

	// assert expectations (processor)
	time.Sleep(100 * time.Millisecond)
	mockProcessor.AssertExpectations(t)
}
