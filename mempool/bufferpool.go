package mempool

import (
	"fmt"
	"sync/atomic"
)

// Buffer achieves a data struct based on bytes block. Buffer uses a bytes block allocated from BufferPool. User can control its life
// circle manually, and if user want to return it to BufferPool, just invokes Buffer.DecRef() to reduce its reference counter to zero.
// Notice that if you want to share a Buffer with others, you must invoke Buffer.AddRef() to add its reference, otherwise it may be
// recycled unexpectedly.
//     1. Generally, user will need a memory buffer with a required size. Eg. maybe user wants a buffer that can hold a data
//        with a fixed length 'm', so we will allocate a memory buffer by 'm', and here 'oriSize' = 'curSize' = m.
//     2. After memory buffer allocated, 'oriSize' will hold the initialise size in case that user changes the 'curSize' and
//        wants to restore the 'curSize' to the initialise state.
//     3. Sometimes, user may want to shorten or expand the data length, user can get to this purpose by shortening or expanding
//        'curSize', and in this case, 'curSize' <= cap(buffer) - resSize.
//     4. 'resSize' is assigned when one buffer is allocated, it can't be changed by any means, and it's invisible in most cases.
//        If user wants to add some data header to the data, user can just access this space and fill the extra data in it.
type Buffer struct {
	owner   *BufferPool // belonger that manages all Buffer with the same buffer capacity. Eg, if cap(Buffer.buf)=1024, the class of BufferPool it belonged to is 1024 too.
	buf     []byte      // real buffer.
	ref     int32       // reference counter
	resSize uint32      // reserved size
	oriSize uint32      // original size
	curSize uint32      // current  size
}

// AddRef increases the reference counter. If you want to share one Buffer with others, don't forget invoke this function.
func (b *Buffer) AddRef() int32 {
	ref := atomic.AddInt32(&b.ref, 1)

	if ref <= 0 { // unexpected error
		panic(fmt.Sprintf("[Buffer.AddRef error] expected Buffer.ref > 0, got Buffer.ref = %d", b.ref))
	}

	return ref
}

// DecRef decreases the reference counter. If you want to return current Buffer to BufferPool, just invoke this function.
func (b *Buffer) DecRef() int32 {
	ref := atomic.AddInt32(&b.ref, -1)

	if ref < 0 { // unexpected error
		panic(fmt.Sprintf("[Buffer.DecRef error] expected Buffer.ref >= 0, got Buffer.ref = %d", b.ref))
	}

	if ref == 0 { // only executed once
		b.owner.put(b)
	}

	return ref
}

// ResSize returns the reserved size of the buffer.
func (b *Buffer) ResSize() int {
	return int(b.resSize)
}

// CurSize returns the current size of buffer.
func (b *Buffer) CurSize() int {
	return int(b.curSize)
}

// OriSize returns the original size of buffer.
func (b *Buffer) OriSize() int {
	return int(b.oriSize)
}

// SetSize changes the current size of buffer.
func (b *Buffer) SetSize(newSize uint32) {
	if newSize == 0 {
		panic(fmt.Sprintf("[Buffer.SetSize error] expected newSize > 0, got newSize = %d", newSize))
	}

	if b.resSize+newSize > uint32(cap(b.buf)) {
		panic(fmt.Sprintf("[Buffer.SetSize error] expected newSize = %d, it's too large.", newSize))
	}

	atomic.SwapUint32(&b.curSize, newSize)
}

// ResetSize resets the current size to original state.
func (b *Buffer) ResetSize() {
	atomic.SwapUint32(&b.curSize, b.oriSize)
}

// ResBuf returns the reserved space of the buffer.
func (b *Buffer) ResBuf() []byte {
	if b.resSize > uint32(cap(b.buf)) {
		panic("[Buffer.ResBuf error] Buffer.resSize is too large!")
	}

	return b.buf[0:b.resSize]
}

// Buf returns the space of the buffer with current size.
func (b *Buffer) Buf() []byte {
	if b.resSize+b.curSize > uint32(cap(b.buf)) {
		panic("[Buffer.Buf error] Buffer.resSize or Buffer.curSize is too large!")
	}

	return b.buf[b.resSize : b.resSize+b.curSize]
}

// BufferPool holds all Buffers with specified block size. The Buffers with same capacity will be placed in the same Pool.
//    Eg, All 1024 bytes blocks will be placed in the pool with the class type 1024. The pools with different class make up
//    the final BytesPool.
type BufferPool struct {
	classes  []int
	allPools []*Pool
}

