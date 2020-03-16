package server

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mendersoftware/go-lib-micro/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/canyanio/rating-agent-hep/client/rabbitmq"
	mock_rabbitmq "github.com/canyanio/rating-agent-hep/client/rabbitmq/mock"
	dconfig "github.com/canyanio/rating-agent-hep/config"
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

	path = filepath.Join(cwd, "..", "testdata", "hep-ack.bin")
	buffAck, _ := ioutil.ReadFile(path)

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
		mock.AnythingOfType("*model.BeginTransaction"),
	).Return(nil)
	mockClient.On("Publish",
		mock.MatchedBy(func(_ context.Context) bool {
			return true
		}),
		rabbitmq.QueueNameEndTransaction,
		mock.AnythingOfType("*model.EndTransaction"),
	).Return(nil)

	// new UDP server with mocked client
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
	time.Sleep(50 * time.Millisecond)

	// write the bufAck
	bytes, err = conn.Write(buffAck)
	assert.Nil(t, err)
	assert.Equal(t, len(buffAck), bytes)
	time.Sleep(50 * time.Millisecond)

	// write the buffBye
	bytes, err = conn.Write(buffBye)
	assert.Nil(t, err)
	assert.Equal(t, len(buffBye), bytes)
	time.Sleep(50 * time.Millisecond)

	// wait the server to process the packet, then shut it down
	srv.Stop()

	// assert expectations (processor)
	mockClient.AssertExpectations(t)
}

func TestUDPServerStartWithRedis(t *testing.T) {
	flag.Parse()
	if testing.Short() {
		t.Skip()
	}

	cwd, _ := os.Getwd()

	path := filepath.Join(cwd, "..", "testdata", "hep-invite.bin")
	buffInvite, _ := ioutil.ReadFile(path)

	path = filepath.Join(cwd, "..", "testdata", "hep-ack.bin")
	buffAck, _ := ioutil.ReadFile(path)

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
		mock.AnythingOfType("*model.BeginTransaction"),
	).Return(nil)
	mockClient.On("Publish",
		mock.MatchedBy(func(_ context.Context) bool {
			return true
		}),
		rabbitmq.QueueNameEndTransaction,
		mock.AnythingOfType("*model.EndTransaction"),
	).Return(nil)

	// get a free UDP port
	udpPort, err := getFreeUDPPort()
	assert.Nil(t, err)
	listen := fmt.Sprintf("localhost:%d", udpPort)

	// new UDP server with mocked processor and redis state manager
	stateManagerType := "redis"
	messagebusURI := config.Config.GetString(dconfig.SettingMessageBusURI)
	redisAddress := config.Config.GetString(dconfig.SettingRedisAddress)
	redisPassword := config.Config.GetString(dconfig.SettingRedisPassword)
	redisDb := config.Config.GetInt(dconfig.SettingRedisDb)

	srv := newUDPServerWithConfig(
		listen,
		messagebusURI,
		stateManagerType,
		redisAddress,
		redisPassword,
		redisDb,
	)
	assert.NotNil(t, srv)

	srv.setClient(mockClient)

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
	time.Sleep(50 * time.Millisecond)

	// write the bufAck
	bytes, err = conn.Write(buffAck)
	assert.Nil(t, err)
	assert.Equal(t, len(buffAck), bytes)
	time.Sleep(50 * time.Millisecond)

	// write the buffBye
	bytes, err = conn.Write(buffBye)
	assert.Nil(t, err)
	assert.Equal(t, len(buffBye), bytes)
	time.Sleep(50 * time.Millisecond)

	// wait the server to process the packet, then shut it down
	srv.Stop()

	// assert expectations (processor)
	mockClient.AssertExpectations(t)
}
