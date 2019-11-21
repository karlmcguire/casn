package casn

import "unsafe"

type casnStatus byte

const (
	undecided casnStatus = iota
	failed
	succeeded
)

type casnDescriptor struct {
	status casnStatus
}

type descriptor struct {
	// control address
	a1 *uint64
	// expected value
	o1 uint64
	// data address
	a2 *uint64
	// old value
	o2 uint64
	// new value
	n2 uint64
}

func (d *descriptor) ptr() uint64 {
	return uint64(uintptr(unsafe.Pointer(d))) | 1<<63
}

func getDescriptor(ptr uint64) *descriptor {
	return (*descriptor)(unsafe.Pointer(uintptr(ptr &^ (1 << 63))))
}

func rdcss(d *descriptor) uint64 {
	r := d.ptr()
	for isDescriptor(r) {
		r, _ = cas(d.a2, d.o2, r)
		if isDescriptor(r) {
			complete(r)
		}
	}
	if r == d.o2 {
		complete(r)
	}
	return r
}

func rdcssRead(d *descriptor) uint64 {
	r := d.ptr()
	for isDescriptor(r) {
		complete(r)
	}
	return r
}

func complete(ptr uint64) {
	d := getDescriptor(ptr)
	v := *(d.a1)
	if v == d.o1 {
		cas(d.a2, ptr, d.n2)
	} else {
		cas(d.a2, ptr, d.o2)
	}
}

func isDescriptor(ptr uint64) bool {
	return ptr>>63 == 1
}

func isCASNDescriptor(ptr uint64) bool {
	return ptr>>62 == 1
}

func cas(ptr *uint64, old, new uint64) (uint64, bool)
