package utils

import (
	"fmt"
	"runtime"
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

func NewSnowflake(epoch time.Time, nodeIDBits, sequenceBits uint) (*Snowflake, error) {

	if nodeIDBits+sequenceBits > 63 {
		return nil, fmt.Errorf("nodeIDBits (%d) + sequenceBits (%d) must not exceed 63", nodeIDBits, sequenceBits)
	}

	sf := &Snowflake{
		epoch:        epoch.UnixMilli(),
		nodeIDBits:   nodeIDBits,
		sequenceBits: sequenceBits,
	}
	sf.maxNodeID = -1 ^ (-1 << nodeIDBits)
	sf.maxSequence = -1 ^ (-1 << sequenceBits)
	return sf, nil
}

func (sf *Snowflake) Generate(nodeID int64) (int64, error) {

	if nodeID < 0 || nodeID > sf.maxNodeID {
		return 0, fmt.Errorf("nodeID %d out of range [0, %d]", nodeID, sf.maxNodeID)
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
				runtime.Gosched()
				continue
			}
		} else {
			if atomic.CompareAndSwapInt64(&sf.sequence, sf.sequence, 0) &&
				atomic.CompareAndSwapInt64(&sf.lastTime, lastTime, currentTime) {
				newSequence = 0
			} else {
				time.Sleep(time.Millisecond)
				continue
			}
		}

		id := ((currentTime - sf.epoch) << (sf.nodeIDBits + sf.sequenceBits)) |
			(nodeID << sf.sequenceBits) |
			newSequence

		return id, nil
	}
}
