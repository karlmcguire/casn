package casn

import (
	"sync/atomic"
	"unsafe"
)

type Update struct {
	Address *uint64
	Old     uint64
	New     uint64
}

func CASN(updates []Update) bool {
	return casn(&casnDescriptor{undecided, updates})
}

func CASNRead(addr *uint64) uint64 {
	return casnRead(addr)
}

const (
	undecided uint64 = iota
	failed
	succeeded
)

type casnDescriptor struct {
	status  uint64
	updates []Update
}

func (d *casnDescriptor) ptr() uint64 {
	return uint64(uintptr(unsafe.Pointer(d))) | 1<<62
}

func getCASNDescriptor(ptr uint64) *casnDescriptor {
	return (*casnDescriptor)(unsafe.Pointer(uintptr(ptr &^ (1 << 62))))
}

func isCASNDescriptor(ptr uint64) bool {
	return ptr>>62 == 1
}

func casn(cd *casnDescriptor) bool {
	if cd.status == undecided {
		status := succeeded
		descs := make([]*rdcssDescriptor, 0, cap(cd.updates))
		for i := 0; i < len(cd.updates) && status == succeeded; i++ {
		retry:
			descs = append(descs, &rdcssDescriptor{
				a1: &cd.status,
				o1: undecided,
				a2: cd.updates[i].Address,
				o2: cd.updates[i].Old,
				n2: cd.ptr(),
			})
			val := rdcss(descs[i])
			if isCASNDescriptor(val) {
				if val != cd.ptr() {
					casn(getCASNDescriptor(val))
					goto retry
				}
			} else if val != cd.updates[i].Old {
				status = failed
			}
		}
		cas(&cd.status, undecided, status)
	}
	success := cd.status == succeeded
	for i := 0; i < len(cd.updates); i++ {
		new := uint64(0)
		if success {
			new = cd.updates[i].New
		} else {
			new = cd.updates[i].Old
		}
		cas(cd.updates[i].Address, cd.ptr(), new)
	}
	return success
}

func casnRead(addr *uint64) uint64 {
	var r uint64
	for {
		r = atomic.LoadUint64(addr)
		if isCASNDescriptor(r) {
			casn(getCASNDescriptor(r))
		} else {
			break
		}
	}
	return r
}

type rdcssDescriptor struct {
	a1 *uint64
	o1 uint64
	a2 *uint64
	o2 uint64
	n2 uint64
}

func (d *rdcssDescriptor) ptr() uint64 {
	return uint64(uintptr(unsafe.Pointer(d))) | 1<<63
}

func getRDCSSDescriptor(ptr uint64) *rdcssDescriptor {
	return (*rdcssDescriptor)(unsafe.Pointer(uintptr(ptr &^ (1 << 63))))
}

func isRDCSSDescriptor(ptr uint64) bool {
	return ptr>>63 == 1
}

func rdcss(d *rdcssDescriptor) uint64 {
	var r uint64
	for {
		r = cas(d.a2, d.o2, d.ptr())
		if isRDCSSDescriptor(r) {
			complete(getRDCSSDescriptor(r))
		} else {
			break
		}
	}
	if r == d.o2 {
		complete(d)
	}
	return r
}

func rdcssRead(addr *uint64) uint64 {
	var r uint64
	for {
		r = atomic.LoadUint64(addr)
		if isRDCSSDescriptor(r) {
			complete(getRDCSSDescriptor(r))
		} else {
			break
		}
	}
	return r
}

func complete(d *rdcssDescriptor) {
	a1 := atomic.LoadUint64(d.a1)
	if a1 == d.o1 {
		cas(d.a2, d.ptr(), d.n2)
	} else {
		cas(d.a2, d.ptr(), d.o2)
	}
}

func cas(ptr *uint64, old, new uint64) uint64
