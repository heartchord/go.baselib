package goblazer

import (
	"container/list"
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

var memoryBlockHoldTime = time.Second * 30
var memoryBlockCleanUpTime = time.Second * 5

// MemoryBlock manages an actual memory block.
type MemoryBlock struct {
	Length int       // 用户需求内存大小 <=内存块长度len(MemBlock.Buffer) <= 内存块容量cap(MemBlock.Buffer)
	Buffer []byte    // 用户数据存放缓冲
	allocT time.Time // 内存块分配时间戳
}

// newMemoryBlock creates a new MemoryBlock instance.
func newMemoryBlock(blockSize int) *MemoryBlock {
	mb := new(MemoryBlock)

	mb.Length = 0
	mb.Buffer = make([]byte, blockSize) // cap = len
	mb.allocT = time.Now()

	return mb
}

// reset recovers memory block to initial state
func (b *MemoryBlock) reset() {
	b.Length = 0
	b.allocT = time.Now()
}

// memoryBlockList stores memory block with same block-size specified by 'memBlockSizeSet'.
type memoryBlockList struct {
	blockSize  int               // 管理的每个内存块大小
	blockList  *list.List        // 管理的实际内存卡链表
	allocChan  chan *MemoryBlock // 分配内存块的管道
	recylChan  chan *MemoryBlock // 回收内存块的管道
	newOpTimes int               // new内存的次数
	allocTimes int               // 分配内存块次数
	recylTimes int               // 回收内存块池数
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
	timeout := time.NewTimer(memoryBlockCleanUpTime)

	for {
		if l.blockList.Len() == 0 {
			b := newMemoryBlock(l.blockSize)
			l.blockList.PushFront(b)
			l.newOpTimes++
		}

		f := l.blockList.Front()

		select {
		case b := <-l.recylChan:
			{ // 回收操作
				if !timeout.Stop() {
					<-timeout.C
				}

				b.reset()
				l.blockList.PushFront(b)

				timeout.Reset(memoryBlockCleanUpTime)
			}
		case l.allocChan <- f.Value.(*MemoryBlock):
			{ // 分配操作
				if !timeout.Stop() {
					<-timeout.C
				}

				l.blockList.Remove(f)

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

	// 创建内存分配和回收管道
	l.allocChan = make(chan *MemoryBlock)
	l.recylChan = make(chan *MemoryBlock)

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

// Statistics :
func (p *MemoryPool) Statistics() {
	for i := 0; i < memoryBlockSizeNum; i++ {
		idx := memoryBlockSizeSet[i]
		if list, ok := p.blockListMap[idx]; ok {
			fmt.Printf("MemoryPool.blockListMap[%d] - new : %d, alloc : %d, recyl : %d\n", idx, list.newOpTimes, list.allocTimes, list.recylTimes)
		}
	}
}
