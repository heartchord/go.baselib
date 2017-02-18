package goblazer

import "encoding/binary"

// BytesToUint32 converts bytes to uint32. If the length of bytes is smaller than 4, the insufficient space will
//     be filled with 0; if the length of bytes is bigger than 4, the first 4 bytes will be used to calculate.
func BytesToUint32(bs []byte) uint32 {
	var b0, b1, b2, b3 byte

	l := len(bs)
	if l >= 1 {
		b0 = bs[0]
	}
	if l >= 2 {
		b1 = bs[1]
	}
	if l >= 3 {
		b2 = bs[2]
	}
	if l >= 4 {
		b3 = bs[3]
	}

	return uint32(b0) | uint32(b1)<<8 | uint32(b2)<<16 | uint32(b3)<<24
}

func BytesToUint64(bs []byte) uint64 {
	var b0, b1, b2, b3, b4, b5, b6, b7 byte

	l := len(bs)
	if l >= 1 {
		b0 = bs[0]
	}
	if l >= 2 {
		b1 = bs[1]
	}
	if l >= 3 {
		b2 = bs[2]
	}
	if l >= 4 {
		b3 = bs[3]
	}
	if l >= 5 {
		b4 = bs[4]
	}
	if l >= 6 {
		b5 = bs[5]
	}
	if l >= 7 {
		b6 = bs[6]
	}
	if l >= 8 {
		b7 = bs[7]
	}

	return uint64(b0) | uint64(b1)<<8 | uint64(b2)<<16 | uint64(b3)<<24 | uint64(b4)<<32 | uint64(b5)<<40 | uint64(b6)<<48 | uint64(b7)<<56
}

// BytesToInt32 converts bytes to uint32. If the length of bytes is smaller than 4, the insufficient space will
//     be filled with 0; if the length of bytes is bigger than 4, the first 4 bytes will be used to calculate.
func BytesToInt32(bs []byte) int32 {
	return int32(BytesToUint32(bs))
}

func BytesToInt64(bs []byte) int64 {
	return int64(BytesToUint64(bs))
}

func BytesToNumber(bs []byte, order binary.ByteOrder, data interface{}) bool {
	buffLen := len(bs)
	dataLen := 0

	switch data.(type) {
	case *int8, *uint8:
		dataLen = 1
	case *int16, *uint16:
		dataLen = 2
	case *int32, *uint32:
		dataLen = 4
	case *int64, *uint64:
		dataLen = 8
	}

	if dataLen <= 0 || buffLen < dataLen {
		return false
	}

	switch data := data.(type) {
	case *int8:
		*data = int8(bs[0])
	case *uint8:
		*data = bs[0]
	case *int16:
		*data = int16(order.Uint16(bs))
	case *uint16:
		*data = order.Uint16(bs)
	case *int32:
		*data = int32(order.Uint32(bs))
	case *uint32:
		*data = order.Uint32(bs)
	case *int64:
		*data = int64(order.Uint64(bs))
	case *uint64:
		*data = order.Uint64(bs)
	default:
		return false
	}

	return true
}
