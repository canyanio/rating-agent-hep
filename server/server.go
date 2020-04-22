package server

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"

	"github.com/mendersoftware/go-lib-micro/config"
	"github.com/mendersoftware/go-lib-micro/log"
	"golang.org/x/sys/unix"

	"github.com/canyanio/rating-agent-hep/client/rabbitmq"
	dconfig "github.com/canyanio/rating-agent-hep/config"
	"github.com/canyanio/rating-agent-hep/processor"
	"github.com/canyanio/rating-agent-hep/state"
)

// Interface is the interface for Server objects
type Interface interface {
	Start() error
}

// Server is the UDP/TCP server
type Server struct {
	processor processor.HEPProcessorInterface
	state     state.ManagerInterface
	client    rabbitmq.ClientInterface
	listenUDP string
	listenTCP string
	quit      chan os.Signal
}

// UDP/TCP packet received by the UDP/TCP server
type packet struct {
	addr    net.Addr
	payload []byte
}

// NewServer initializes a new UDP/TCP server
func NewServer() *Server {
	listenUDP := config.Config.GetString(dconfig.SettingListenUDP)
	listenTCP := config.Config.GetString(dconfig.SettingListenTCP)
	messagebusURI := config.Config.GetString(dconfig.SettingMessageBusURI)
	stateManagerType := config.Config.GetString(dconfig.SettingStateManager)
	redisAddress := config.Config.GetString(dconfig.SettingRedisAddress)
	redisPassword := config.Config.GetString(dconfig.SettingRedisPassword)
	redisDb := config.Config.GetInt(dconfig.SettingRedisDb)
	return newServerWithConfig(
		listenUDP,
		listenTCP,
		messagebusURI,
		stateManagerType,
		redisAddress,
		redisPassword,
		redisDb,
	)
}

func newServerWithConfig(listenUDP, listenTCP, messagebusURI, stateManagerType, redisAddress, redisPassword string, redisDb int) *Server {
	var stateManager state.ManagerInterface
	if stateManagerType == dconfig.StateManagerRedis {
		stateManager = state.NewRedisManager(redisAddress, redisPassword, redisDb)
	} else {
		stateManager = state.NewMemoryManager()
	}

	quit := make(chan os.Signal)
	processor := processor.NewHEPProcessor()
	client := rabbitmq.NewClient(messagebusURI)

	signal.Notify(quit, unix.SIGINT, unix.SIGTERM)

	return &Server{
		processor: processor,
		client:    client,
		state:     stateManager,
		quit:      quit,
		listenUDP: listenUDP,
		listenTCP: listenTCP,
	}
}

func (s *Server) setListenUDP(listen string) {
	s.listenUDP = listen
}

func (s *Server) setListenTCP(listen string) {
	s.listenTCP = listen
}

func (s *Server) setClient(c rabbitmq.ClientInterface) {
	s.client = c
}

// Start starts the UDP/TCP server which receives the HEP packats
func (s *Server) Start() error {
	ctx := context.Background()
	l := log.FromContext(ctx)

	if err := s.state.Connect(ctx); err != nil {
		l.Error(err)
		return err
	}
	defer s.state.Close(ctx)

	if err := s.client.Connect(ctx); err != nil {
		l.Error(err)
		return err
	}
	defer s.client.Close(ctx)

	listenUDP := s.listenUDP
	var pc net.PacketConn
	if listenUDP != "" {
		l.Infof("Listening on udp:%v", listenUDP)
		var err error
		pc, err = net.ListenPacket("udp", listenUDP)
		if err != nil {
			l.Error(err)
			return err
		}
	}

	listenTCP := s.listenTCP
	var li net.Listener
	if listenTCP != "" {
		l.Infof("Listening on tcp:%v", listenTCP)
		var err error
		li, err = net.Listen("tcp", listenTCP)
		if err != nil {
			l.Error(err)
			return err
		}
	}

	if listenUDP == "" && listenTCP == "" {
		err := errors.New("neither listen_tcp nor listen_udp are set, exiting")
		l.Error(err)
		return err
	}

	packets := make(chan packet)

	if listenUDP != "" {
		go func() {
			for {
				buf := make([]byte, 65536)
				n, addr, err := pc.ReadFrom(buf)
				if err != nil {
					continue
				}
				packets <- packet{
					addr:    addr,
					payload: buf[:n],
				}
			}
		}()
	}

	if listenTCP != "" {
		go func() {
			for {
				conn, err := li.Accept()
				if err != nil {
					continue
				}
				go func() {
					for {
						buf := make([]byte, 65536)
						n, err := conn.Read(buf)
						if err != nil {
							return
						}
						packets <- packet{
							addr:    conn.RemoteAddr(),
							payload: buf[:n],
						}
					}
				}()
			}
		}()
	}

	for {
		select {
		case pkt := <-packets:
			go s.handle(ctx, pkt.addr, pkt.payload)
			break

		case <-s.quit:
			if listenUDP != "" {
				pc.Close()
			}
			if listenTCP != "" {
				li.Close()
			}
			return nil
		}
	}
}

// Stop stops the UDP/TCP server
func (s *Server) Stop() {
	close(s.quit)

	ctx := context.Background()
	s.client.Close(ctx)
}
