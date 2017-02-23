package mempool

import (
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"

	. "github.com/heartchord/goblazer"
)

// defaultBytesBlockSizeSet defines all sizes of bytes-block with fixed-length.
var defaultBytesBlockSizeSet = []int{
	1 * 8,     // 8     bytes
	2 * 8,     // 16    bytes
	4 * 8,     // 32    bytes
	8 * 8,     // 64    bytes
	16 * 8,    // 128   bytes
	32 * 8,    // 256   bytes
	64 * 8,    // 512   bytes
	1 * 1024,  // 1024  bytes(1  KB)
	2 * 1024,  // 2048  bytes(2  KB)
	3 * 1024,  // 3072  bytes(3  KB)
	4 * 1024,  // 4096  bytes(4  KB)
	5 * 1024,  // 5120  bytes(5  KB)
	6 * 1024,  // 6144  bytes(6  KB)
	7 * 1024,  // 7168  bytes(7  KB)
	8 * 1024,  // 8192  bytes(8  KB)
	16 * 1024, // 16384 bytes(16 KB)
	32 * 1024, // 32768 bytes(32 KB)
	64 * 1024, // 65536 bytes(64 KB)
}

var defaultBytesBlockSizeNum = len(defaultBytesBlockSizeSet)
var defaultBytesBlockReservedHeaderSize = 8

type localBytesPool struct {
	private []*bytesBlockList
	shared  []*bytesBlockList
	sync.Mutex
}

func (lbp *localBytesPool) getPrivate(blockSize int) (bb *BytesBlock) {
	idx := BinarySearchIntsGE(defaultBytesBlockSizeSet, blockSize)
	if idx >= 0 {
		bb = lbp.private[idx].getFront()
	}
	return
}

func (lbp *localBytesPool) putPrivate(bb *BytesBlock) {
	if bb.idx >= 0 {
		lbp.private[bb.idx].putBack(bb)
	}
}

func (lbp *localBytesPool) getShared(blockSize int) (bb *BytesBlock) {
	idx := BinarySearchIntsGE(defaultBytesBlockSizeSet, blockSize)
	if idx >= 0 {
		lbp.Lock()
		bb = lbp.shared[idx].getFront()
		defer lbp.Unlock()
	}
	return
}

func (lbp *localBytesPool) putShared(bb *BytesBlock) {
	if bb.idx >= 0 {
		lbp.Lock()
		defer lbp.Unlock()
		lbp.shared[bb.idx].putBack(bb)
	}
}

func newLocalBytesPool() (lbp *localBytesPool) {
	lbp = new(localBytesPool)

	lbp.private = make([]*bytesBlockList, defaultBytesBlockSizeNum)
	lbp.shared = make([]*bytesBlockList, defaultBytesBlockSizeNum)

	for i := 0; i < defaultBytesBlockSizeNum; i++ {
		blockSize := defaultBytesBlockSizeSet[i]
		lbp.private[i] = newBytesBlockList(blockSize)
		lbp.shared[i] = newBytesBlockList(blockSize)
	}

	return
}

type BytesPool struct {
	local     unsafe.Pointer
	localSize uintptr
}

func NewBytesPool() *BytesPool {
	bp := new(BytesPool)

	size := runtime.GOMAXPROCS(0)
	local := make([]*localBytesPool, size)

	atomic.StorePointer(&bp.local, unsafe.Pointer(&local[0]))
	atomic.StoreUintptr(&bp.localSize, uintptr(size))

	for i := 0; i < size; i++ {
		local[i] = newLocalBytesPool()
	}

	return bp
}

func (p *BytesPool) Put(bb *BytesBlock) {
	l := p.pin()
	l.putPrivate(bb)
	sync_runtime_procUnpin()
}

func (p *BytesPool) pin() *localBytesPool {
	pid := sync_runtime_procPin()
	s := atomic.LoadUintptr(&p.localSize)
	l := p.local
	if uintptr(pid) < s {
		return indexLocalPool(l, pid)
	}
	return nil
}

func (p *BytesPool) Get(blockSize int) *BytesBlock {
	l := p.pin()
	x := l.getPrivate(blockSize)
	sync_runtime_procUnpin()
	if x == nil {
		x = l.getShared(blockSize)
		if x == nil {
			x = p.getSlow(blockSize)
		}
	}

	if x == nil {
		idx := BinarySearchIntsGE(defaultBytesBlockSizeSet, blockSize)
		x = newBytesBlock(blockSize, idx)
	}
	return x
}

func (p *BytesPool) getSlow(blockSize int) (x *BytesBlock) {
	size := atomic.LoadUintptr(&p.localSize)
	local := p.local

	pid := sync_runtime_procPin()
	sync_runtime_procUnpin()
	for i := 0; i < int(size); i++ {
		l := indexLocalPool(local, (pid+i+1)%int(size))
		x = l.getShared(blockSize)
		if x != nil {
			break
		}
	}
	return x
}

func indexLocalPool(l unsafe.Pointer, i int) *localBytesPool {
	return (*[1000000]*localBytesPool)(l)[i]
}

//go:linkname sync_runtime_procPin sync.runtime_procPin
func sync_runtime_procPin() int

//go:linkname sync_runtime_procUnpin sync.runtime_procUnpin
func sync_runtime_procUnpin()
