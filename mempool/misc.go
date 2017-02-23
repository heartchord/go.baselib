package mempool

// NewFunc is a function type to create a new object.
type NewFunc func() interface{}

type PoolStats struct {
	PutOkTimes         int64
	PutDiscardTimes    int64
	GetOkTimes         int64
	GetByNewTimes      int64
	GetFromOtherPTimes int64
}
