package snowfake

import (
	"fmt"
	"time"
)

type snowfake struct {
	config

	node  uint64
	epoch uint64
	time  uint64
	step  uint64
}

func New(node uint64) *snowfake {
	nodeBits := uint8(8)
	epoch := uint64(1577836800) // epoch start from 01/01/2020 @ 12:00am (UTC)
	sf, _ := NewWithConfig(node, epoch, nodeBits, maxBits-timeBits-nodeBits)
	return sf
}

func NewWithConfig(node, epoch uint64, nodeBits, stepBits uint8) (*snowfake, error) {
	if timeBits+nodeBits+stepBits > maxBits {
		expectedBits := maxBits - timeBits
		actualBits := nodeBits + stepBits
		return nil, fmt.Errorf("nodeBits + stepBits should has %d in total, got %d", expectedBits, actualBits)
	}

	if node >= 1<<nodeBits {
		return nil, fmt.Errorf("node should below %d", 1<<nodeBits)
	}

	s := &snowfake{}

	s.node = node
	s.epoch = epoch
	s.time = 0
	s.step = 0

	s.nodeBits = nodeBits
	s.stepBits = stepBits

	s.timeShift = nodeBits + stepBits
	s.nodeShift = stepBits

	s.stepMask = 1<<stepBits - 1
	s.nodeMask = (1<<nodeBits - 1) << s.nodeShift
	s.timeMask = (1<<timeBits - 1) << s.timeShift

	return s, nil
}
