// share memory by communicating, don't communicate by sharing memory.

package goblazer

import (
	"fmt"
	"time"
)

// memoryBlockSizeSet holds all sizes of memory block with fixed length.
var memoryBlockSizeSet = []int{
	1 * 8,     // 1  byte
	2 * 8,     // 2  byte
	3 * 8,     // 3  byte
	4 * 8,     // 4  byte
	5 * 8,     // 5  byte
	6 * 8,     // 6  byte
	7 * 8,     // 7  byte
	8 * 8,     // 8  byte
	16 * 8,    // 16 byte
	32 * 8,    // 32 byte
	64 * 8,    // 64 byte
	96 * 8,    // 96 byte
	1 * 1024,  // 1  KB
	2 * 1024,  // 2  KB
	3 * 1024,  // 3  KB
	4 * 1024,  // 4  KB
	5 * 1024,  // 5  KB
	6 * 1024,  // 6  KB
	7 * 1024,  // 7  KB
	8 * 1024,  // 8  KB
	16 * 1024, // 16 KB
	32 * 1024, // 32 KB
	64 * 1024, // 64 KB
}

// memoryBlockSizeNum is the length of 'memoryBlockSizeSet'
var memoryBlockSizeNum = len(memoryBlockSizeSet)

const memoryBlockReservedHeadSize = 8
const memoryBlockHoldTime = time.Second * 30
const memoryBlockCleanUpTime = time.Second * 5

// MemoryBlock manages an actual memory block.
type MemoryBlock struct {
	Length    int              // 用户需求内存大小 <=内存块长度len(MemBlock.Buffer) <= 内存块容量cap(MemBlock.Buffer)
	Header    []byte           // 缓冲保留头部空间
	Buffer    []byte           // 用户数据存放缓冲
	buffer    []byte           // 真正数据空间 = mb.Header + mb.Buffer
	flag      int              // 是否在使用中
	allocT    time.Time        // 内存块分配时间戳
	nextBlock *MemoryBlock     // 下一个内存块
	prevBlock *MemoryBlock     // 上一个内存块
	blockList *memoryBlockList // 所属内存块链
}

func (b *MemoryBlock) next() *MemoryBlock {
	if p := b.nextBlock; b.blockList != nil && p != b.blockList.listRoot {
		return p
	}
	return nil
}

// newMemoryBlock creates a new MemoryBlock instance.
func newMemoryBlock(blockSize int) *MemoryBlock {
	b := new(MemoryBlock)

	blockSize += memoryBlockReservedHeadSize
	b.buffer = make([]byte, blockSize)
	b.Header = b.buffer[0:memoryBlockReservedHeadSize]
	b.Buffer = b.buffer[memoryBlockReservedHeadSize:blockSize]
	b.Length = 0
	b.allocT = time.Now()

	return b
}

// reset recovers memory block to initial state
func (b *MemoryBlock) reset() {
	b.Length = 0
	b.allocT = time.Now()
}

// memoryBlockList stores memory block with same block-size specified by 'memBlockSizeSet'.
type memoryBlockList struct {
	blockLen   int               // 管理的每个内存块大小
	listRoot   *MemoryBlock      // 内存块链表根结点
	listSize   int               //
	preallocs  int               // 预先分配的内存块个数
	allocChan  chan *MemoryBlock // 分配内存块的管道
	recylChan  chan *MemoryBlock // 回收内存块的管道
	newOpTimes int               // new内存的次数
	allocTimes int               // 分配内存块次数
	recylTimes int               // 回收内存块池数
}

func (l *memoryBlockList) front() *MemoryBlock {
	if l.listSize == 0 {
		return nil
	}
	return l.listRoot.nextBlock
}

func (l *memoryBlockList) push(b *MemoryBlock) *MemoryBlock {
	r := l.listRoot

	if r.nextBlock == nil {
		r.nextBlock = r
		r.prevBlock = r
		l.listSize = 0
	}

	n := r.nextBlock
	r.nextBlock = b
	b.prevBlock = r
	b.nextBlock = n
	n.prevBlock = b
	b.blockList = l
	l.listSize++

	return b
}

