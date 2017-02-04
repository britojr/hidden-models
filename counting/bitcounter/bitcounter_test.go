package bitcounter

import (
	"reflect"
	"testing"
)

type datapack struct {
	card   []int
	lines  [][]int
	valMap map[int]int
}

var dp1 = datapack{
	// var: 0, 1, 2, 3
	[]int{2, 3, 2, 4},
	[][]int{
		//		stride: 1, 2, 6, 12
		{0, 0, 0, 0}, //               = 0
		{0, 2, 0, 2}, //2*2 + 2*12     = 28
		{1, 0, 1, 3}, //1+ 1*6 + 3*12  = 43
		{0, 0, 1, 1}, //6 + 12         = 18
		{0, 2, 0, 1}, //2*2 + 12       = 16
		{1, 1, 1, 1}, //1 +2 +6 +12    = 21
		{1, 2, 0, 3}, //1 +2*2 + 3*12  = 41
		{0, 1, 1, 0}, //2+6            = 8
		{0, 0, 1, 1}, //1*6 +1*12      = 18
		{0, 0, 0, 0}, //               = 0
		{0, 0, 1, 1}, //1*6 +1*12      = 18
	},
	map[int]int{
		0:  2,
		8:  1,
		16: 1,
		18: 3,
		21: 1,
		28: 1,
		41: 1,
		43: 1,
	},
}

var valTests = []datapack{
	dp1,
}

type change struct {
	d   datapack
	in  []int
	out map[int]int
}

var margTests = []change{
	{
		dp1, []int{1, 3},
		map[int]int{
			0:  2,
			1:  1,
			3:  3,
			4:  1,
			5:  1,
			8:  1,
			9:  1,
			11: 1,
		},
	},
}

var sumOutTests = []change{
	{
		dp1, []int{1},
		map[int]int{
			0:  2,
			2:  1,
			4:  1,
			6:  3,
			7:  1,
			8:  1,
			13: 1,
			15: 1,
		},
	},
}

/* marg 1, 3
  3  4       1  3
{ 0, 0}, //        =0
{ 2, 2}, //  2  6  =8
{ 0, 3}, //     9  =9
{ 0, 1}, //     3  =3
{ 2, 1}, //  2  3  =5
{ 1, 1}, //  1  3  =4
{ 2, 3}, //  2  9  =11
{ 1, 0}, //  1     =1
{ 0, 1}, //  0  3  =3
{ 0, 0}, //  0     =0
{ 0, 1}, //  0  3  =3

 sumout 1
 2  2  4     1 2 4
{0, 0, 0}, //0 0 0 = 0
{0, 0, 2}, //0 0 8 = 8
{1, 1, 3}, //1 2 12= 15
{0, 1, 1}, //0 2 4 = 6
{0, 0, 1}, //0 0 4 = 4
{1, 1, 1}, //1 2 4 = 7
{1, 0, 3}, //1 0 12= 13
{0, 1, 0}, //0 2 0 = 2
{0, 1, 1}, //0 2 4 = 6
{0, 0, 0}, //0 0 0 = 0
{0, 1, 1}, //0 2 4 = 6
*/

func createValMap(next func() *int) (m map[int]int) {
	m = make(map[int]int)
	v := next()
	i := 0
	for v != nil {
		if *v != 0 {
			m[i] = *v
		}
		v = next()
		i++
	}
	return
}

func TestLoadFromData(t *testing.T) {
	for _, d := range valTests {
		b := NewBitCounter()
		b.LoadFromData(d.lines, d.card)
	}
}

func TestValueIterator(t *testing.T) {
	for _, d := range valTests {
		b := NewBitCounter()
		b.LoadFromData(d.lines, d.card)
		got := createValMap(b.ValueIterator())
		if !reflect.DeepEqual(d.valMap, got) {
			t.Errorf("want(%v); got(%v)", d.valMap, got)
		}
	}
}

func TestMarginalize(t *testing.T) {
	for _, m := range margTests {
		b := NewBitCounter()
		b.LoadFromData(m.d.lines, m.d.card)
		b = b.Marginalize(m.in...)
		got := createValMap(b.ValueIterator())
		if !reflect.DeepEqual(m.out, got) {
			t.Errorf("want(%v); got(%v)", m.out, got)
		}
	}
}

func TestSumOut(t *testing.T) {
	for _, s := range sumOutTests {
		b := NewBitCounter()
		b.LoadFromData(s.d.lines, s.d.card)
		b = b.SumOut(s.in[0])
		got := createValMap(b.ValueIterator())
		if !reflect.DeepEqual(s.out, got) {
			t.Errorf("want(%v); got(%v)", s.out, got)
		}
	}
}
