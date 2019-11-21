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

type rdcssDescriptor struct {
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

func (d *rdcssDescriptor) ptr() uint64 {
	return uint64(uintptr(unsafe.Pointer(d))) | 1<<63
}

func getRDCSSDescriptor(ptr uint64) *rdcssDescriptor {
	return (*rdcssDescriptor)(unsafe.Pointer(uintptr(ptr &^ (1 << 63))))
}

func rdcss(d *rdcssDescriptor) uint64 {
	r := d.ptr()
	for isRDCSSDescriptor(r) {
		r, _ = cas(d.a2, d.o2, r)
		if isRDCSSDescriptor(r) {
			complete(r)
		}
	}
	if r == d.o2 {
		complete(r)
	}
	return r
}

func rdcssRead(d *rdcssDescriptor) uint64 {
	r := d.ptr()
	for isRDCSSDescriptor(r) {
		complete(r)
	}
	return r
}

func complete(ptr uint64) {
	d := getRDCSSDescriptor(ptr)
	v := *(d.a1)
	if v == d.o1 {
		cas(d.a2, ptr, d.n2)
	} else {
		cas(d.a2, ptr, d.o2)
	}
}

func isRDCSSDescriptor(ptr uint64) bool {
	return ptr>>63 == 1
}

func isCASNDescriptor(ptr uint64) bool {
	return ptr>>62 == 1
}

func cas(ptr *uint64, old, new uint64) (uint64, bool)
