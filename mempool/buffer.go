package mempool

import (
	"fmt"
	"sync/atomic"
)

// Buffer wraps a struct based on BytesBlock.
//     Some tips about resSize, oriSize and curSize:
//     1.Generally, we will allocate memory buffer by required data size. Eg. If we have a data with a
//       fixed length 'm', we will allocate a memory buffer by required size 'm', so oriSize = curSize = m.
//     2.When memory buffer allocated, the size of buffer won't be changed, it means oriSize can't be change
//       after allocated.
//     3.But, sometimes, we may want to shorten the data length, we can get to this purpose by shortening curSize
//        and in this case, curSize <= oriSize.
//     4.resSize is allocated when memory buffer is allocating, this reserved sapce can't be changed by any means,
//       and it's invisible when using Buffer normally. But if you want to and some data header to data, you can
//       just access this space and fill the extra data.
type Buffer struct {
	resSize uint32      // reserved size
	oriSize uint32      // original size
	curSize uint32      // current size
	byBlock *BytesBlock // bytes block
}

// NewBuffer create a new instance of Buffer.
func NewBuffer(bp *BytesPool, reservedSize uint32, requiredSize uint32) *Buffer {
	if bp == nil {
		panic("[NewBuffer] expected a BytesPool instance, got nil!")
	}

	return &Buffer{
		resSize: reservedSize,
		oriSize: requiredSize,
		curSize: requiredSize,
		byBlock: bp.Get(int(reservedSize + requiredSize)),
	}
}

// AddRef increases the reference counter. If you want to share this Buffer with others, don't forget invoking this function.
func (b *Buffer) AddRef() int32 {
	if b.byBlock == nil {
		panic("[Buffer.AddRef error] expected Buffer.byBlock != nil, got Buffer.byBlock nil")
	}
	return b.byBlock.AddRef()
}

// DecRef decreases the reference counter. If you want to give back current Buffer to BytesPool, just invoke this function.
func (b *Buffer) DecRef() int32 {
	if b.byBlock == nil {
		panic("[Buffer.AddRef error] expected Buffer.byBlock != nil, got Buffer.byBlock nil")
	}

	ref := b.byBlock.DecRef()
	if ref == int32(0) {
		b.byBlock = nil // avoid leaking
	}

	return ref
}

// OriSize returns the original size of buffer.
func (b *Buffer) OriSize() uint32 {
	return b.oriSize
}

// ResSize returns the reserved size of buffer.
func (b *Buffer) ResSize() uint32 {
	return b.resSize
}

// CurSize returns the current size of buffer.
func (b *Buffer) CurSize() uint32 {
	return b.curSize
}

// SetSize changes the current size of buffer.
func (b *Buffer) SetSize(newSize uint32) {
	if newSize == 0 {
		panic(fmt.Sprintf("[Buffer.SetSize error] expected newSize > 0, got newSize = %d", newSize))
	}

	if newSize > b.oriSize {
		panic(fmt.Sprintf("[Buffer.SetSize error] expected newSize <= Buffer.oriSize, got newSize = %d", newSize))
	}

	atomic.SwapUint32(&b.curSize, newSize)
}

// ResetSize resets the current size of buffer to original.
func (b *Buffer) ResetSize() {
	atomic.SwapUint32(&b.curSize, b.oriSize)
}

// ResBuf returns the reserved space of the buffer.
func (b *Buffer) ResBuf() []byte {
	if b.resSize > uint32(cap(b.byBlock.buf)) {
		panic("[Buffer.ResBuf error] Buffer.resSize is too large!")
	}
	return b.byBlock.buf[0:b.resSize]
}

// Buf returns the user space of the buffer.
func (b *Buffer) Buf() []byte {
	if b.resSize+b.oriSize > uint32(cap(b.byBlock.buf)) {
		panic("[Buffer.GetBuf error] Buffer.resSize or Buffer.oriSize is too large!")
	}

	if b.curSize > b.oriSize {
		panic("[Buffer.GetBuf error] Buffer.curSize is too large!")
	}

	return b.byBlock.buf[b.resSize : b.resSize+b.curSize]
}
