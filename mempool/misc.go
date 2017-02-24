package mempool

import "fmt"

// NewFunc is a function type to create a new object.
type NewFunc func() interface{}

// PoolStats contais some statistics about pool put & get operation.
type PoolStats struct {
	Enabled            bool  // true - open statistics, close - close statistics.
	PutOkTimes         int64 // times that puts memory block back to memory pool successfully.
	PutDiscardTimes    int64 // times that drops memory block to floor(not back to memory pool).
	GetOkTimes         int64 // times that gets memory block from memory pool successfully.
	GetByNewTimes      int64 // times that gets memory block from memory pool by new operation.
	GetFromOtherPTimes int64 // times that gets memory block from memory pool in other P.
}

func (ps PoolStats) PrintAllInfo() {
	fmt.Printf("Stats = %#v\n", ps)
}
