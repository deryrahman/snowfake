package snowfake_test

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/deryrahman/snowfake"
)

func TestNew(t *testing.T) {
	resetConfig()
	name := func(nodeID uint64) string {
		return fmt.Sprintf("nodeID=%d", nodeID)
	}
	tests := []struct {
		nodeID   uint64
		nodeBits uint8

		expectedErr error
	}{
		{
			nodeID:      uint64(1),
			nodeBits:    16,
			expectedErr: nil,
		},
		{
			nodeID:      uint64(4),
			nodeBits:    2,
			expectedErr: errors.New("nodeID should less than 4"),
		},
	}

	for _, tt := range tests {
		t.Run(name(tt.nodeID), func(t *testing.T) {
			snowfake.NodeBits = tt.nodeBits
			_ = snowfake.Init()

			sf, err := snowfake.New(tt.nodeID)

			assertEqual(t, tt.expectedErr, err)
			if err == nil {
				assertNotNil(t, sf)
			}
		})
	}
}

func TestSnowfake_GenerateID(t *testing.T) {
	resetConfig()
	nodeID := uint64(29)

	sf, _ := snowfake.New(nodeID)

	estimateTimeFromID := uint64(time.Now().Unix()) - snowfake.Epoch
	expectedNodeFromID := uint64(29)
	expectedSeqFromID := uint64(0)

	assertNotNil(t, sf)
	if sf != nil {
		id := sf.GenerateID()
		assertTrue(t, estimateTimeFromID <= ((id&snowfake.GetTimeMask())>>snowfake.GetTimeShift()))
		assertEqual(t, expectedNodeFromID, (id&snowfake.GetNodeMask())>>snowfake.GetNodeShift())
		assertEqual(t, expectedSeqFromID, id&snowfake.GetSeqMask())
	}

}

func TestSnowfake_GenerateID_Collision(t *testing.T) {
	resetConfig()
	nodeID := uint64(1)
	concurrent := 10000

	sf, _ := snowfake.New(nodeID)

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

func resetConfig() {
	snowfake.Epoch = 1577836800
	snowfake.NodeBits = 5
	snowfake.SeqBits = 27
	_ = snowfake.Init()
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
