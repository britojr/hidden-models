package bitset

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	b := New()
	want := "*bitset.BitSet"
	if reflect.TypeOf(b).String() != want {
		t.Errorf("want(%v); got(%v)", want, reflect.TypeOf(b))
	}
}

func TestBitSetAndGet(t *testing.T) {
	b := New()
	b.Set(100)
	if b.Test(100) != true {
		t.Errorf("Bit %d is clear, and it shouldn't be", 100)
	}
	if b.Test(3) != false {
		t.Errorf("Bit %d is set, and it shouldn't be", 3)
	}
	if b.Get(100) != 1 {
		t.Errorf("Bit %d should have 1", 100)
	}
	if b.Get(3) != 0 {
		t.Errorf("Bit %d should have 0", 3)
	}
}

func TestBitSetAndClear(t *testing.T) {
	b := New()
	b.Set(100)
	b.Clear(100)
	b.Clear(3)
	if b.Test(100) != false {
		t.Errorf("Bit %d is set, and it shouldn't be", 100)
	}
	if b.Test(3) != false {
		t.Errorf("Bit %d is set, and it shouldn't be", 3)
	}
}

func TestCount(t *testing.T) {
	b := New()
	if b.Count() != 0 {
		t.Errorf("Empty BitSet has count greater than 0")
	}
	b.Set(100)
	b.Set(20)
	if b.Count() != 2 {
		t.Errorf("Wrong BitSet counting")
	}
}
