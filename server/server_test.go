package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/canyanio/rating-agent-hep/client/rabbitmq"
	mock_rabbitmq "github.com/canyanio/rating-agent-hep/client/rabbitmq/mock"
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
	cwd, _ := os.Getwd()
	path := filepath.Join(cwd, "..", "testdata", "hep-invite.bin")
	buffInvite, _ := ioutil.ReadFile(path)
	path = filepath.Join(cwd, "..", "testdata", "hep-bye.bin")
	buffBye, _ := ioutil.ReadFile(path)

	// mock rabbitmq client
	mockClient := &mock_rabbitmq.Client{}
	mockClient.On("Connect",
		mock.MatchedBy(func(_ context.Context) bool {
			return true
		}),
	).Return(nil)
	mockClient.On("Close",
		mock.MatchedBy(func(_ context.Context) bool {
			return true
		}),
	).Return(nil)
	mockClient.On("Publish",
		mock.MatchedBy(func(_ context.Context) bool {
			return true
		}),
		rabbitmq.QueueNameBeginTransaction,
		mock.AnythingOfType("*model.BeginTransactionRequest"),
	).Return(nil)
	mockClient.On("Publish",
		mock.MatchedBy(func(_ context.Context) bool {
			return true
		}),
		rabbitmq.QueueNameEndTransaction,
		mock.AnythingOfType("*model.EndTransactionRequest"),
	).Return(nil)

	// new UDP server with mocked processor
	srv := NewUDPServer()
	assert.NotNil(t, srv)

	srv.setClient(mockClient)

	// get a free UDP port
	udpPort, err := getFreeUDPPort()
	assert.Nil(t, err)
	listen := fmt.Sprintf("localhost:%d", udpPort)
	srv.setListen(listen)

	// start the server
	go srv.Start()

	// wait the server to start-up
	time.Sleep(100 * time.Millisecond)

	// connect to the server
	raddr, err := net.ResolveUDPAddr("udp", listen)
	assert.Nil(t, err)
	conn, err := net.DialUDP("udp", nil, raddr)
	assert.Nil(t, err)
	defer conn.Close()

	// write the buffInvite
	bytes, err := conn.Write(buffInvite)
	assert.Nil(t, err)
	assert.Equal(t, len(buffInvite), bytes)

	// write the buffBye
	bytes, err = conn.Write(buffBye)
	assert.Nil(t, err)
	assert.Equal(t, len(buffBye), bytes)

	// wait the server to process the packet, then shut it down
	time.Sleep(100 * time.Millisecond)
	srv.Stop()

	// assert expectations (processor)
	mockClient.AssertExpectations(t)
}
