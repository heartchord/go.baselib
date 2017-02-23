package mempool

// BytesBlock :
type BytesBlock struct {
	Header    []byte      // 缓冲保留头部空间
	Buffer    []byte      // 用户数据存放缓冲
	ApplySize int         // 用户申请内存大小
	idx       int         // 用户实际内存大小
	buffer    []byte      // 真正数据存放空间 = bb.Header + bb.Buffer
	next      *BytesBlock // 下一个字节内存块
	prev      *BytesBlock // 上一个字节内存块
}

func newBytesBlock(blockSize int, idx int) (bb *BytesBlock) {
	bb = new(BytesBlock)
	totalSize := blockSize + defaultBytesBlockReservedHeaderSize
	bb.buffer = make([]byte, totalSize)
	bb.Header = bb.buffer[0:defaultBytesBlockReservedHeaderSize]
	bb.Buffer = bb.buffer[defaultBytesBlockReservedHeaderSize:totalSize]
	bb.ApplySize = 0
	bb.idx = idx
	return
}

func (bb *BytesBlock) reset() {
	bb.next = nil
	bb.prev = nil
	bb.ApplySize = 0
}
