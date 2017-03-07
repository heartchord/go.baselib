package mempool

import (
	"fmt"
	"sync/atomic"
)

// poolClasses defines all types of bytes-block with fixed-length.
var poolClasses = []int{
	2 * 8,     // 16    bytes        <00-01>
	4 * 8,     // 32    bytes        <01-02>
	6 * 8,     // 48    bytes        <02-03>
	8 * 8,     // 64    bytes        <03-04>
	12 * 8,    // 96    bytes        <04-05>
	16 * 8,    // 128   bytes        <05-06>
	20 * 8,    // 160   bytes        <06-07>
	24 * 8,    // 192   bytes        <07-08>
	32 * 8,    // 256   bytes        <08-09>
	40 * 8,    // 320   bytes        <09-10>
	48 * 8,    // 384   bytes        <10-11>
	56 * 8,    // 448   bytes        <11-12>
	64 * 8,    // 512   bytes        <12-13>
	96 * 8,    // 768   bytes        <13-14>
	1 * 1024,  // 1024  bytes(1  KB) <14-15>
	2 * 1024,  // 2048  bytes(2  KB) <15-16>
	3 * 1024,  // 3072  bytes(3  KB) <16-17>
	4 * 1024,  // 4096  bytes(4  KB) <17-18>
	5 * 1024,  // 5120  bytes(5  KB) <18-19>
	6 * 1024,  // 6144  bytes(6  KB) <18-20>
	7 * 1024,  // 7168  bytes(7  KB) <20-21>
	8 * 1024,  // 8192  bytes(8  KB) <21-22>
	16 * 1024, // 16384 bytes(16 KB) <22-23>
	32 * 1024, // 32768 bytes(32 KB) <23-24>
	64 * 1024, // 65536 bytes(64 KB) <24-25>
}

var poolClassCount = len(poolClasses)

// BytesBlock represents a bytes buffer allocated from BytesPool. User can control its life circle manually, if you want to return it
//     to BytesPool, just invoke BytesBlock.DecRef() to reduce its reference counter to zero. Notice that if you want to share a BytesBlock
//     with others, you must invoke BytesBlock.AddRef() to add its reference, otherwise it will be recycled unexpectedly.
type BytesBlock struct {
	owner *BytesPool // belonger that manages all BytesBlock with the same buffer size. Eg, if cap(BytesBlock.buf)=1024, the class of BytesPool it belonged to is 1024 too.
	buf   []byte     // required buffer, notice that len(bb.buf) <= cap(bb.buf), eg: original-size = 1024, user wants 1000, bb.buf = bb.buf[0:1000]
	ref   int32      // reference counter
}

// Buf returns the real buffer([]byte).
func (bb *BytesBlock) Buf() []byte {
	return bb.buf
}

// AddRef increases the reference counter. If you want to share this BytesBlock with others, don't forget invoking this function.
func (bb *BytesBlock) AddRef() int32 {
	ref := atomic.AddInt32(&bb.ref, 1)

	if ref <= 0 { // error
		panic(fmt.Sprintf("[BytesBlock.AddRef error] expected BytesBlock.ref > 0, got BytesBlock.ref = %d", bb.ref))
	}

	return ref
}

// DecRef decreases the reference counter. If you want to give back current BytesBlock to BytesPool, just invoke this function.
func (bb *BytesBlock) DecRef() int32 {
	ref := atomic.AddInt32(&bb.ref, -1)

	if ref < 0 {
		panic(fmt.Sprintf("[BytesBlock.DecRef error] expected BytesBlock.ref >= 0, got BytesBlock.ref = %d", bb.ref))
	}

	if ref == 0 {
		bb.owner.put(bb)
	}

	return ref
}

// BytesPool holds all BytesBlocks with specified block size. The BytesBlocks with same size will be placed in one Pool.
//    Eg, All 1024 BytesBlocks will be placed in the pool with the class type 1024. The pools with different class make up
//    the final BytesPool.
type BytesPool struct {
	classes  []int   // represents all kinds of Pool, default classes can be specified by 'poolClasses'
	allPools []*Pool // represents all Pools with different class.
}

