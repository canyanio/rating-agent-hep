package processor

import (
	"github.com/sipcapture/heplify-server/decoder"

	"github.com/canyanio/rating-agent-hep/model"
)

// HEPProcessorInterface is the interface for Server objects
type HEPProcessorInterface interface {
	Process(packet []byte) (*model.SIPMessage, error)
}

// HEPProcessor is the HEP processor
type HEPProcessor struct {
}

// NewHEPProcessor initializes a new HEP processor
func NewHEPProcessor() *HEPProcessor {
	return &HEPProcessor{}
}

// Process raw bytes containing a HEP packet
func (s *HEPProcessor) Process(packet []byte) (*model.SIPMessage, error) {
	hepPacket, err := s.hepFromBytes(packet)
	if err != nil {
		return nil, err
	}
	return model.SIPMessageFromHEP(hepPacket), nil
}

func (s *HEPProcessor) hepFromBytes(packet []byte) (*decoder.HEP, error) {
	hepPacket, err := decoder.DecodeHEP(packet)
	if err != nil {
		return nil, err
	}
	return hepPacket, nil
}
