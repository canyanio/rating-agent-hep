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
	MethodInvite          = "INVITE"
	MethodAck             = "ACK"
	MethodBye             = "BYE"
	StateManagerTTLInvite = 600
	StateManagerTTLCall   = 3600 * 6
)

func (s *UDPServer) handle(ctx context.Context, pc net.PacketConn, addr net.Addr, packet []byte) {
	l := log.FromContext(ctx)
	l.Debugf("received a new packet from %v, length %d bytes", addr.String(), len(packet))

	msg, err := s.processor.Process(packet)
	if err != nil {
		l.Error(errors.Wrap(err, "unable to decode the HEP package"))
		return
	}

	requestMethod := string(msg.Req.Method)
	callID := string(msg.CallId.Value)
	l.Debugf("method: %s, call-id: %s", requestMethod, callID)

	var routingKey string
	var req interface{}
	if requestMethod == MethodInvite {
		call := &model.Call{
			Tenant:                config.Config.GetString(dconfig.SettingTenant),
			TransactionTag:        callID,
			AccountTag:            string(msg.From.User),
			DestinationAccountTag: string(msg.To.User),
			Source:                "sip:" + string(msg.From.User) + "@" + string(msg.From.Host),
			Destination:           "sip:" + string(msg.To.User) + "@" + string(msg.To.Host),
			TimestampInvite:       msg.Timestamp,
			CSeq:                  string(msg.Cseq.Id),
		}
		s.state.Set(ctx, callID, call, StateManagerTTLInvite)
	} else if requestMethod == MethodAck {
		var call model.Call
		err := s.state.Get(ctx, callID, &call)
		if err != nil {
			l.Error(errors.Wrapf(err, "unable to retrieve status for Call-Id: %s", callID))
			return
		}

		if string(msg.Cseq.Id) == call.CSeq && call.TimestampAck.IsZero() {
			call.TimestampAck = msg.Timestamp
			s.state.Set(ctx, call.TransactionTag, call, StateManagerTTLCall)

			routingKey = rabbitmq.QueueNameBeginTransaction
			req = &model.BeginTransaction{
				Request: model.BeginTransactionRequest{
					Tenant:                call.Tenant,
					TransactionTag:        call.TransactionTag,
					AccountTag:            call.AccountTag,
					DestinationAccountTag: call.DestinationAccountTag,
					Source:                call.Source,
					Destination:           call.Destination,
					TimestampBegin:        msg.Timestamp.Format(time.RFC3339),
				},
			}
		}
	} else if requestMethod == MethodBye {
		s.state.Delete(ctx, callID)

		routingKey = rabbitmq.QueueNameEndTransaction
		req = &model.EndTransaction{
			Request: model.EndTransactionRequest{
				Tenant:                config.Config.GetString(dconfig.SettingTenant),
				TransactionTag:        callID,
				AccountTag:            string(msg.From.User),
				DestinationAccountTag: string(msg.To.User),
				TimestampEnd:          msg.Timestamp.Format(time.RFC3339),
			},
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
