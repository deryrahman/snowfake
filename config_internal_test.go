package snowfake

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestInit(t *testing.T) {
	name := func(nodeBits, seqBits uint8) string {
		return fmt.Sprintf("nodeBits=%d.stepBits=%d", nodeBits, seqBits)
	}
	tests := []struct {
		epoch    uint64
		nodeBits uint8
		seqBits  uint8

		expectedEpoch     uint64
		expectedTimeShift uint8
		expectedNodeShift uint8
		expectedSeqShift  uint8
		expectedTimeMask  uint64
		expectedNodeMask  uint64
		expectedSeqMask   uint64
		expectedMaxNode   uint64
		expectedErr       error
	}{
		{
			epoch:    uint64(1577836800),
			nodeBits: 16,
			seqBits:  16,

			expectedEpoch:     uint64(1577836800),
			expectedTimeShift: uint8(32),
			expectedNodeShift: uint8(16),
			expectedSeqShift:  uint8(0),
			expectedTimeMask:  0xFFFFFFFF00000000,
			expectedNodeMask:  0xFFFF0000,
			expectedSeqMask:   0xFFFF,
			expectedMaxNode:   1 << 16,
			expectedErr:       nil,
		},
		{
			epoch:    uint64(1577836800),
			nodeBits: 8,
			seqBits:  4,

			expectedEpoch:     uint64(1577836800),
			expectedTimeShift: uint8(12),
			expectedNodeShift: uint8(4),
			expectedSeqShift:  uint8(0),
			expectedTimeMask:  0xFFFFFFFF000,
			expectedNodeMask:  0xFF0,
			expectedSeqMask:   0xF,
			expectedMaxNode:   1 << 8,
			expectedErr:       nil,
		},
		{
			epoch:    uint64(1577836800),
			nodeBits: 2,
			seqBits:  0,

			expectedEpoch:     uint64(1577836800),
			expectedTimeShift: uint8(2),
			expectedNodeShift: uint8(0),
			expectedSeqShift:  uint8(0),
			expectedTimeMask:  0b1111111111111111111111111111111100,
			expectedNodeMask:  0b11,
			expectedSeqMask:   0b0,
			expectedMaxNode:   1 << 2,
			expectedErr:       nil,
		},
		{
			epoch:       uint64(1577836800),
			nodeBits:    17,
			seqBits:     16,
			expectedErr: errors.New("nodeBits + seqBits should has 32 in total"),
		},
	}

	for _, tt := range tests {
		t.Run(name(tt.nodeBits, tt.seqBits), func(t *testing.T) {
			SetEpoch(tt.epoch)
			SetNodeBits(tt.nodeBits)
			SetSeqBits(tt.seqBits)
			err := Init()
			assertEqual(t, tt.expectedErr, err)
			if err == nil {
				assertEqual(t, tt.expectedEpoch, epoch)
				assertEqual(t, tt.expectedTimeShift, timeShift)
				assertEqual(t, tt.expectedNodeShift, nodeShift)
				assertEqual(t, tt.expectedSeqShift, seqShift)
				assertEqual(t, tt.expectedTimeMask, timeMask)
				assertEqual(t, tt.expectedNodeMask, nodeMask)
				assertEqual(t, tt.expectedSeqMask, seqMask)
				assertEqual(t, tt.expectedMaxNode, maxNode)
			}
		})
	}
}

func assertEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()
	if !isEqual(expected, actual) {
		t.Errorf("got %v, expected %v", actual, expected)
	}
}

func isEqual(expected, actual interface{}) bool {
	if expected == nil || actual == nil {
		return expected == actual
	}

	exp, ok := expected.([]byte)
	if !ok {
		return reflect.DeepEqual(expected, actual)
	}

	act, ok := actual.([]byte)
	if !ok {
		return false
	}
	if exp == nil || act == nil {
		return exp == nil && act == nil
	}
	return bytes.Equal(exp, act)
}
