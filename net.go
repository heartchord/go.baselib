package goblazer

import (
	"net"
)

// 临时内存池方案，待替换
var tempPool = NewMemoryPool()

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
		b := buff[recvd : recvd+torecv] // 剩余数据空间
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

const (
	defaultPackHeaderSize = 2
	defaultRecvChanLength = 100
	defaultSendChanLength = 100
)

const (
	socketStreamStatusInit = iota
	socketStreamStatusStart
	socketStreamStatusStop
)

// SocketStream receives complete protocol packages into recv channel and sends complete protocol packages from send channel.
type SocketStream struct {
	sockStat       int
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

	ss.sockStat = socketStreamStatusInit
	ss.sockConn = conn

	ss.packHeaderSize = defaultPackHeaderSize
	ss.recvChanLength = defaultRecvChanLength
	ss.sendChanLength = defaultSendChanLength

	ss.recvChan = make(chan *MemoryBlock, ss.recvChanLength)
	ss.sendChan = make(chan *MemoryBlock, ss.sendChanLength)

	return ss
}

// SetPackHeaderSize sets the package header size. This function must be invoked before ss.Start().
func (ss *SocketStream) SetPackHeaderSize(packHeaderSize int) (result bool) {
	if ss.sockStat != socketStreamStatusStart {
		ss.packHeaderSize = packHeaderSize
		result = true
	}
	return
}

// GetPackHeaderSize returns the package header size.
func (ss *SocketStream) GetPackHeaderSize() int {
	return ss.packHeaderSize
}

// GetStatus returns the status of current SocketStream.
func (ss *SocketStream) GetStatus() int {
	return ss.sockStat
}

// Start starts recv and send coroutine to recv or send complete package.
func (ss *SocketStream) Start() {
	if ss.sockStat != socketStreamStatusStart {
		go ss.recvCoroutine()
		go ss.sendCoroutine()
		ss.sockStat = socketStreamStatusStart
	}
}

// Close :
func (ss *SocketStream) Close() {

}

// RecvPackage gets a complete protocol package from current SocketStream.
func (ss *SocketStream) RecvPackage() (pack *MemoryBlock) {
	if ss.sockStat != socketStreamStatusStart {
		return
	}

	if mb, ok := <-ss.recvChan; ok {
		pack = mb
	}
	return
}

// SendPackage :
func (ss *SocketStream) SendPackage(mb *MemoryBlock) (result bool) {
	if ss.sockStat != socketStreamStatusStart {
		return
	}

	ss.sendChan <- mb
	result = true
	return
}

func (ss *SocketStream) recvCoroutine() {
	var ok bool
	var n, packSize int
	var header, pack *MemoryBlock

	if header, ok = tempPool.Allocate(ss.packHeaderSize); !ok {
		goto Exit0
	}

	for {
		// 读取包头数据（包头存储后续完整包长度）
		if n, ok = RecvDataFromConn(ss.sockConn, header.Buffer, ss.packHeaderSize); !ok || n != ss.packHeaderSize {
			goto Exit0
		}

		// 获取完整数据包长度
		if packSize = int(BytesToUint32(header.Buffer)); packSize <= 0 {
			goto Exit0
		}

		// 分配数据包空间
		if pack, ok = tempPool.Allocate(packSize); !ok {
			goto Exit0
		}

		// 读取完整数据包
		if n, ok = RecvDataFromConn(ss.sockConn, pack.Buffer, packSize); !ok || n != packSize {
			goto Exit0
		}

		// 转移内存块所有权
		ss.recvChan <- pack
		pack = nil
	}

Exit0:
	if header != nil {
		tempPool.Recycle(header)
	}
	if pack != nil {
		tempPool.Recycle(pack)
	}
}

func (ss *SocketStream) sendCoroutine() {
	for {
		pack := <-ss.sendChan
		n, ok := SendDataToConn(ss.sockConn, pack.Buffer, pack.Length)
		if !ok || n != pack.Length { // error
		}
	}
}

// SocketAcceptor :
type SocketAcceptor struct {
	listener net.Listener
}

// NewSocketAcceptor :
func NewSocketAcceptor(network string, addr string) *SocketAcceptor {
	l, err := net.Listen(network, addr)
	if err != nil {
		return nil
	}

	sa := new(SocketAcceptor)
	sa.listener = l

	return sa
}

//func (sa *SocketAcceptor) Accept() *SocketStream {
//
//}
