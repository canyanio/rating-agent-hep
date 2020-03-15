package server

import (
	"context"
	"net"

	"github.com/mendersoftware/go-lib-micro/log"
	"github.com/pkg/errors"
)

func (s *UDPServer) handle(ctx context.Context, pc net.PacketConn, addr net.Addr, packet []byte) {
	l := log.FromContext(ctx)
	msg, err := s.processor.Process(packet)
	if err != nil {
		l.Error(errors.Wrap(err, "unable to decode the HEP package"))
		return
	}
	l.Infof("Message from %s, to %s, call-id %s", msg.From.User, msg.To.User, msg.CallId.Src)
}