func (l *memoryBlockList) pop(b *MemoryBlock) {
	if b.blockList == l {
		b.prevBlock.nextBlock = b.nextBlock // 断链操作
		b.nextBlock.prevBlock = b.prevBlock // 断链操作
		b.nextBlock = nil                   // 避免对象引用造成的内存泄漏
		b.prevBlock = nil                   // 避免对象引用造成的内存泄漏
		b.blockList = nil                   // 避免对象引用造成的内存泄漏
		l.listSize--                        // 长度减一
	}
}

// recycle returns the memory block to memory block list.
func (l *memoryBlockList) recycle(mb *MemoryBlock) {
	l.recylChan <- mb
	l.recylTimes++
}

// allocate fetches a new memory block from memory block list.
func (l *memoryBlockList) allocate() *MemoryBlock {
	mb := <-l.allocChan
	l.allocTimes++
	return mb
}

// workCoroutine is responsible to allocate, recycle and abandon memory blocks
func (l *memoryBlockList) workCoroutine() {
	for {
		if l.listSize == 0 {
			b := newMemoryBlock(l.blockLen)
			l.push(b)
			l.newOpTimes++
		}

		f := l.front()

		select {
		case b := <-l.recylChan:
			{ // 回收操作
				b.reset()
				l.push(b)
			}
		case l.allocChan <- f:
			{ // 分配操作
				l.pop(f)
			}
		}
	}
}

func newMemBlockList(blockSize int) *memoryBlockList {
	l := new(memoryBlockList)

	l.blockLen = blockSize

	l.listRoot = newMemoryBlock(l.blockLen)
	l.listRoot.nextBlock = l.listRoot
	l.listRoot.prevBlock = l.listRoot
	l.listSize = 0

	// 创建内存分配和回收管道
	l.allocChan = make(chan *MemoryBlock, 10)
	l.recylChan = make(chan *MemoryBlock, 10)

	// 预先分配内存块
	l.preallocs = 10
	for i := 0; i < l.preallocs; i++ {
		b := newMemoryBlock(l.blockLen)
		l.push(b)
	}

	go l.workCoroutine()

	return l
}

// MemoryPool is a kind of memory pool that uses channel 'Allocate' or 'Recycle' a []byte buffer.
type MemoryPool struct {
	blockListMap map[int]*memoryBlockList
}

// NewMemoryPool :
func NewMemoryPool() *MemoryPool {
	mp := new(MemoryPool)

	// 创建并初始化内存块链表集合实例
	mp.blockListMap = make(map[int]*memoryBlockList)

	for i := 0; i < memoryBlockSizeNum; i++ {
		blockSize := memoryBlockSizeSet[i]
		mp.blockListMap[blockSize] = newMemBlockList(memoryBlockSizeSet[i])
	}

	return mp
}

// Allocate :
func (mp *MemoryPool) Allocate(requiredSize int) (*MemoryBlock, bool) {
	// 遍历memBlockSizeSet找到合适的内存块索引
	idx := -1
	for i := 0; i < memoryBlockSizeNum; i++ {
		if requiredSize <= memoryBlockSizeSet[i] {
			idx = memoryBlockSizeSet[i]
			break
		}
	}

	if idx > 0 {
		if list, ok := mp.blockListMap[idx]; ok {
			mb := list.allocate()
			mb.Length = requiredSize
			return mb, true
		}
	}

	return nil, false
}

// Recycle :
func (mp *MemoryPool) Recycle(b *MemoryBlock) bool {
	c := cap(b.buffer) // 容量
	l := len(b.buffer) // 长度

	if l != c { // 非法内存块
		return false
	}

	if blockList, ok := mp.blockListMap[l-memoryBlockReservedHeadSize]; ok {
		blockList.recycle(b)
		return true
	}

	fmt.Println("not recycle!")

	return false
}

// Statistics :
func (mp *MemoryPool) Statistics() {
	for i := 0; i < memoryBlockSizeNum; i++ {
		idx := memoryBlockSizeSet[i]
		if list, ok := mp.blockListMap[idx]; ok {
			fmt.Printf("MemoryPool.blockListMap[%d] - new : %d, alloc : %d, recyl : %d\n", idx, list.newOpTimes, list.allocTimes, list.recylTimes)
		}
	}
}
