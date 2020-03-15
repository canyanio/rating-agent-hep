package mock

import (
	"github.com/stretchr/testify/mock"

	"github.com/canyanio/rating-agent-hep/model"
)

// HEPProcessor is the HEP processor
type HEPProcessor struct {
	mock.Mock
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
