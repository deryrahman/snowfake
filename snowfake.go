package snowfake

import (
	"fmt"
)

type snowfake struct {
	config
	stepBit uint32
}

func New() *snowfake {
	nodeBits := uint8(8)
	sf, _ := NewWithConfig(nodeBits, maxBits-timeBits-nodeBits)
	return sf
}

func NewWithConfig(nodeBits, stepBits uint8) (*snowfake, error) {
	if timeBits+nodeBits+stepBits > maxBits {
		expectedBits := maxBits - timeBits
		actualBits := nodeBits + stepBits
		return nil, fmt.Errorf("nodeBits + stepBits should has %d in total, got %d", expectedBits, actualBits)
	}

	s := &snowfake{}
	s.nodeBits = nodeBits
	s.stepBits = stepBits

	s.stepBit = 0
	s.stepMask = 1<<stepBits - 1
	s.nodeMask = (1<<nodeBits - 1) << stepBits
	s.timeMask = (1<<timeBits - 1) << (nodeBits + stepBits)

	return s, nil
}