// Get fetches one BytesBlock fromm BytesPool. The parameter 'requiredSize' represents the buffer size user demands, so requiredSize
//    must <= cap(BytesBlock.buf). When user applies one block, BytesPool will make slicing operation to make sure the space user got
//    is just sufficient.
func (bp *BytesPool) Get(requiredSize int) (bb *BytesBlock) {
	if requiredSize < 0 {
		panic(fmt.Sprintf("[BytesPool.Get error] expected requiredSize >= 0, got requiredSize = %d", requiredSize))
	}

	idx := bp.getClassIdx(requiredSize)
	if idx >= 0 {
		bb = bp.allPools[idx].Get().(*BytesBlock)
		bb.buf = bb.buf[:requiredSize] // reduce the buffer length to user required.
		bb.ref = 1                     // the reference counter is 1 when allocated.
		return
	}

	// New operation. The class can't be found, so this requiredSize is biger than all blockSize in bp.classes.
	// When recycled, because the cap is bigger than all pools, the BytesBlock will just be dropped to the floor.
	bb = &BytesBlock{
		owner: bp,
		ref:   1,
		buf:   make([]byte, requiredSize),
	}
	return
}

// Reset will set the status of all pools to initialization.
func (bp *BytesPool) Reset() {
	for _, pool := range bp.allPools {
		pool.Reset()
	}
}

// put returns the BytesBlock to the appropriate pool, if any appropriate pool can't be found, it will just be dropped to the floor.
func (bp *BytesPool) put(bb *BytesBlock) {
	idx := bp.getClassIdx(cap(bb.buf))
	if idx >= 0 {
		bp.allPools[idx].Put(bb)
	}
}

func checkAllClasses(classes []int) bool {
	l := len(classes)

	if l <= 0 {
		return false
	}

	for i := 0; i < l; i++ {
		if classes[i] <= 0 {
			return false
		}

		if i+1 < l && classes[i] >= classes[i+1] {
			return false
		}
	}

	return true
}

// getClassIdx will find the appropriate index of BytesPool.classes by 'requiredSize'.
//     It will search the first index where BytesPool.classes[idx] >= requiredSize.
//func (bp *BytesPool) getClassIdx(requiredSize int) (idx int) {
//	idx = goblazer.BinarySearchIntsGE(bp.classes, requiredSize)
//	return
//}

// getClassIdx will find the appropriate index of BytesPool.classes by 'requiredSize'.
func (bp *BytesPool) getClassIdx(requiredSize int) (idx int) {
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

// NewBytesPool create a new BytesPool.
//     classes      : specify all kinds of pool. You can use default 'poolClasses' or specify a new classes.
//     localPoolCap : specify the local pool capacity of every 'P'. It specifies the max value of the sum of
//                    pool.localPool.private + pool.localPool.shared.
func NewBytesPool(classes []int, localPoolCap int) *BytesPool {
	if !checkAllClasses(classes) {
		panic(fmt.Sprintf("[NewBytesPool error] expected a ascend classes, got classes is %#v", classes))
	}

	// create an instance.
	classCount := len(classes)
	bp := &BytesPool{
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
				return &BytesBlock{
					owner: bp,
					ref:   0,
					buf:   make([]byte, l),
				}
			},
		)
	}

	return bp
}

var defaultBytesPool *BytesPool

func init() {
	defaultBytesPool = NewBytesPool(poolClasses, 8)
}

// GetFromDefaultBytesPool fetches one BytesBlock from defaultBytesPool.
func GetFromDefaultBytesPool(requiredSize int) *BytesBlock {
	return defaultBytesPool.Get(requiredSize)
}

// ResetDefaultBytesPool resets defaultBytesPool.
func ResetDefaultBytesPool() {
	defaultBytesPool.Reset()
}
