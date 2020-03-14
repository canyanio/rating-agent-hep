package server

import (
	"log"
	"net"

	"github.com/mendersoftware/go-lib-micro/config"
	"github.com/pkg/errors"

	"github.com/canyanio/rating-agent-hep/client/rabbitmq"
	dconfig "github.com/canyanio/rating-agent-hep/config"
	"github.com/canyanio/rating-agent-hep/processor"
)

// UDPServerInterface is the interface for Server objects
type UDPServerInterface interface {
	Start() error
}

// UDPServer is the UDP server
type UDPServer struct {
	processor processor.HEPProcessorInterface
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
	p := processor.NewHEPProcessor()
	quit := make(chan interface{})
	listen := config.Config.GetString(dconfig.SettingListen)
	return &UDPServer{
		processor: p,
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
	log.Printf("Listening on %v", s.listen)
	pc, err := net.ListenPacket("udp", s.listen)
	if err != nil {
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
			go s.serve(pkt.pc, pkt.addr, pkt.payload)
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
}

func (s *UDPServer) serve(pc net.PacketConn, addr net.Addr, packet []byte) {
	msg, err := s.processor.Process(packet)
	if err != nil {
		log.Print(errors.Wrap(err, "unable to decode the HEP package"))
	}
	log.Printf("Message from %s, to %s, call-id %s", msg.From.User, msg.To.User, msg.CallId.Src)
}
