package snowfake_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/deryrahman/snowfake"
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
			expectedErr: errors.New("NodeBits + SeqBits should has 32 in total"),
		},
	}

	for _, tt := range tests {
		t.Run(name(tt.nodeBits, tt.seqBits), func(t *testing.T) {
			snowfake.Epoch = tt.epoch
			snowfake.NodeBits = tt.nodeBits
			snowfake.SeqBits = tt.seqBits
			err := snowfake.Init()
			assertEqual(t, tt.expectedErr, err)
			if err == nil {
				assertEqual(t, tt.expectedEpoch, snowfake.Epoch)
				assertEqual(t, tt.expectedTimeShift, snowfake.GetTimeShift())
				assertEqual(t, tt.expectedNodeShift, snowfake.GetNodeShift())
				assertEqual(t, tt.expectedSeqShift, snowfake.GetSeqShift())
				assertEqual(t, tt.expectedTimeMask, snowfake.GetTimeMask())
				assertEqual(t, tt.expectedNodeMask, snowfake.GetNodeMask())
				assertEqual(t, tt.expectedSeqMask, snowfake.GetSeqMask())
				assertEqual(t, tt.expectedMaxNode, snowfake.GetMaxNode())
			}
		})
	}
}
