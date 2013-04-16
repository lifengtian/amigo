package myutil

import (
	"testing"
)

func TestAdd1 (t *testing.T ) {
	a := 2
	b := Add1(a )

	if b != 3 {
		t.Errorf("Add1(2) failed. Got %d, expected 3", b)
	}
}

func BenchmarkAdd1 (b *testing.B ) {
	for i := 0 ; i < b.N ; i++ {
		Add1(3)
	}
}

