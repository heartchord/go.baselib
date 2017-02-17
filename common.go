package goblazer

import "encoding/binary"

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
