package server

import (
	"context"
	"net"
	"strings"
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

	requestMethod := msg.FirstMethod
	callID := msg.CallID
	CSeqParts := strings.SplitN(msg.Cseq.Val, " ", 2)
	CSeqID := CSeqParts[0]
	l.Debugf("method: %s, call-id: %s", requestMethod, callID)

	var routingKey string
	var req interface{}
	if requestMethod == MethodInvite {
		call := &model.Call{
			Tenant:                config.Config.GetString(dconfig.SettingTenant),
			TransactionTag:        callID,
			AccountTag:            msg.FromUser,
			DestinationAccountTag: msg.ToUser,
			Source:                "sip:" + msg.FromUser + "@" + msg.FromHost,
			Destination:           "sip:" + msg.ToUser + "@" + msg.ToHost,
			TimestampInvite:       msg.Timestamp,
			CSeq:                  CSeqID,
		}
		s.state.Set(ctx, callID, call, StateManagerTTLInvite)
	} else if requestMethod == MethodAck {
		var call model.Call
		err := s.state.Get(ctx, callID, &call)
		if err != nil {
			l.Error(errors.Wrapf(err, "unable to retrieve status for Call-Id: %s", callID))
			return
		}

		if CSeqID == call.CSeq && call.TimestampAck.IsZero() &&
			(call.AccountTag != "" || call.DestinationAccountTag != "") {
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
				AccountTag:            msg.FromUser,
				DestinationAccountTag: msg.ToUser,
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
