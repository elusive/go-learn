package main

import (
	"fmt"
	"unsafe"
)

type _type struct {
	size       uintptr
	ptrdata    uintptr
	hash       uint32
	align      uint8
	fieldalign uint8
	kind       uint8
}

type functype struct {
	typ      _type
	inCount  uint16
	outCount uint16
}

type eface struct {
	_type *_type
	data  unsafe.Pointer
}

func out(str string) uintptr {
	fmt.Println("out: ", str)
	code := uintptr(1)
	return code
}

func main() {
	fmt.Println("Entering sandbox...")
	fn := (*eface)(unsafe.Pointer(out))
	ft := (*functype)(unsafe.Pointer(fn._type))
	fmt.Println("function type: ", ft)
	uintptrSize := unsafe.Sizeof(uintptr(0))

	if ft.out()[0].size != uintptrSize {
		fmt.Println("size does not MATCH!")
	}

	fmt.Println("exiting sandbox")
}
