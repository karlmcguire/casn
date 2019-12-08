package casn

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestGetRDCSSDescriptor(t *testing.T) {
	data := []uint64{0, 0}
	d := &rdcssDescriptor{
		a1: &data[0],
		o1: 0,
		a2: &data[1],
		o2: 0,
		n2: 1,
	}
	if d != getRDCSSDescriptor(d.ptr()) {
		t.Fatal("error")
	}
}

func TestRDCSS(t *testing.T) {
	control := uint64(0)
	data := uint64(0)
	old := rdcss(&rdcssDescriptor{
		a1: &control,
		o1: 0,
		a2: &data,
		o2: 0,
		n2: 1,
	})
	if old != 0 && data != 1 {
		t.Fatal("RDCSS failed")
	}
}

func BenchmarkRDCSS(b *testing.B) {
	data := []uint64{0, 0}
	b.SetBytes(1)
	for n := uint64(0); n < uint64(b.N); n++ {
		rdcss(&rdcssDescriptor{
			a1: &data[0],
			o1: 0,
			a2: &data[1],
			o2: n + 0,
			n2: n + 1,
		})
	}
}

func BenchmarkRDCSSParallel(b *testing.B) {
	data := []uint64{0, 0}
	desc := make([]*rdcssDescriptor, b.N)
	for i := uint64(0); i < uint64(b.N); i++ {
		desc[i] = &rdcssDescriptor{
			a1: &data[0],
			o1: 0,
			a2: &data[1],
			o2: i + 0,
			n2: i + 1,
		}
	}
	b.SetBytes(1)
	b.RunParallel(func(pb *testing.PB) {
		for n := uint64(0); pb.Next(); n++ {
			rdcss(desc[n])
		}
	})
}

func TestCASN(t *testing.T) {
	data := []uint64{0, 1, 2, 3}
	if CASN([]Update{
		{&data[0], 0, 1},
		{&data[1], 1, 2},
		{&data[2], 2, 3},
		{&data[3], 3, 4},
	}) != true {
		t.Fatal("CASN should be successful")
	}
	if data[0] != 1 || data[1] != 2 || data[2] != 3 || data[3] != 4 {
		t.Fatal("CASN didn't swap values")
	}
	if CASN([]Update{
		{&data[0], 0, 1},
		{&data[1], 1, 2},
		{&data[2], 2, 3},
		{&data[3], 3, 4},
	}) != false {
		t.Fatal("CASN should have failed")
	}
	if data[0] != 1 || data[1] != 2 || data[2] != 3 || data[3] != 4 {
		t.Fatal("CASN shouldn't have swapped values")
	}
}

func BenchmarkMutex(b *testing.B) {
	data := []uint64{0, 1, 2, 3}
	mu := &sync.Mutex{}
	b.SetBytes(4)
	for n := uint64(0); n < uint64(b.N); n++ {
		mu.Lock()
		data[0] = n + 1
		data[1] = n + 2
		data[2] = n + 3
		data[3] = n + 4
		mu.Unlock()
	}
}

func BenchmarkMutexParallel(b *testing.B) {
	data := []uint64{0, 1, 2, 3}
	mu := &sync.Mutex{}
	b.SetBytes(4)
	b.RunParallel(func(pb *testing.PB) {
		for n := uint64(0); pb.Next(); n++ {
			mu.Lock()
			data[0] = n + 1
			data[1] = n + 2
			data[2] = n + 3
			data[3] = n + 4
			mu.Unlock()
		}
	})
}

func BenchmarkCASN(b *testing.B) {
	data := []uint64{0, 1, 2, 3}
	b.SetBytes(4)
	for n := uint64(0); n < uint64(b.N); n++ {
		CASN([]Update{
			{&data[0], n + 0, n + 1},
			{&data[1], n + 1, n + 2},
			{&data[2], n + 2, n + 3},
			{&data[3], n + 3, n + 4},
		})
	}
}

func BenchmarkCASNParallel(b *testing.B) {
	data := []uint64{0, 1, 2, 3}
	b.SetBytes(4)
	b.RunParallel(func(pb *testing.PB) {
		for n := uint64(0); pb.Next(); n++ {
			CASN([]Update{
				{&data[0], n + 0, n + 1},
				{&data[1], n + 1, n + 2},
				{&data[2], n + 2, n + 3},
				{&data[3], n + 3, n + 4},
			})
		}
	})
}

func TestCas(t *testing.T) {
	num := uint64(0)
	if old := cas(&num, 0, 1); old != 0 {
		t.Fatal("Cas didn't swap")
	}
	if old := cas(&num, 2, 0); old != 1 {
		t.Fatal("Cas shouldn't have swapped")
	}
}

func BenchmarkCAS(b *testing.B) {
	data := uint64(0)
	b.SetBytes(1)
	for n := uint64(0); n < uint64(b.N); n++ {
		atomic.CompareAndSwapUint64(&data, n, n+1)
	}
}

func BenchmarkCASParallel(b *testing.B) {
	data := uint64(0)
	b.SetBytes(1)
	b.RunParallel(func(pb *testing.PB) {
		for n := uint64(0); pb.Next(); n++ {
			atomic.CompareAndSwapUint64(&data, n, n+1)
		}
	})
}

// BenchmarkEvaluation runs a benchmark most similar to the one found in the
// original paper in order to get an idea of how closely we've followed the
// correct implementation.
//
// TODO
func BenchmarkEvaluation(b *testing.B) {
	b.ReportMetric(1.00, "testing")
}
