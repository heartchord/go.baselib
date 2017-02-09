package goblazer

import (
	"container/list"
	"fmt"
	"time"
)

var memoryBlockSizeSet = []int{
	1 * 8,     // 1 byte
	2 * 8,     // 2 byte
	3 * 8,     // 3 byte
	4 * 8,     // 4 byte
	5 * 8,     // 5 byte
	6 * 8,     // 6 byte
	7 * 8,     // 7 byte
	8 * 8,     // 8 byte
	16 * 8,    // 16 byte
	32 * 8,    // 32 byte
	64 * 8,    // 64 byte
	96 * 8,    // 96 byte
	1 * 1024,  // 1 KB
	2 * 1024,  // 2 KB
	3 * 1024,  // 3 KB
	4 * 1024,  // 4 KB
	5 * 1024,  // 5 KB
	6 * 1024,  // 6 KB
	7 * 1024,  // 7 KB
	8 * 1024,  // 8 KB
	16 * 1024, // 16 KB
	32 * 1024, // 32 KB
	64 * 1024, // 64 KB
}

var memoryBlockSizeSetSize = len(memoryBlockSizeSet)
var memoryBlockHoldTime = time.Second * 30
var memoryBlockCleanUpTime = time.Second * 5

// MemoryBlock stores an actual memory block.
type MemoryBlock struct {
	Length int       // 用户需求内存大小（内存块的长度和容量可以通过len(MemBlock.Buffer)和cap(MemBlock.Buffer)求出）
	Buffer []byte    // 用户数据存放缓冲
	allocT time.Time // 内存块分配时间戳
}

func newMemoryBlock(blockSize int) *MemoryBlock {
	b := new(MemoryBlock)

	b.Length = 0
	b.Buffer = make([]byte, blockSize)
	b.allocT = time.Now()

	return b
}

func (b *MemoryBlock) reset() {
	b.Length = 0
	b.allocT = time.Now()
}

// memoryBlockList stores memory block with specified block size by 'memBlockSizeSet'
type memoryBlockList struct {
	blockSize    int               // 内存块链表中每个内存块大小
	blockList    *list.List        // 内存块链表实际存储数据
	alloc        chan *MemoryBlock // 分配channel
	recyl        chan *MemoryBlock // 回收channel
	newOpTimes   int               // new内存次数
	allocOpTimes int               // 分配内存块次数
	recylOpTimes int               // 回收内存块池数
}

func (l *memoryBlockList) recycle(b *MemoryBlock) {
	l.recyl <- b
}

func (l *memoryBlockList) allocate() *MemoryBlock {
	return <-l.alloc
}

func (l *memoryBlockList) workerCoroutine() {
	timeout := time.NewTimer(memoryBlockCleanUpTime)

	for {
		if l.blockList.Len() == 0 {
			b := newMemoryBlock(l.blockSize)
			l.blockList.PushFront(b)
			l.newOpTimes++
		}

		f := l.blockList.Front()

		select {
		case b := <-l.recyl:
			{ // 回收操作
				if !timeout.Stop() {
					<-timeout.C
				}

				b.reset()
				l.blockList.PushFront(b)
				l.recylOpTimes++

				timeout.Reset(memoryBlockCleanUpTime)
			}
		case l.alloc <- f.Value.(*MemoryBlock):
			{ // 分配操作
				if !timeout.Stop() {
					<-timeout.C
				}

				l.blockList.Remove(f)
				l.allocOpTimes++

				timeout.Reset(memoryBlockCleanUpTime)
			}
		case <-timeout.C: // 空闲时段进行回收
			{
				b := l.blockList.Front()
				for b != nil {
					n := b.Next()
					if time.Since(b.Value.(*MemoryBlock).allocT) > memoryBlockHoldTime {
						l.blockList.Remove(b)
						b.Value = nil
					}
					b = n
				}
				timeout.Reset(memoryBlockCleanUpTime)
			}
		}
	}
}

func newMemBlockList(blockSize int) *memoryBlockList {
	l := new(memoryBlockList)

	l.blockSize = blockSize
	l.blockList = list.New()

	// 创建并初始化内存分配和回收管道
	l.alloc = make(chan *MemoryBlock)
	l.recyl = make(chan *MemoryBlock)

	go l.workerCoroutine()

	return l
}

// MemoryPool is a kind of memory pool that uses channel 'Allocate' or 'Recycle' a []byte buffer.
type MemoryPool struct {
	blockListMap map[int]*memoryBlockList
}

// NewMemoryPool :
func NewMemoryPool() *MemoryPool {
	return new(MemoryPool)
}

// InitPool must be invoked before using ChanMemPool
func (p *MemoryPool) InitPool() {
	// 创建并初始化内存块链表集合实例
	p.blockListMap = make(map[int]*memoryBlockList)

	for i := 0; i < memoryBlockSizeSetSize; i++ {
		blockSize := memoryBlockSizeSet[i]
		p.blockListMap[blockSize] = newMemBlockList(memoryBlockSizeSet[i])
	}
}

// Allocate :
func (p *MemoryPool) Allocate(requiredSize int) (*MemoryBlock, bool) {
	// 遍历memBlockSizeSet找到合适的内存块索引
	idx := -1
	for i := 0; i < memoryBlockSizeSetSize; i++ {
		if requiredSize <= memoryBlockSizeSet[i] {
			idx = memoryBlockSizeSet[i]
			break
		}
	}

	if idx > 0 {
		if list, ok := p.blockListMap[idx]; ok {
			mb := list.allocate()
			mb.Length = requiredSize
			return mb, true
		}
	}

	return nil, false
}

// Recycle :
func (p *MemoryPool) Recycle(b *MemoryBlock) bool {
	c := cap(b.Buffer) // 容量
	l := len(b.Buffer) // 长度

	if l != c { // 非法内存块
		return false
	}

	if blockList, ok := p.blockListMap[l]; ok {
		blockList.recycle(b)
		return true
	}

	return false
}

func (p *MemoryPool) Statistics() {
	for i := 0; i < memoryBlockSizeSetSize; i++ {
		idx := memoryBlockSizeSet[i]
		if list, ok := p.blockListMap[idx]; ok {
			fmt.Printf("MemoryPool.blockListMap[%d] - new : %d, alloc : %d, recyl : %d\n", idx, list.newOpTimes, list.allocOpTimes, list.recylOpTimes)
		}
	}
}

func init() {
}
