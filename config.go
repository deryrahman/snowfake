package snowfake

import "fmt"

const (
	timeBits = uint8(32)
	maxBits  = uint8(64)
)

var (
	epoch    uint64
	nodeBits uint8
	seqBits  uint8

	timeShift uint8
	nodeShift uint8
	seqShift  uint8

	timeMask uint64
	nodeMask uint64
	seqMask  uint64

	maxNode uint64
)

func init() {
	// default configuration
	SetEpoch(1577836800) // epoch start from 01/01/2020 @ 12:00am (UTC)
	SetNodeBits(5)
	SetSeqBits(27)

	_ = Init()
}

// Init loads configuration based on nodeBits and seqBits
func Init() error {
	if timeBits+nodeBits+seqBits > maxBits {
		return fmt.Errorf("nodeBits + seqBits should has %d in total", maxBits-timeBits)
	}

	timeShift = nodeBits + seqBits
	nodeShift = seqBits
	seqShift = 0

	timeMask = (1<<timeBits - 1) << timeShift
	nodeMask = (1<<nodeBits - 1) << nodeShift
	seqMask = (1<<seqBits - 1) << seqShift

	maxNode = 1 << nodeBits

	return nil
}

// SetEpoch sets epoch configuration
func SetEpoch(e uint64) {
	epoch = e
}

// SetNodeBits sets nodeBits configuration
func SetNodeBits(n uint8) {
	nodeBits = n
}

// SetSeqBits sets seqBits configuration
func SetSeqBits(s uint8) {
	seqBits = s
}
