package snowfake

import "fmt"

const (
	// TimeBits is allocation number of timestamp
	TimeBits = uint8(32)
	// MaxBits is allocation number of snowfake ID
	MaxBits = uint8(64)
)

var (
	// Epoch is an origin timestamp in second
	Epoch uint64 = 1577836800 // epoch start from 01/01/2020 @ 12:00am (UTC)
	// NodeBits is allocation number of machines
	NodeBits uint8 = 5
	// SeqBits is allocation number of sequence
	SeqBits uint8 = 27

	timeShift uint8
	nodeShift uint8
	seqShift  uint8

	timeMask uint64
	nodeMask uint64
	seqMask  uint64

	maxNode uint64
)

func init() {
	_ = Init()
}

// Init loads configuration based on NodeBits and SeqBits
func Init() error {
	if TimeBits+NodeBits+SeqBits > MaxBits {
		return fmt.Errorf("NodeBits + SeqBits should has %d in total", MaxBits-TimeBits)
	}

	timeShift = NodeBits + SeqBits
	nodeShift = SeqBits
	seqShift = 0

	timeMask = (1<<TimeBits - 1) << timeShift
	nodeMask = (1<<NodeBits - 1) << nodeShift
	seqMask = (1<<SeqBits - 1) << seqShift

	maxNode = 1 << NodeBits

	return nil
}

// GetTimeShift returns timeShift from config
func GetTimeShift() uint8 {
	return timeShift
}

// GetNodeShift returns nodeShift from config
func GetNodeShift() uint8 {
	return nodeShift
}

// GetSeqShift returns seqShift from config
func GetSeqShift() uint8 {
	return seqShift
}

// GetTimeMask returns timeMask from config
func GetTimeMask() uint64 {
	return timeMask
}

// GetNodeMask returns nodeMask from config
func GetNodeMask() uint64 {
	return nodeMask
}

// GetSeqMask returns seqMask from config
func GetSeqMask() uint64 {
	return seqMask
}

// GetMaxNode returns maxNode from config
func GetMaxNode() uint64 {
	return maxNode
}
