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
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/canyanio/rating-agent-hep/client/rabbitmq"
	mock_rabbitmq "github.com/canyanio/rating-agent-hep/client/rabbitmq/mock"
	dconfig "github.com/canyanio/rating-agent-hep/config"
	"github.com/canyanio/rating-agent-hep/model"
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

func getFreeTCPPort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func TestNewServer(t *testing.T) {
	srv := NewServer()
	assert.NotNil(t, srv)
}

func TestServerStartWithoutListenTCPorListenUDP(t *testing.T) {
	srv := NewServer()
	assert.NotNil(t, srv)

	srv.setListenTCP("")
	srv.setListenUDP("")

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
	srv.setClient(mockClient)

	err := srv.Start()
	assert.NotNil(t, err)
}

func TestServerStartTCP(t *testing.T) {
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
		mock.MatchedBy(func(req *model.BeginTransaction) bool {
			assert.Equal(t, "2020-03-14T08:56:08Z", req.Request.TimestampBegin)

			return true
		}),
	).Return(nil)
	mockClient.On("Publish",
		mock.MatchedBy(func(_ context.Context) bool {
			return true
		}),
		rabbitmq.QueueNameEndTransaction,
		mock.MatchedBy(func(req *model.EndTransaction) bool {
			assert.Equal(t, "2020-03-14T08:56:09Z", req.Request.TimestampEnd)

			return true
		}),
	).Return(nil)

	// new UDP server with mocked client
	srv := NewServer()
	assert.NotNil(t, srv)

	srv.setClient(mockClient)

	// get a free UDP port
	tcpPort, err := getFreeTCPPort()
	assert.Nil(t, err)
	listen := fmt.Sprintf("localhost:%d", tcpPort)
	srv.setListenTCP(listen)

	// start the server
	go srv.Start()

	// wait the server to start-up
	time.Sleep(100 * time.Millisecond)

	// connect to the server
	raddr, err := net.ResolveTCPAddr("tcp", listen)
	assert.Nil(t, err)
	conn, err := net.DialTCP("tcp", nil, raddr)
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

func TestServerStart(t *testing.T) {
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
		mock.MatchedBy(func(req *model.BeginTransaction) bool {
			assert.Equal(t, "2020-03-14T08:56:08Z", req.Request.TimestampBegin)

			return true
		}),
	).Return(nil)
	mockClient.On("Publish",
		mock.MatchedBy(func(_ context.Context) bool {
			return true
		}),
		rabbitmq.QueueNameEndTransaction,
		mock.MatchedBy(func(req *model.EndTransaction) bool {
			assert.Equal(t, "2020-03-14T08:56:09Z", req.Request.TimestampEnd)

			return true
		}),
	).Return(nil)

	// new UDP server with mocked client
	srv := NewServer()
	assert.NotNil(t, srv)

	srv.setClient(mockClient)

	// get a free UDP port
	udpPort, err := getFreeUDPPort()
	assert.Nil(t, err)
	listen := fmt.Sprintf("localhost:%d", udpPort)
	srv.setListenUDP(listen)

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

func TestServerStartWithRedis(t *testing.T) {
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

	srv := newServerWithConfig(
		listen,
		"",
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

func TestServerStartPublishFailure(t *testing.T) {
	cwd, _ := os.Getwd()

	path := filepath.Join(cwd, "..", "testdata", "hep-invite.bin")
	buffInvite, _ := ioutil.ReadFile(path)

	path = filepath.Join(cwd, "..", "testdata", "hep-ack.bin")
	buffAck, _ := ioutil.ReadFile(path)

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
		mock.MatchedBy(func(req *model.BeginTransaction) bool {
			assert.Equal(t, "2020-03-14T08:56:08Z", req.Request.TimestampBegin)

			return true
		}),
	).Return(errors.New("generic error"))

	// new UDP server with mocked client
	srv := NewServer()
	assert.NotNil(t, srv)

	srv.setClient(mockClient)

	// get a free UDP port
	udpPort, err := getFreeUDPPort()
	assert.Nil(t, err)
	listen := fmt.Sprintf("localhost:%d", udpPort)
	srv.setListenUDP(listen)

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

	// wait the server to process the packet, then shut it down
	srv.Stop()

	// assert expectations (processor)
	mockClient.AssertExpectations(t)
}

func TestServerStartAckWithoutInvite(t *testing.T) {
	cwd, _ := os.Getwd()

	path := filepath.Join(cwd, "..", "testdata", "hep-ack.bin")
	buffAck, _ := ioutil.ReadFile(path)

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

	// new UDP server with mocked client
	srv := NewServer()
	assert.NotNil(t, srv)

	srv.setClient(mockClient)

	// get a free UDP port
	udpPort, err := getFreeUDPPort()
	assert.Nil(t, err)
	listen := fmt.Sprintf("localhost:%d", udpPort)
	srv.setListenUDP(listen)

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

	// write the bufAck
	bytes, err := conn.Write(buffAck)
	assert.Nil(t, err)
	assert.Equal(t, len(buffAck), bytes)
	time.Sleep(50 * time.Millisecond)

	// wait the server to process the packet, then shut it down
	srv.Stop()

	// assert expectations (processor)
	mockClient.AssertExpectations(t)
}

func TestServerStartAckErrorInProcessing(t *testing.T) {
	buff := []byte("DUMMY")

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

	// new UDP server with mocked client
	srv := NewServer()
	assert.NotNil(t, srv)

	srv.setClient(mockClient)

	// get a free UDP port
	udpPort, err := getFreeUDPPort()
	assert.Nil(t, err)
	listen := fmt.Sprintf("localhost:%d", udpPort)
	srv.setListenUDP(listen)

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

	// write the bufAck
	bytes, err := conn.Write(buff)
	assert.Nil(t, err)
	assert.Equal(t, len(buff), bytes)
	time.Sleep(50 * time.Millisecond)

	// wait the server to process the packet, then shut it down
	srv.Stop()

	// assert expectations (processor)
	mockClient.AssertExpectations(t)
}
