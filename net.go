package goblazer

import (
	"fmt"
	"net"
)

// RecvDataFromConn recvs data with specific 'size' into 'buff'.
//    parameter    : conn - 套接字连接, buff - 数据接收缓冲, size - 期望接收数据长度
//    return value : recvd - 接收数据长度, ret - 执行结果，如果发生错误，返回false
func RecvDataFromConn(conn net.Conn, buff []byte, size int) (recvd int, result bool) {
	if size < 0 || len(buff) < size { // 如果期望长度为负或数据接收缓存不够
		result = false
		return
	}

	torecv := size // 还需收取数据总字节数
	for torecv > 0 {
		b := buff[recvd : recvd+torecv] // 剩余空间
		n, err := conn.Read(b)
		if err != nil { // 读取数据错误
			result = false
			return
		}

		recvd += n
		torecv -= n
	}

	result = true
	return
}

// SendDataToConn sends data with specific 'size' into 'buff'.
//    parameter    : conn - 套接字连接, buff - 数据发送缓冲, size - 期望发送数据长度
//    return value : sent - 发送数据长度, ret - 执行结果，如果发生错误，返回false
func SendDataToConn(conn net.Conn, buff []byte, size int) (sent int, result bool) {
	if size < 0 || len(buff) < size { // 如果期望长度为负或数据接发送缓存不够
		result = false
		return
	}

	tosend := size // 还需发送数据总字节数
	for tosend > 0 {
		b := buff[sent : sent+tosend] // 剩余空间
		n, err := conn.Write(b)
		if err != nil { // 写入数据错误
			result = false
			return
		}

		sent += n
		tosend -= n
	}

	result = true
	return
}

var tempPool = NewMemoryPool()
var defaultPackHeaderSize = 2

// SocketStream receives complete protocol packages into recv channel and sends complete protocol packages from send channel.
type SocketStream struct {
	sockConn       net.Conn
	recvChan       chan *MemoryBlock
	sendChan       chan *MemoryBlock
	packHeaderSize int
	recvChanLength int
	sendChanLength int
}

// NewSocketStream creates a new SocketStream instance.
func NewSocketStream(conn net.Conn) *SocketStream {
	if conn == nil {
		return nil
	}

	ss := new(SocketStream)

	ss.sockConn = conn
	ss.recvChan = make(chan *MemoryBlock, 100)
	ss.sendChan = make(chan *MemoryBlock, 100)
	ss.packHeaderSize = defaultPackHeaderSize

	go ss.recvCoroutine()

	return ss
}

// SetPackHeaderSize sets the package header size.
func (ss *SocketStream) SetPackHeaderSize(packHeaderSize int) {
	ss.packHeaderSize = packHeaderSize
}

// GetPackHeaderSize returns the package header size.
func (ss *SocketStream) GetPackHeaderSize() int {
	return ss.packHeaderSize
}

// Close :
func (ss *SocketStream) Close() {

}

// RecvPackage :
func (ss *SocketStream) RecvPackage() *MemoryBlock {
	if mb, ok := <-ss.recvChan; ok {
		return mb
	}

	return nil
}

// SendPackage :
func (ss *SocketStream) SendPackage(mb *MemoryBlock) {
	ss.sendChan <- mb
}

func (ss *SocketStream) recvCoroutine() {
	header, ok := tempPool.Allocate(ss.packHeaderSize)
	if !ok {
		// err
	}
	defer tempPool.Recycle(header)

	for {
		// 先读出包头
		n, ok := RecvDataFromConn(ss.sockConn, header.Buffer, ss.packHeaderSize)
		if !ok || n != ss.packHeaderSize {
			// err
		}

		packBodySize := int(BytesToUint32(header.Buffer))
		fmt.Println(packBodySize)

		body, ok := tempPool.Allocate(packBodySize)
		if !ok {
			// error
		}

		n, ok = RecvDataFromConn(ss.sockConn, body.Buffer, packBodySize)
		if !ok || n != packBodySize {
			// err
		}

		ss.recvChan <- body
	}
}

func (ss *SocketStream) sendCoroutine() {
	for {
		// 先读出包头

		// 在根据包头指定数据长度读出数据包
		for {

		}
	}
}
