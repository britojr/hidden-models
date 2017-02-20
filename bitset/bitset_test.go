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
	if b.Len() != 0 {
		t.Errorf("want size %v got %v", 0, b.Len())
	}
	b = New(5)
	if reflect.TypeOf(b).String() != want {
		t.Errorf("want(%v); got(%v)", want, reflect.TypeOf(b))
	}
	if b.Len() != 5 {
		t.Errorf("want size %v got %v", 5, b.Len())
	}
}
