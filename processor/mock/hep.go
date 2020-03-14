package mock

import (
	"github.com/stretchr/testify/mock"

	"github.com/sipcapture/heplify-server/decoder"

	"github.com/canyanio/rating-agent-hep/client/rabbitmq"
	"github.com/canyanio/rating-agent-hep/model"
)

// HEPProcessor is the HEP processor
type HEPProcessor struct {
	mock.Mock
}

// SetClient sets the message bus client for the processor
func (s *HEPProcessor) SetClient(client rabbitmq.ClientInterface) {
	s.Called(client)
}

// Process raw bytes containing a HEP packet
func (s *HEPProcessor) Process(packet []byte) (*model.SIPMessage, error) {
	ret := s.Called(packet)

	var r0 *model.SIPMessage
	if rf, ok := ret.Get(0).(func([]byte) *model.SIPMessage); ok {
		r0 = rf(packet)
	} else {
		r0 = ret.Get(0).(*model.SIPMessage)
	}

	var r1 error
	if rf, ok := ret.Get(01).(func([]byte) error); ok {
		r1 = rf(packet)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// HEPFromBytes returns a HEP packet from the raw bytes
func (s *HEPProcessor) HEPFromBytes(packet []byte) (*decoder.HEP, error) {
	ret := s.Called(packet)

	var r0 *decoder.HEP
	if rf, ok := ret.Get(0).(func([]byte) *decoder.HEP); ok {
		r0 = rf(packet)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*decoder.HEP)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(packet)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}
