package assignment

import "testing"

//                 A, B, C, D, E, F, G, H, I, J, K, L
var cardin = []int{2, 3, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2}

func TestNew(t *testing.T) {
	cases := []struct {
		cardin  []int
		varlist []int
		values  []int
		next    bool
	}{
		{cardin, []int{}, []int{}, false},
		{[]int{}, []int{0}, []int{}, false},
		{[]int{0}, []int{0}, []int{}, false},
		{[]int{1}, []int{0}, []int{0}, true},
		{cardin, []int{1}, []int{0}, true},
		{cardin, []int{0, 1, 2}, []int{0, 0, 0}, true},
		{cardin, []int{3, 7, 10}, []int{0, 0, 0}, true},
		{cardin, []int{6, 2, 11, 9}, []int{0, 0, 0, 0}, true},
	}
	for _, tt := range cases {
		assig := New(tt.varlist, tt.cardin)
		for i := range tt.varlist {
			if assig.Variables()[i] != tt.varlist[i] {
				t.Errorf("Missing variables on assignment: %v", tt.varlist[i])
			}
		}
		next := assig.Next()
		if tt.next != next {
			t.Errorf("Wrong response in next, want %v, got %v", tt.next, next)
		}
		if tt.next {
			for i := range tt.varlist {
				if assig.Values()[i] != tt.values[i] {
					t.Errorf("Initialized with wrong value: %v", tt.values[i])
				}
			}
		}
	}
}

func TestNext(t *testing.T) {
	cases := []struct {
		cardin  []int
		varlist []int
		next    int
		values  []int
	}{
		{cardin, []int{0}, 2, []int{1}},
		{cardin, []int{1}, 3, []int{2}},
		{cardin, []int{0, 1, 2}, 6, []int{1, 2, 0}},
		{cardin, []int{0, 1, 2}, 1, []int{0, 0, 0}},
		{cardin, []int{0, 1, 2}, 2, []int{1, 0, 0}},
		{cardin, []int{3, 7, 10}, 8, []int{1, 1, 1}},
		{cardin, []int{6, 2, 11, 9}, 15, []int{0, 1, 1, 1}},
	}
	for _, v := range cases {
		assig := New(v.varlist, cardin)
		for i := 0; i < v.next; i++ {
			if !assig.Next() {
				t.Errorf("Should have next value: %v", assig)
			}
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
	assig.Next()
	hasnext := assig.Next()
	if hasnext {
		t.Errorf("Want end of assig, got %v", assig)
	}

	assig = New([]int{0}, []int{1})
	assig.Next()
	hasnext = assig.Next()
	if hasnext {
		t.Errorf("Want end of assig, got %v", assig)
	}

	assig = New([]int{0}, []int{0})
	hasnext = assig.Next()
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
	{[]int{0}, 2, []int{1}, []int{}, []int{0}},
	{[]int{1}, 3, []int{2}, []int{1, 2}, []int{1, 1}},
	{[]int{0, 1, 2}, 6, []int{1, 2, 0}, []int{1, 2, 0}, []int{1, 1}},
	{[]int{0, 1, 2}, 1, []int{0, 0, 0}, []int{-1, 0, 0, 0}, []int{0, 0, 1, 1, 1}},
	{[]int{3, 7, 10}, 8, []int{1, 1, 1}, []int{0, 2, 0, 1, 0, 0, 0, 1, 0, 0, 1}, []int{0, 2, 0, 1, 0, 0, 0, 0}},
	{[]int{6, 2, 11, 9}, 15, []int{0, 1, 1, 1}, []int{1, 1, 1, 1, 1, 1, 0}, []int{1, 1, 1, 1, 1, 0, 1}},
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
}{{
	[]int{0, 1},
	[]int{2, 2},
	[]int{0, 1},
	map[int]int{0: 1, 1: 2},
	2,
}, {
	[]int{0, 1, 2},
	[]int{2, 3, 2},
	[]int{0, 2, 1},
	map[int]int{0: 1, 1: 2},
	4,
}, {
	[]int{0, 1, 2},
	[]int{2, 3, 2},
	[]int{0, 2, 1},
	map[int]int{1: 2, 2: 6},
	10,
}, {
	[]int{},
	[]int{2, 2, 2},
	[]int{},
	map[int]int{1: 2, 2: 6},
	0,
}}

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
