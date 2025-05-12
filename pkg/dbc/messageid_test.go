package dbc

import (
	"fmt"
	"testing"

	"gotest.tools/v3/assert"
)

func TestMessageID_Validate(t *testing.T) {
	for _, tt := range []MessageID{
		0,
		1,
		maxID,
		0 | messageIDExtendedFlag,
		1 | messageIDExtendedFlag,
		maxID | messageIDExtendedFlag,
		maxExtendedID | messageIDExtendedFlag,
		messageIDIndependentSignals,
	} {
		t.Run(fmt.Sprintf("%v", tt), func(t *testing.T) {
			assert.NilError(t, tt.Validate())
		})
	}
}

func TestMessageID_Validate_Error(t *testing.T) {
	for _, tt := range []MessageID{
		maxID + 1,
		(maxExtendedID + 1) | messageIDExtendedFlag,
		0xffffffff,
	} {
		t.Run(fmt.Sprintf("%v", tt), func(t *testing.T) {
			assert.ErrorContains(t, tt.Validate(), "invalid")
		})
	}
}

func TestMessageID_ToCAN(t *testing.T) {
	for _, tt := range []struct {
		messageID MessageID
		expected  uint32
	}{
		{messageID: 1, expected: 1},
		{messageID: messageIDIndependentSignals, expected: 0x40000000},
		{messageID: 2566857156, expected: 419373508},
	} {
		t.Run(fmt.Sprintf("%v", tt.messageID), func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.messageID.ToCAN())
		})
	}
}

func TestMessageID_IsExtended(t *testing.T) {
	for _, tt := range []struct {
		messageID MessageID
		expected  bool
	}{
		{messageID: 1, expected: false},
		{messageID: messageIDIndependentSignals, expected: false},
		{messageID: 2566857156, expected: true},
	} {
		t.Run(fmt.Sprintf("%v", tt.messageID), func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.messageID.IsExtended())
		})
	}
}
