package server

import (
	"context"
	"net"
	"strings"
	"time"

	uuid "github.com/google/uuid"
	"github.com/mendersoftware/go-lib-micro/config"
	"github.com/mendersoftware/go-lib-micro/log"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/canyanio/rating-agent-hep/client/rabbitmq"
	dconfig "github.com/canyanio/rating-agent-hep/config"
	"github.com/canyanio/rating-agent-hep/model"
)

// Server handler specific constants
const (
	MethodInvite          = "INVITE"
	MethodAck             = "ACK"
	MethodBye             = "BYE"
	MethodCancel          = "CANCEL"
	StateManagerTTLInvite = 600
	StateManagerTTLCall   = 3600 * 6
)

func (s *Server) handle(ctx context.Context, addr net.Addr, packet []byte) {
	reqID := uuid.New()

	l := log.FromContext(ctx)

	msg, err := s.processor.Process(packet)
	if err != nil {
		l.Error(errors.Wrap(err, "unable to decode the HEP package"))
		l.WithFields(logrus.Fields{
			"req-id": reqID,
			"source": addr.String(),
			"length": len(packet),
		}).Error("unable to decode the HEP package")
		return
	}

	requestMethod := msg.FirstMethod
	callID := msg.CallID
	CSeqParts := strings.SplitN(msg.Cseq.Val, " ", 2)
	CSeqID := CSeqParts[0]

	productTag := config.Config.GetString(dconfig.SettingProductTag)
	transactionTags := config.Config.GetStringSlice(dconfig.SettingTransactionTags)

	l.WithFields(logrus.Fields{
		"req-id":        reqID,
		"source":        addr.String(),
		"length":        len(packet),
		"requestMethod": requestMethod,
		"callID":        callID,
		"CSeqID":        CSeqID,
	}).Debug("received msg")

	var routingKey string
	var req interface{}
	if requestMethod == MethodInvite && (msg.AccountTag != "" || msg.DestinationAccountTag != "") {
		call := &model.Call{
			Tenant:                config.Config.GetString(dconfig.SettingTenant),
			TransactionTag:        callID,
			AccountTag:            msg.AccountTag,
			DestinationAccountTag: msg.DestinationAccountTag,
			Source:                "sip:" + msg.FromUser + "@" + msg.FromHost,
			Destination:           "sip:" + msg.ToUser + "@" + msg.ToHost,
			TimestampInvite:       msg.Timestamp,
			CSeq:                  CSeqID,
		}
		err := s.state.Set(ctx, callID, call, StateManagerTTLInvite)
		if err != nil {
			l.WithFields(logrus.Fields{
				"req-id":  reqID,
				"method":  requestMethod,
				"call-id": callID,
				"ts":      msg.Timestamp,
				"err":     err.Error(),
			}).Error("unable to set the call status")
		} else {
			l.WithFields(logrus.Fields{
				"req-id":  reqID,
				"source":  addr.String(),
				"length":  len(packet),
				"method":  requestMethod,
				"call-id": callID,
				"ts":      msg.Timestamp,
			}).Debug("call status set in the state manager, waiting for the ACK")
		}
	} else {
		var call model.Call
		err := s.state.Get(ctx, callID, &call)
		if err != nil || call.CSeq == "" {
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			l.WithFields(logrus.Fields{
				"req-id":  reqID,
				"source":  addr.String(),
				"length":  len(packet),
				"method":  requestMethod,
				"call-id": callID,
				"ts":      msg.Timestamp,
				"err":     errStr,
			}).Error("unable to retrieve the call status, INVITE has not been received for this calls")
			return
		}

		if requestMethod == MethodAck {
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
						ProductTag:            productTag,
						Tags:                  transactionTags,
						TimestampBegin:        msg.Timestamp.UTC().Format(time.RFC3339),
					},
				}

				l.WithFields(logrus.Fields{
					"req-id":  reqID,
					"source":  addr.String(),
					"length":  len(packet),
					"method":  requestMethod,
					"call-id": callID,
					"ts":      msg.Timestamp,
				}).Debug("call start detected: begin transaction")
			}
		} else if requestMethod == MethodBye || requestMethod == MethodCancel {
			s.state.Delete(ctx, callID)

			routingKey = rabbitmq.QueueNameEndTransaction
			req = &model.EndTransaction{
				Request: model.EndTransactionRequest{
					Tenant:                config.Config.GetString(dconfig.SettingTenant),
					TransactionTag:        callID,
					AccountTag:            call.AccountTag,
					DestinationAccountTag: call.DestinationAccountTag,
					TimestampEnd:          msg.Timestamp.UTC().Format(time.RFC3339),
				},
			}

			l.WithFields(logrus.Fields{
				"req-id":  reqID,
				"source":  addr.String(),
				"length":  len(packet),
				"method":  requestMethod,
				"call-id": callID,
				"ts":      msg.Timestamp,
			}).Debug("call end detected: end transaction")
		}
	}

	if req != nil {
		err = s.client.Publish(ctx, routingKey, req)
		if err != nil {
			l.WithFields(logrus.Fields{
				"req-id":  reqID,
				"source":  addr.String(),
				"length":  len(packet),
				"method":  requestMethod,
				"call-id": callID,
				"err":     err.Error(),
			}).Error("unable to publish the request")

			return
		}
	}
}