// Get fetches one Buffer from BufferPool.
// The parameter 'reservedSize' represents the reserved space capacity.
// The parameter 'requiredSize' represents the buffer size user demands, so 'reservedSize' + 'requiredSize' <= cap(Buffer.buf).
// When user applies one buffer, BufferPool will make slicing operation to make sure the space user got is just sufficient.
func (bp *BufferPool) Get(reservedSize int, requiredSize int) (b *Buffer) {
	if reservedSize < 0 {
		panic(fmt.Sprintf("[BufferPool.Get error] expected reservedSize >= 0, got reservedSize = %d", reservedSize))
	}

	if requiredSize < 0 {
		panic(fmt.Sprintf("[BufferPool.Get error] expected requiredSize >= 0, got requiredSize = %d", requiredSize))
	}

	idx := bp.getClassIdx(reservedSize + requiredSize)
	if idx >= 0 {
		b = bp.allPools[idx].Get().(*Buffer)
		b.buf = b.buf[:reservedSize+requiredSize] // reduce the buffer length to user required.
		b.ref = 1                                 // the reference counter is 1 when allocated.
		b.resSize = uint32(reservedSize)
		b.curSize = uint32(requiredSize)
		b.oriSize = uint32(requiredSize)
		return
	}

	// New operation. The class can't be found, so this size is biger than all blockSize in bp.classes.
	// When recycled, because the cap is bigger than all pools, the BytesBlock will just be dropped to the floor.
	b = &Buffer{
		owner:   bp,
		ref:     1,
		buf:     make([]byte, reservedSize+requiredSize),
		resSize: uint32(reservedSize),
		oriSize: uint32(requiredSize),
		curSize: uint32(requiredSize),
	}
	return
}

// Reset will set the status of all pools to initialization.
func (bp *BufferPool) Reset() {
	for _, pool := range bp.allPools {
		pool.Reset()
	}
}

// put returns the BytesBlock to the appropriate pool, if any appropriate pool can't be found, it will just be dropped to the floor.
func (bp *BufferPool) put(b *Buffer) {
	idx := bp.getClassIdx(cap(b.buf))
	if idx >= 0 {
		bp.allPools[idx].Put(b)
	}
}

// getClassIdx will find the appropriate index of BufferPool.classes by 'requiredSize'.
func (bp *BufferPool) getClassIdx(requiredSize int) (idx int) {
	idx = -1
	if requiredSize <= bp.classes[len(bp.classes)-1] {
		for idx = 0; idx < len(bp.classes); idx++ {
			len := bp.classes[idx]
			if requiredSize <= len {
				break
			}
		}
	}
	return
}

// NewBufferPool create a new BufferPool container.
//     classes      : specify all kinds of pool. You can use default 'poolClasses' or specify a new classes.
//     localPoolCap : specify the local pool capacity of every 'P'. It specifies the max value of the sum of
//                    pool.localPool.private + pool.localPool.shared.
func NewBufferPool(classes []int, localPoolCap int) *BufferPool {
	if !checkAllClasses(classes) {
		panic(fmt.Sprintf("[NewBufferPool error] expected a ascend classes, got classes is %#v", classes))
	}

	// create an instance.
	classCount := len(classes)
	bp := &BufferPool{
		classes:  make([]int, classCount),
		allPools: make([]*Pool, classCount),
	}

	// Initialization.
	for i := 0; i < classCount; i++ {
		idx := i // used in closure.
		bp.classes[i] = classes[i]
		bp.allPools[i] = NewPool(
			localPoolCap,
			func() interface{} {
				l := classes[idx]
				return &Buffer{
					owner:   bp,
					ref:     0,
					buf:     make([]byte, l),
					resSize: 0,
					oriSize: 0,
					curSize: 0,
				}
			},
		)
	}

	return bp
}

var defaultBufferPool *BufferPool

func init() {
	defaultBufferPool = NewBufferPool(poolClasses, 8)
}

// GetFromDefaultBufferPool fetches one Buffer from defaultBufferPool.
func GetFromDefaultBufferPool(reservedSize int, requiredSize int) *Buffer {
	return defaultBufferPool.Get(reservedSize, requiredSize)
}

// ResetDefaultBufferPool resets defaultBufferPool.
func ResetDefaultBufferPool() {
	defaultBufferPool.Reset()
}
