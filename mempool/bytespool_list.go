package mempool

import (
	. "github.com/heartchord/goblazer"
)

type bytesBlockList struct {
	root      *BytesBlock // 字节内存块链表根结点
	size      int         // 字节内存块链表结点数
	blockSize int         // 管理的每个内存块大小
}

func newBytesBlockList(blockSize int) (l *bytesBlockList) {
	l = new(bytesBlockList)
	idx := BinarySearchIntsGE(defaultBytesBlockSizeSet, blockSize)
	l.root = newBytesBlock(blockSize, idx)
	l.root.next = l.root
	l.root.prev = l.root
	l.blockSize = blockSize

	for i := 0; i < 100; i++ {
		l.putBack(newBytesBlock(blockSize, idx))
	}
	return
}

func (l *bytesBlockList) getBack() (bb *BytesBlock) {
	r := l.root
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

func (l *bytesBlockList) putBack(e *BytesBlock) (result bool) {
	r := l.root
	if e != nil {
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

func (l *bytesBlockList) getFront() (bb *BytesBlock) {
	r := l.root
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

func (l *bytesBlockList) putFront(e *BytesBlock) (result bool) {
	r := l.root
	if e != nil {
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
