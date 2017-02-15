package goblazer

import (
	"net"
)

// SocketStream :
type SocketStream struct {
	conn net.Conn
	pack chan *MemoryBlock
}

// NewSocketStream :
func NewSocketStream(conn net.Conn) *SocketStream {
	ss := new(SocketStream)

	ss.conn = conn
	ss.pack = make(chan *MemoryBlock, 100)

	go ss.read()

	return ss
}

// Close :
func (ss *SocketStream) Close() {

}

// RecvPackage :
func (ss *SocketStream) RecvPackage() *MemoryBlock {
	if mb, ok := <-ss.pack; ok {
		return mb
	}

	return nil
}

// SendPackage :
func (ss *SocketStream) SendPackage(mb *MemoryBlock) {

}

func (ss *SocketStream) read() {
	defer close(ss.pack)

	for {
		// 先读出包头

		// 在根据包头指定数据长度读出数据包
		for {

		}
	}
}
