package snowfake

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	expectedTimeMask := uint64(0xFFFFFFFF00000000)
	expectedNodeMask := uint64(0xFF000000)
	expectedStepMask := uint64(0xFFFFFF)

	sf := New()
	assertNotNil(t, sf)
	assertEqual(t, expectedTimeMask, sf.timeMask)
	assertEqual(t, expectedNodeMask, sf.nodeMask)
	assertEqual(t, expectedStepMask, sf.stepMask)
}

func TestNewWithConfig(t *testing.T) {
	name := func(nodeBits, stepBits uint8) string {
		return fmt.Sprintf("nodeBits=%d.stepBits=%d", nodeBits, stepBits)
	}
	tests := []struct {
		epoch    int64
		nodeBits uint8
		stepBits uint8

		expectedEpoch    time.Time
		expectedTimeMask uint64
		expectedNodeMask uint64
		expectedStepMask uint64
		expectedErr      error
	}{
		{
			epoch:            int64(1577836800),
			nodeBits:         16,
			stepBits:         16,
			expectedEpoch:    time.Unix(1577836800, 0),
			expectedTimeMask: 0xFFFFFFFF00000000,
			expectedNodeMask: 0xFFFF0000,
			expectedStepMask: 0xFFFF,
			expectedErr:      nil,
		},
		{
			epoch:            int64(807148800),
			nodeBits:         8,
			stepBits:         4,
			expectedEpoch:    time.Unix(807148800, 0),
			expectedTimeMask: 0xFFFFFFFF000,
			expectedNodeMask: 0xFF0,
			expectedStepMask: 0xF,
			expectedErr:      nil,
		},
		{
			epoch:            int64(1577836800),
			nodeBits:         2,
			stepBits:         0,
			expectedEpoch:    time.Unix(1577836800, 0),
			expectedTimeMask: 0b1111111111111111111111111111111100,
			expectedNodeMask: 0b11,
			expectedStepMask: 0b0,
			expectedErr:      nil,
		},
		{
			epoch:       int64(1577836800),
			nodeBits:    17,
			stepBits:    16,
			expectedErr: errors.New("nodeBits + stepBits should has 32 in total, got 33"),
		},
	}

	for _, tt := range tests {
		t.Run(name(tt.nodeBits, tt.stepBits), func(t *testing.T) {
			sf, err := NewWithConfig(tt.epoch, tt.nodeBits, tt.stepBits)
			assertEqual(t, tt.expectedErr, err)
			if err == nil {
				assertNotNil(t, sf)
				assertEqual(t, tt.expectedEpoch, sf.epoch)
				assertEqual(t, tt.expectedTimeMask, sf.timeMask)
				assertEqual(t, tt.expectedNodeMask, sf.nodeMask)
				assertEqual(t, tt.expectedStepMask, sf.stepMask)
			}
		})
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
