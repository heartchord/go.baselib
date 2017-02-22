package mempool

import (
	"runtime"
	"sort"
	"sync"
	"sync/atomic"
	"unsafe"
)

// defaultBytesBlockSizeSet defines all sizes of bytes-block with fixed-length.
var defaultBytesBlockSizeSet = []int{
	/*
		1 * 8,    // 8     bytes
		2 * 8,    // 16    bytes
		4 * 8,    // 32    bytes
		8 * 8,    // 64    bytes
		16 * 8,   // 128   bytes
		32 * 8,   // 256   bytes
		64 * 8,   // 512   bytes
	*/
	1 * 1024, // 1024  bytes(1  KB)
	/*
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
	*/
}

func getIndex(blockSize int) int {
	return sort.SearchInts(defaultBytesBlockSizeSet, blockSize)
}

var defaultBytesBlockReservedHeaderSize = 8

// BytesBlock :
type BytesBlock struct {
	Length int         // 用户申请内存大小 <= 用户内存块实际长度len(BytesBlock.Buffer)
	Header []byte      // 缓冲保留头部空间
	Buffer []byte      // 用户数据存放缓冲
	buffer []byte      // 真正数据存放空间 = bb.Header + bb.Buffer
	next   *BytesBlock // 下一个字节内存块
	prev   *BytesBlock // 上一个字节内存块
}

func newBytesBlock(blockSize int) (bb *BytesBlock) {
	bb = new(BytesBlock)
	blockSize += defaultBytesBlockReservedHeaderSize
	bb.buffer = make([]byte, blockSize)
	bb.Header = bb.buffer[0:defaultBytesBlockReservedHeaderSize]
	bb.Buffer = bb.buffer[defaultBytesBlockReservedHeaderSize:blockSize]
	bb.Length = 0
	return
}

func (bb *BytesBlock) reset() {
	bb.next = nil
	bb.prev = nil
	bb.Length = 0
}

type bytesBlockList struct {
	root      *BytesBlock // 字节内存块链表根结点
	size      int         // 字节内存块链表结点数
	blockSize int         // 管理的每个内存块大小
}

func newBytesBlockList(blockSize int) (l *bytesBlockList) {
	l = new(bytesBlockList)
	l.root = newBytesBlock(blockSize)
	l.root.next = l.root
	l.root.prev = l.root
	l.blockSize = blockSize

	for i := 0; i < 1; i++ {
		l.pushBack(newBytesBlock(blockSize))
	}
	return
}

func (l *bytesBlockList) pushBack(e *BytesBlock) (result bool) {
	r := l.root
	if r.next == nil {
		r.next = r
		r.prev = r
		l.size = 0
	}

	if e != nil && len(e.buffer) == l.blockSize+defaultBytesBlockReservedHeaderSize {
		p := r.prev
		r.prev = e
		e.next = r
		e.prev = p
		p.next = e
		l.size++
		result = true
	}

	return
}

func (l *bytesBlockList) pushFront(e *BytesBlock) (result bool) {
	r := l.root
	if r.next == nil {
		r.next = r
		r.prev = r
		l.size = 0
	}

	if e != nil && len(e.buffer) == l.blockSize+defaultBytesBlockReservedHeaderSize {
		n := r.next
		r.next = e
		e.prev = r
		e.next = n
		n.prev = e
		l.size++
		result = true
	}

	return
}

func (l *bytesBlockList) popBack() (bb *BytesBlock) {
	r := l.root
	if r.next == nil {
		panic("unexpected value")
	}

	if l.size == 0 {
		return nil
	}

	bb = r.prev
	p := r.prev.prev
	p.next = r
	r.prev = p
	bb.next = nil
	bb.prev = nil
	l.size--

	return
}

func (l *bytesBlockList) popFront() (bb *BytesBlock) {
	r := l.root
	if r.next == nil {
		panic("unexpected value")
	}

	if l.size == 0 {
		return nil
	}

	bb = r.next
	p := r.next.next
	p.prev = r
	r.next = p
	bb.next = nil
	bb.prev = nil
	l.size--
	return
}

type localBytesPool struct {
	private map[int]*bytesBlockList
	shared  map[int]*bytesBlockList
	sync.Mutex
}

func (lbp *localBytesPool) getPrivate(blockSize int) *BytesBlock {
	//idx := getIndex(blockSize)
	idx := blockSize
	return lbp.private[idx].popFront()
}

func (lbp *localBytesPool) putPrivate(bb *BytesBlock) {
	//idx := getIndex(len(bb.Buffer))
	idx := len(bb.Buffer)
	lbp.private[idx].pushBack(bb)
}

func (lbp *localBytesPool) getShared(blockSize int) *BytesBlock {
	//idx := getIndex(blockSize)
	idx := blockSize
	lbp.Lock()
	defer lbp.Unlock()
	return lbp.shared[idx].popFront()
}

func (lbp *localBytesPool) putShared(bb *BytesBlock) {
	//idx := getIndex(len(bb.Buffer))
	idx := len(bb.Buffer)
	lbp.Lock()
	defer lbp.Unlock()
	lbp.shared[idx].pushBack(bb)
}

func newLocalBytesPool() (lbp *localBytesPool) {
	lbp = new(localBytesPool)

	lbp.private = make(map[int]*bytesBlockList)
	lbp.shared = make(map[int]*bytesBlockList)

	for i := 0; i < len(defaultBytesBlockSizeSet); i++ {
		blockSize := defaultBytesBlockSizeSet[i]
		//idx := getIndex(blockSize)
		idx := blockSize
		lbp.private[idx] = newBytesBlockList(blockSize)
		lbp.shared[idx] = newBytesBlockList(blockSize)
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
	if bb == nil {
		return
	}

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
		x = newBytesBlock(blockSize)
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
//go:nosplit
func sync_runtime_procPin() int

//go:linkname sync_runtime_procUnpin sync.runtime_procUnpin
//go:nosplit
func sync_runtime_procUnpin()
