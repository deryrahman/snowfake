package snowfake

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	node := uint64(1)

	expectedNode := uint64(1)
	expectedEpoch := uint64(1577836800)
	expectedTimeShift := uint8(32)
	expectedNodeShift := uint8(24)
	expectedTime := uint64(0)
	expectedStep := uint64(0)
	expectedTimeMask := uint64(0xFFFFFFFF00000000)
	expectedNodeMask := uint64(0xFF000000)
	expectedStepMask := uint64(0xFFFFFF)

	sf := New(node)
	assertNotNil(t, sf)

	assertEqual(t, expectedNode, sf.node)
	assertEqual(t, expectedEpoch, sf.epoch)
	assertEqual(t, expectedTime, sf.time)
	assertEqual(t, expectedStep, sf.step)
	assertEqual(t, expectedTimeShift, sf.timeShift)
	assertEqual(t, expectedNodeShift, sf.nodeShift)
	assertEqual(t, expectedTimeMask, sf.timeMask)
	assertEqual(t, expectedNodeMask, sf.nodeMask)
	assertEqual(t, expectedStepMask, sf.stepMask)
}

func TestNewWithConfig(t *testing.T) {
	name := func(nodeBits, stepBits uint8) string {
		return fmt.Sprintf("nodeBits=%d.stepBits=%d", nodeBits, stepBits)
	}
	tests := []struct {
		node     uint64
		epoch    uint64
		nodeBits uint8
		stepBits uint8

		expectedNode  uint64
		expectedEpoch uint64
		expectedTime  uint64
		expectedStep  uint64

		expectedTimeShift uint8
		expectedNodeShift uint8
		expectedTimeMask  uint64
		expectedNodeMask  uint64
		expectedStepMask  uint64
		expectedErr       error
	}{
		{
			node:     uint64(1),
			epoch:    uint64(1577836800),
			nodeBits: 16,
			stepBits: 16,

			expectedNode:  uint64(1),
			expectedEpoch: uint64(1577836800),
			expectedTime:  uint64(0),
			expectedStep:  uint64(0),

			expectedTimeShift: uint8(32),
			expectedNodeShift: uint8(16),
			expectedTimeMask:  0xFFFFFFFF00000000,
			expectedNodeMask:  0xFFFF0000,
			expectedStepMask:  0xFFFF,
			expectedErr:       nil,
		},
		{
			node:     uint64(2),
			epoch:    uint64(807148800),
			nodeBits: 8,
			stepBits: 4,

			expectedNode:  uint64(2),
			expectedEpoch: uint64(807148800),
			expectedTime:  uint64(0),
			expectedStep:  uint64(0),

			expectedTimeShift: uint8(12),
			expectedNodeShift: uint8(4),
			expectedTimeMask:  0xFFFFFFFF000,
			expectedNodeMask:  0xFF0,
			expectedStepMask:  0xF,
			expectedErr:       nil,
		},
		{
			node:     uint64(0),
			epoch:    uint64(1577836800),
			nodeBits: 2,
			stepBits: 0,

			expectedNode:  uint64(0),
			expectedEpoch: uint64(1577836800),
			expectedTime:  uint64(0),
			expectedStep:  uint64(0),

			expectedTimeShift: uint8(2),
			expectedNodeShift: uint8(0),
			expectedTimeMask:  0b1111111111111111111111111111111100,
			expectedNodeMask:  0b11,
			expectedStepMask:  0b0,
			expectedErr:       nil,
		},
		{
			node:        uint64(0),
			epoch:       uint64(1577836800),
			nodeBits:    17,
			stepBits:    16,
			expectedErr: errors.New("nodeBits + stepBits should has 32 in total, got 33"),
		},
		{
			node:        uint64(4),
			epoch:       uint64(1577836800),
			nodeBits:    2,
			stepBits:    16,
			expectedErr: errors.New("node should below 4"),
		},
	}

	for _, tt := range tests {
		t.Run(name(tt.nodeBits, tt.stepBits), func(t *testing.T) {
			sf, err := NewWithConfig(tt.node, tt.epoch, tt.nodeBits, tt.stepBits)
			assertEqual(t, tt.expectedErr, err)
			if err == nil {
				assertNotNil(t, sf)
				assertEqual(t, tt.expectedNode, sf.node)
				assertEqual(t, tt.expectedEpoch, sf.epoch)
				assertEqual(t, tt.expectedTime, sf.time)
				assertEqual(t, tt.expectedStep, sf.step)
				assertEqual(t, tt.expectedTimeShift, sf.timeShift)
				assertEqual(t, tt.expectedNodeShift, sf.nodeShift)
				assertEqual(t, tt.expectedTimeMask, sf.timeMask)
				assertEqual(t, tt.expectedNodeMask, sf.nodeMask)
				assertEqual(t, tt.expectedStepMask, sf.stepMask)
			}
		})
	}
}

func TestSnowfake_GenerateID(t *testing.T) {
	node := uint64(29)

	sf := New(node)

	estimateTimeFromID := uint64(time.Now().Unix()) - sf.epoch
	expectedNodeFromID := uint64(29)
	expectedStepFromID := uint64(0)

	assertNotNil(t, sf)
	if sf != nil {
		id := sf.GenerateID()
		assertTrue(t, estimateTimeFromID <= ((id&sf.timeMask)>>sf.timeShift))
		assertEqual(t, expectedNodeFromID, (id&sf.nodeMask)>>sf.nodeShift)
		assertEqual(t, expectedStepFromID, id&sf.stepMask)
	}

}

func TestSnowfake_GenerateID_Collision(t *testing.T) {
	node := uint64(29)
	concurrent := 10000

	sf := New(node)

	assertNotNil(t, sf)
	if sf != nil {
		var wg sync.WaitGroup
		c := make(chan uint64, concurrent)

		wg.Add(concurrent)
		for i := 0; i < concurrent; i++ {
			go func(c chan uint64, wg *sync.WaitGroup) {
				id := sf.GenerateID()
				c <- id
				wg.Done()
			}(c, &wg)
		}

		wg.Wait()
		close(c)

		mp := make(map[uint64]bool)
		for ch := range c {
			mp[ch] = true
		}

		assertEqual(t, concurrent, len(mp))
	}

}

func assertTrue(t *testing.T, cond bool) {
	if !cond {
		t.Errorf("got false, expected true")
	}
}

func assertNotNil(t *testing.T, actual interface{}) {
	t.Helper()
	if isEqual(nil, actual) {
		t.Errorf("got %v, expected nil", actual)
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
