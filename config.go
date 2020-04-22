package snowfake

import "time"

const (
	timeBits = uint8(32)
	maxBits  = uint8(64)
)

// In total, NodeBits + StepBits less than or equal to 32
// timeBits has fixed 32 bits allocation
type config struct {
	timeBits uint8
	nodeBits uint8
	stepBits uint8

	epoch    time.Time
	timeMask uint64
	nodeMask uint64
	stepMask uint64
}
