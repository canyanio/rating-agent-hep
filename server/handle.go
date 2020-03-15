package server

import (
	"context"
	"net"
	"time"

	"github.com/mendersoftware/go-lib-micro/config"
	"github.com/mendersoftware/go-lib-micro/log"
	"github.com/pkg/errors"

	"github.com/canyanio/rating-agent-hep/client/rabbitmq"
	dconfig "github.com/canyanio/rating-agent-hep/config"
	"github.com/canyanio/rating-agent-hep/model"
)

// Server handler specific constants
const (
	MethodInvite = "INVITE"
	MethodBye    = "BYE"
)

func (s *UDPServer) handle(ctx context.Context, pc net.PacketConn, addr net.Addr, packet []byte) {
	l := log.FromContext(ctx)
	msg, err := s.processor.Process(packet)
	if err != nil {
		l.Error(errors.Wrap(err, "unable to decode the HEP package"))
		return
	}

	var routingKey string
	var req interface{}
	if string(msg.Req.Method) == MethodInvite {
		routingKey = rabbitmq.QueueNameBeginTransaction
		req = &model.BeginTransactionRequest{
			Tenant:                config.Config.GetString(dconfig.SettingTenant),
			TransactionTag:        string(msg.CallId.Src),
			AccountTag:            string(msg.From.User),
			DestinationAccountTag: string(msg.To.User),
			Source:                "sip:" + string(msg.From.User) + "@" + string(msg.From.Host),
			Destination:           "sip:" + string(msg.To.User) + "@" + string(msg.To.Host),
			TimestampBegin:        msg.Timestamp.Format(time.RFC3339),
		}
	} else if string(msg.Req.Method) == MethodBye {
		routingKey = rabbitmq.QueueNameEndTransaction
		req = &model.EndTransactionRequest{
			Tenant:                config.Config.GetString(dconfig.SettingTenant),
			TransactionTag:        string(msg.CallId.Src),
			AccountTag:            string(msg.From.User),
			DestinationAccountTag: string(msg.To.User),
			TimestampEnd:          msg.Timestamp.Format(time.RFC3339),
		}
	}

	if req != nil {
		err = s.client.Publish(ctx, routingKey, req)
		if err != nil {
			l.Error(errors.Wrap(err, "unable to publish the request"))
			return
		}
	}
}
