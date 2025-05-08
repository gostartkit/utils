package utils

import (
	"fmt"
	"sync/atomic"
	"time"
)

type Snowflake struct {
	epoch        int64
	sequence     int64
	lastTime     int64
	nodeIDBits   uint
	sequenceBits uint
	maxNodeID    int64
	maxSequence  int64
}

func NewSnowflake(epoch time.Time) *Snowflake {
	const (
		defaultNodeIDBits   = 10
		defaultSequenceBits = 12
	)

	sf := &Snowflake{
		epoch:        epoch.UnixMilli(),
		nodeIDBits:   defaultNodeIDBits,
		sequenceBits: defaultSequenceBits,
	}

	sf.maxNodeID = -1 ^ (-1 << sf.nodeIDBits)
	sf.maxSequence = -1 ^ (-1 << sf.sequenceBits)

	return sf
}

func (sf *Snowflake) Generate(nodeID int64) (int64, error) {

	if nodeID < 0 || nodeID > sf.maxNodeID {
		return 0, fmt.Errorf("nodeID must be between 0 and %d", sf.maxNodeID)
	}

	for {
		currentTime := time.Now().UnixMilli()
		lastTime := atomic.LoadInt64(&sf.lastTime)

		if currentTime < lastTime {
			return 0, fmt.Errorf("clock moved backwards")
		}

		var newSequence int64
		if currentTime == lastTime {
			newSequence = atomic.AddInt64(&sf.sequence, 1) & sf.maxSequence
			if newSequence == 0 {
				continue
			}
		} else {
			if atomic.CompareAndSwapInt64(&sf.sequence, sf.sequence, 0) &&
				atomic.CompareAndSwapInt64(&sf.lastTime, lastTime, currentTime) {
				newSequence = 0
			} else {
				continue
			}
		}

		id := ((currentTime - sf.epoch) << (sf.nodeIDBits + sf.sequenceBits)) |
			(nodeID << sf.sequenceBits) |
			newSequence

		return id, nil
	}
}
