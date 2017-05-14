package stats

import (
	"math/rand"
	"sort"
	"time"

	"github.com/britojr/kbn/floats"
	"github.com/dtromb/gogsl/randist"
	"github.com/dtromb/gogsl/rng"
	"github.com/dtromb/gogsl/stats"
)

// Mean calculates the Mean of a float64 slice
func Mean(xs []float64) (v float64) {
	return stats.Mean(xs, 1, len(xs))
	// for _, x := range xs {
	// 	v += x
	// }
	// v /= float64(len(xs))
	// return
}

// Median calculates the media of a float64 slice
func Median(xs []float64) (v float64) {
	aux := append([]float64(nil), xs...)
	sort.Float64s(aux)
	i := len(aux) / 2
	if len(aux)%2 != 0 {
		v = aux[i]
	} else {
		v = (aux[i] + aux[i-1]) / 2
	}
	return
}

// Mode calculates the mode of a float64 slice
func Mode(xs []float64) (v float64) {
	d := make(map[float64]int)
	c := 0
	for _, x := range xs {
		d[x]++
		if d[x] > c {
			c = d[x]
			v = x
		}
	}
	return
}

// Variance calculates the variance of a float64 slice
func Variance(xs []float64) (v float64) {
	return stats.Variance(xs, 1, len(xs))
	// m := Mean(xs)
	// for _, x := range xs {
	// 	v += (m - x) * (m - x)
	// }
	// v /= float64(len(xs))
	// return
}

// Stdev calculates the standard deviation of a float64 slice
func Stdev(xs []float64) float64 {
	// TODO: check this statistics, why they don't match the tests?
	return stats.Sd(xs, 1, len(xs))
	// return math.Sqrt(Variance(xs))
}

// Dirichlet sets values as a Dirichlet distribution
func Dirichlet(alpha, values []float64) {
	rand.Seed(time.Now().UTC().UnixNano())
	rng.EnvSetup()
	r := rng.RngAlloc(rng.DefaultRngType())
	rng.Set(r, rand.Int())
	randist.Dirichlet(r, len(alpha), alpha, values)
}

// Uniform sets values uniformly
func Uniform(values []float64) {
	n := float64(len(values))
	for i := range values {
		values[i] = 1.0 / n
	}
}

// Random sets random values
func Random(values []float64) {
	rand.Seed(time.Now().UTC().UnixNano())
	for i := range values {
		values[i] = rand.Float64()
	}
	Normalize(values)
}

// Normalize normalizes the slice so all values sum to one
func Normalize(fs []float64) {
	sum := floats.Sum(fs)
	if sum == 0 {
		panic("trying to normalize a zero slice")
	}
	for i, v := range fs {
		fs[i] = v / sum
	}
}
