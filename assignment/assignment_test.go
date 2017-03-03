package assignment

import "testing"

//                 A, B, C, D, E, F, G, H, I, J, K, L
var cardin = []int{2, 3, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2}
var testNew = []struct {
	varlist []int
	values  []int
}{
	{[]int{}, []int{}},
	{[]int{1}, []int{0}},
	{[]int{0, 1, 2}, []int{0, 0, 0}},
	{[]int{3, 7, 10}, []int{0, 0, 0}},
	{[]int{6, 2, 11, 9}, []int{0, 0, 0, 0}},
}

func TestNew(t *testing.T) {
	for _, v := range testNew {
		assig := New(v.varlist, cardin)
		for i := range v.varlist {
			if assig.Variables()[i] != v.varlist[i] {
				t.Errorf("Missing variables on assignment: %v", v.varlist[i])
			}
			if assig.Values()[i] != v.values[i] {
				t.Errorf("Initialized with wrong value: %v", v.values[i])
			}
		}
	}
}

var testNext = []struct {
	varlist []int
	next    int
	values  []int
}{
	{[]int{0}, 1, []int{1}},
	{[]int{1}, 2, []int{2}},
	{[]int{0, 1, 2}, 5, []int{1, 2, 0}},
	{[]int{0, 1, 2}, 0, []int{0, 0, 0}},
	{[]int{0, 1, 2}, 1, []int{1, 0, 0}},
	{[]int{3, 7, 10}, 7, []int{1, 1, 1}},
	{[]int{6, 2, 11, 9}, 14, []int{0, 1, 1, 1}},
	{[]int{}, 1, []int{}},
}

func TestNext(t *testing.T) {
	for _, v := range testNext {
		assig := New(v.varlist, cardin)
		for i := 0; i < v.next; i++ {
			assig.Next()
		}
		for i := range v.varlist {
			if assig.Var(i) != v.varlist[i] {
				t.Errorf("Missing variables on assignment: %v", v.varlist[i])
			}
			if assig.Value(i) != v.values[i] {
				t.Errorf("Initialized with wrong value. want %v, got %v", v.values, assig)
			}
		}
	}
	assig := New([]int{0}, []int{2})
	assig.Next()
	hasnext := assig.Next()
	if hasnext {
		t.Errorf("Want end of assig, got %v", assig)
	}

}

var testConsistent = []struct {
	varlist   []int
	next      int
	values    []int
	consist   []int
	inconsist []int
}{
	{[]int{0}, 1, []int{1}, []int{}, []int{0}},
	{[]int{1}, 2, []int{2}, []int{1, 2}, []int{1, 1}},
	{[]int{0, 1, 2}, 5, []int{1, 2, 0}, []int{1, 2, 0}, []int{1, 1}},
	{[]int{0, 1, 2}, 0, []int{0, 0, 0}, []int{-1, 0, 0, 0}, []int{0, 0, 1, 1, 1}},
	{[]int{3, 7, 10}, 7, []int{1, 1, 1}, []int{0, 2, 0, 1, 0, 0, 0, 1, 0, 0, 1}, []int{0, 2, 0, 1, 0, 0, 0, 0}},
	{[]int{6, 2, 11, 9}, 14, []int{0, 1, 1, 1}, []int{1, 1, 1, 1, 1, 1, 0}, []int{1, 1, 1, 1, 1, 0, 1}},
}

func TestConsistent(t *testing.T) {
	for _, v := range testConsistent {
		assig := New(v.varlist, cardin)
		for i := 0; i < v.next; i++ {
			assig.Next()
		}
		if !assig.Consistent(v.consist) {
			t.Errorf("Assignment should be consistent. assig %v, val %v", assig, v.consist)
		}
		if assig.Consistent(v.inconsist) {
			t.Errorf("Assignment shouldn't be consistent. assig %v, val %v", assig, v.inconsist)
		}
	}
}

var testIndex = []struct {
	varlist []int
	cardin  []int
	values  []int
	stride  map[int]int
	result  int
}{
	{
		[]int{0, 1},
		[]int{2, 2},
		[]int{0, 1},
		map[int]int{0: 1, 1: 2},
		2,
	},
	{
		[]int{0, 1, 2},
		[]int{2, 3, 2},
		[]int{0, 2, 1},
		map[int]int{0: 1, 1: 2},
		4,
	},
	{
		[]int{0, 1, 2},
		[]int{2, 3, 2},
		[]int{0, 2, 1},
		map[int]int{1: 2, 2: 6},
		10,
	},
	{
		[]int{},
		[]int{2, 2, 2},
		[]int{},
		map[int]int{1: 2, 2: 6},
		0,
	},
}

func TestIndex(t *testing.T) {
	for _, v := range testIndex {
		a := New(v.varlist, v.cardin)
		for i, k := range v.values {
			a.values[i] = k
		}
		got := a.Index(v.stride)
		if got != v.result {
			t.Errorf("want %v, got %v", v.result, got)
		}
	}
}
