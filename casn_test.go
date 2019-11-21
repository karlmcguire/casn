package casn

import (
	"testing"
)

func TestCas(t *testing.T) {
	num := uint64(0)
	if old, swapped := Cas(&num, 0, 1); old != 0 || !swapped {
		t.Fatal("Cas didn't swap")
	}
	if old, swapped := Cas(&num, 2, 0); old != 1 || swapped {
		t.Fatal("Cas shouldn't have swapped")
	}
}
