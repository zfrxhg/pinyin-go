package main

/*
#include <string.h>
*/
import "C"
import (
	"reflect"
	"unsafe"
)

// GoString returns a string which address is same as p.
func GoString(p *C.char) (result string) {
	header := (*reflect.StringHeader)(unsafe.Pointer(&result))
	header.Data = uintptr(unsafe.Pointer(p))
	header.Len = int(C.strlen(p))
	return
}

// GoBytes returns a byte array which address is same as p and length is size.
func GoBytes(p unsafe.Pointer, size C.size_t) (result []byte) {
	header := (*reflect.SliceHeader)(unsafe.Pointer(&result))
	header.Data = uintptr(p)
	header.Len = int(size)
	header.Cap = header.Len
	return
}

// MemoryCopy copies go bytes into C memory and returns the count of bytes has been copied.
func MemoryCopy(dst unsafe.Pointer, dstCap C.size_t, src []byte) C.size_t {
	srcLen := C.size_t(len(src))
	if srcLen == 0 {
		return 0
	}
	if srcLen > dstCap {
		srcLen = dstCap
	}
	C.memcpy(dst, unsafe.Pointer(&src[0]), srcLen)
	return srcLen
}

// MemoryAllocateAndCopy allocates C memory and copies go bytes into it. Then it returns the count of bytes has been copied.
func MemoryAllocateAndCopy(dst unsafe.Pointer, src []byte) C.size_t {
	srcLen := C.size_t(len(src))
	if srcLen == 0 {
		return 0
	}
	p := (*unsafe.Pointer)(dst)
	*p = C.CBytes(src)
	return srcLen
}
