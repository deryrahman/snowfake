package snowfake

import (
	"fmt"
	"sync"
	"time"
)

// Snowfake is an object to generate the ID
type Snowfake struct {
	config

	mu *sync.Mutex

	node  uint64
	epoch uint64
	time  uint64
	step  uint64
}

// New creates new snowfake object with default config.
// Args: node should be within 8bit
// By default, epoch starts from 01/01/2020 @ 12:00am (UTC),
// 256 slots of node (0-255), and 24bit of stepBits which can provide 16777216 req/s when generating
// ID. See GenerateID()
func New(node uint64) *Snowfake {
	nodeBits := uint8(8)
	epoch := uint64(1577836800) // epoch start from 01/01/2020 @ 12:00am (UTC)
	sf, _ := NewWithConfig(node, epoch, nodeBits, maxBits-timeBits-nodeBits)
	return sf
}

// NewWithConfig creates new snowfake object.
// Args: node should be within nodeBits range, epoch should be within 32bit
// nodeBits+stepBits should be less than or equal to 32
// Note: use large stepBits if you want to provide high rate per second when generating
// ID. See GenerateID()
func NewWithConfig(node, epoch uint64, nodeBits, stepBits uint8) (*Snowfake, error) {
	if timeBits+nodeBits+stepBits > maxBits {
		expectedBits := maxBits - timeBits
		actualBits := nodeBits + stepBits
		return nil, fmt.Errorf("nodeBits + stepBits should has %d in total, got %d", expectedBits, actualBits)
	}

	if node >= 1<<nodeBits {
		return nil, fmt.Errorf("node should below %d", 1<<nodeBits)
	}

	s := &Snowfake{}

	s.mu = &sync.Mutex{}

	s.node = node
	s.epoch = epoch
	s.time = 0
	s.step = 0

	s.timeBits = timeBits
	s.nodeBits = nodeBits
	s.stepBits = stepBits

	s.timeShift = nodeBits + stepBits
	s.nodeShift = stepBits
	s.stepShift = 0

	s.timeMask = (1<<timeBits - 1) << s.timeShift
	s.nodeMask = (1<<nodeBits - 1) << s.nodeShift
	s.stepMask = (1<<stepBits - 1) << s.stepShift

	return s, nil
}

// GenerateID generates new ID within 64bit
// it's not guarantee collision if you use small stepBits.
// Rule of thumb, 1024 req/s can be safely generated without collision
// if you use 10 stepBits
func (s *Snowfake) GenerateID() uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	t := s.now()
	if s.time == t {
		// not guarantee if you use small stepBits
		// since it will probably collide within same
		// timestamp
		s.step++
		s.step &= 1<<s.stepBits - 1
	} else {
		s.step = 0
	}

	s.time = t

	r := (s.time << s.timeShift) & s.timeMask
	r |= (s.node << s.nodeShift) & s.nodeMask
	r |= (s.step << s.stepShift) & s.stepMask

	return r
}

func (s *Snowfake) now() uint64 {

	t := uint64(time.Now().Unix())
	t -= s.epoch

	if ((1<<timeBits - 1) & t) == t {
		return t
	}

	return 0
}
