package goblazer

import (
	"reflect"
	"unsafe"
)

func GetByteSliceAddress(b []byte) unsafe.Pointer {
	p := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	return unsafe.Pointer(p.Data)
}
