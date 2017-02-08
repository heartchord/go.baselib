package goblazer

import (
	"container/list"
	"time"
)

// Mem : Memory
var memBlockSizeSet = []int{
	1*8 + 32,
	2*8 + 32,
	3*8 + 32,
	4*8 + 32,
	5*8 + 32,
	6*8 + 32,
	7*8 + 32,
	8*8 + 32,
	16*8 + 32,
	32*8 + 32,
	64*8 + 32,
	96*8 + 32,
	1*1024 + 32,
	2*1024 + 32,
	3*1024 + 32,
	4*1024 + 32,
	5*1024 + 32,
	6*1024 + 32,
	7*1024 + 32,
	8*1024 + 32,
	16*1024 + 32,
	32*1024 + 32,
	64*1024 + 32,
}

var memBlockSizeSetSize = len(memBlockSizeSet)

func init() {
}

// memBlock is a struct stores actual memory block
type memBlock struct {
	when time.Time // 分配时时间戳
	size int       // 实际使用大小
	data []byte    // 数据使用缓冲
}

func newMemBlock(maxSize int) *memBlock {
	b := new(memBlock)

	b.when = time.Now()
	b.size = 0
	b.data = make([]byte, maxSize)

	return b
}

// memBlockList is a struct stores memory block with specified block size by 'memBlockSizeSet'
type memBlockList struct {
	eachBlockSize int        // 内存块链表中每个内存块大小
	blockList     *list.List // 内存块链表实际存储数据
	alloc         chan []byte
	recyl         chan []byte
}

func (l *memBlockList) workerCoroutine() {
	timeout := time.NewTimer(time.Second * 5)

	for {
		if l.blockList.Len() == 0 {
			b := newMemBlock(l.eachBlockSize)
			l.blockList.PushFront(b)
		}

		f := l.blockList.Front()

		select {
		case data := <-l.recyl: // 回收操作

			timeout.Stop()

			b := new(memBlock)
			b.when = time.Now()
			b.size = 0
			b.data = data
			l.blockList.PushFront(b)

			timeout.Reset(time.Second * 5)

		case l.alloc <- f.Value.(*memBlock).data: // 分配操作

			timeout.Stop()
			l.blockList.Remove(f)
			timeout.Reset(time.Second * 5)

		case <-timeout.C: // 空闲时段进行回收

			b := l.blockList.Front()
			for b != nil {
				n := b.Next()
				if time.Since(b.Value.(*memBlock).when) > time.Second*5 {
					l.blockList.Remove(b)
					b.Value = nil
				}
				b = n
			}
			timeout.Reset(time.Second * 5)
		}
	}
}

func newMemBlockList(maxSize int) *memBlockList {
	l := new(memBlockList)

	l.eachBlockSize = maxSize
	l.blockList = list.New()

	// 创建并初始化内存分配和回收管道
	l.alloc = make(chan []byte)
	l.recyl = make(chan []byte)

	go l.workerCoroutine()

	return l
}

// ChanMemPool is a kind of memory pool that uses channel 'Allocate' or 'Recycle' a []byte buffer.
type ChanMemPool struct {
	blockListSet []*memBlockList
}

// InitPool must be invoked before using ChanMemPool
func (p *ChanMemPool) InitPool() {
	// 创建并初始化内存块链表集合实例
	p.blockListSet = make([]*memBlockList, memBlockSizeSetSize)
	for i := 0; i < memBlockSizeSetSize; i++ {
		p.blockListSet[i] = newMemBlockList(memBlockSizeSet[i])
	}
}

/*func (p *ChanMemPool) Allocate(requiredSize int) []byte {

}*/

func (p *ChanMemPool) Recycle(data []byte) bool {
	capacity := cap()
	return true
}
