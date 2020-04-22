package snowfake

import (
	"fmt"
	"time"
)

type snowfake struct {
	config
	stepBit uint32
}

func New() *snowfake {
	nodeBits := uint8(8)
	epoch := int64(1577836800) // epoch start from 01/01/2020 @ 12:00am (UTC)
	sf, _ := NewWithConfig(epoch, nodeBits, maxBits-timeBits-nodeBits)
	return sf
}

func NewWithConfig(epoch int64, nodeBits, stepBits uint8) (*snowfake, error) {
	if timeBits+nodeBits+stepBits > maxBits {
		expectedBits := maxBits - timeBits
		actualBits := nodeBits + stepBits
		return nil, fmt.Errorf("nodeBits + stepBits should has %d in total, got %d", expectedBits, actualBits)
	}

	s := &snowfake{}
	s.nodeBits = nodeBits
	s.stepBits = stepBits

	s.stepBit = 0

	s.epoch = time.Unix(epoch, 0)
	s.stepMask = 1<<stepBits - 1
	s.nodeMask = (1<<nodeBits - 1) << stepBits
	s.timeMask = (1<<timeBits - 1) << (nodeBits + stepBits)

	return s, nil
}
