package canjson

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/blueinnovationsgroup/can-go"
	"github.com/blueinnovationsgroup/can-go/pkg/descriptor"
	"github.com/blueinnovationsgroup/can-go/pkg/generated"
)

// preAllocatedBytesPerSignal is an estimate of how many bytes each signal needs.
const preAllocatedBytesPerSignal = 40

// Marshal a CAN message to JSON.
func Marshal(m generated.Message) ([]byte, error) {
	f := m.Frame()
	bytes := make([]byte, 0, len(m.Descriptor().Signals)*preAllocatedBytesPerSignal)
	bytes = append(bytes, '{')
	for i, s := range m.Descriptor().Signals {
		s := s
		bytes = append(bytes, '"')
		bytes = append(bytes, s.Name...)
		bytes = append(bytes, `":`...)
		sig := &signal{}
		sig.set(s, f)
		jsonSig, err := json.Marshal(sig)
		if err != nil {
			return nil, fmt.Errorf("marshal json: %w", err)
		}
		bytes = append(bytes, jsonSig...)
		if i < len(m.Descriptor().Signals)-1 {
			bytes = append(bytes, ',')
		}
	}
	bytes = append(bytes, '}')
	return bytes, nil
}

type signal struct {
	Raw         json.Number
	Physical    json.Number
	Unit        string `json:",omitempty"`
	Description string `json:",omitempty"`
}

func (s *signal) set(desc *descriptor.Signal, f can.Frame) {
	switch {
	case desc.Length == 1: // bool
		s.setBoolValue(desc.UnmarshalBool(f.Data), desc)
	case desc.IsSigned: // signed
		s.setSignedValue(desc.UnmarshalSigned(f.Data), desc)
	default: // unsigned
		s.setUnsignedValue(desc.UnmarshalUnsigned(f.Data), desc)
	}
}

func (s *signal) setUnsignedValue(value uint64, desc *descriptor.Signal) {
	s.Raw = uintToJSON(value)
	s.Physical = floatToJSON(desc.ToPhysical(float64(value)))
	s.Unit = desc.Unit
	if value, ok := desc.ValueDescription(int64(value)); ok {
		s.Description = value
	}
}

func (s *signal) setSignedValue(value int64, desc *descriptor.Signal) {
	s.Raw = intToJSON(value)
	s.Physical = floatToJSON(desc.ToPhysical(float64(value)))
	s.Unit = desc.Unit
	if value, ok := desc.ValueDescription(value); ok {
		s.Description = value
	}
}

func (s *signal) setBoolValue(value bool, desc *descriptor.Signal) {
	if value {
		s.Raw = "1"
		s.Physical = floatToJSON(desc.ToPhysical(1))
	} else {
		s.Raw = "0"
		s.Physical = floatToJSON(desc.ToPhysical(0))
	}
	s.Unit = desc.Unit
	var intValue int64
	if value {
		intValue = 1
	}
	if value, ok := desc.ValueDescription(intValue); ok {
		s.Description = value
	}
}

func floatToJSON(f float64) json.Number {
	return json.Number(strconv.FormatFloat(f, 'f', -1, 64))
}

func intToJSON(i int64) json.Number {
	return json.Number(strconv.Itoa(int(i)))
}

func uintToJSON(i uint64) json.Number {
	return json.Number(strconv.Itoa(int(i)))
}
