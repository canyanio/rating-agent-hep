package server

import (
	"context"
	"net"

	"github.com/mendersoftware/go-lib-micro/config"
	"github.com/mendersoftware/go-lib-micro/log"

	"github.com/canyanio/rating-agent-hep/client/rabbitmq"
	dconfig "github.com/canyanio/rating-agent-hep/config"
	"github.com/canyanio/rating-agent-hep/processor"
	"github.com/canyanio/rating-agent-hep/state"
)

// UDPServerInterface is the interface for Server objects
type UDPServerInterface interface {
	Start() error
}

// UDPServer is the UDP server
type UDPServer struct {
	processor processor.HEPProcessorInterface
	state     state.ManagerInterface
	client    rabbitmq.ClientInterface
	listen    string
	quit      chan interface{}
}

// UDP packet received by the UDP server
type packet struct {
	pc      net.PacketConn
	addr    net.Addr
	payload []byte
}

// NewUDPServer initializes a new UDP server
func NewUDPServer() *UDPServer {
	listen := config.Config.GetString(dconfig.SettingListen)
	messagebusURI := config.Config.GetString(dconfig.SettingMessageBusURI)
	stateManagerType := config.Config.GetString(dconfig.SettingStateManager)

	var stateManager state.ManagerInterface
	if stateManagerType == dconfig.StateManagerRedis {
		redisAddress := config.Config.GetString(dconfig.SettingRedisAddress)
		redisPassword := config.Config.GetString(dconfig.SettingRedisPassword)
		redisDb := config.Config.GetInt(dconfig.SettingRedisDb)
		stateManager = state.NewRedisManager(redisAddress, redisPassword, redisDb)
	} else {
		stateManager = state.NewMemoryManager()
	}

	quit := make(chan interface{})
	processor := processor.NewHEPProcessor()
	client := rabbitmq.NewClient(messagebusURI)

	return &UDPServer{
		processor: processor,
		client:    client,
		state:     stateManager,
		quit:      quit,
		listen:    listen,
	}
}

func (s *UDPServer) setListen(listen string) {
	s.listen = listen
}

func (s *UDPServer) setProcessor(p processor.HEPProcessorInterface) {
	s.processor = p
}

func (s *UDPServer) setClient(c rabbitmq.ClientInterface) {
	s.client = c
}

// Start starts the UDP server which receives the HEP packats
func (s *UDPServer) Start() error {
	ctx := context.Background()
	l := log.FromContext(ctx)

	if err := s.client.Connect(ctx); err != nil {
		l.Error(err)
		return err
	}
	defer s.client.Close(ctx)

	l.Infof("Listening on %v", s.listen)
	pc, err := net.ListenPacket("udp", s.listen)
	if err != nil {
		l.Error(err)
		return err
	}
	defer pc.Close()

	packets := make(chan packet)
	go func() {
		for {
			buf := make([]byte, 65536)
			n, addr, err := pc.ReadFrom(buf)
			if err != nil {
				continue
			}
			packets <- packet{
				pc:      pc,
				addr:    addr,
				payload: buf[:n],
			}
		}
	}()

	for {
		select {
		case pkt := <-packets:
			go s.handle(ctx, pkt.pc, pkt.addr, pkt.payload)
			break

		case <-s.quit:
			pc.Close()
			return nil
		}
	}
}

// Stop stops the UDP server
func (s *UDPServer) Stop() {
	close(s.quit)

	ctx := context.Background()
	s.client.Close(ctx)
}
