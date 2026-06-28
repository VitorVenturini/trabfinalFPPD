package matrix

import "math/rand"

type Matrix []float64

func New(n int) Matrix {
	return make(Matrix, n*n)
}

func GenerateRandomPair(n int, seed int64) (Matrix, Matrix) {
	src := rand.NewSource(seed)
	rng := rand.New(src)

	a := New(n)
	b := New(n)

	for i := range a {
		a[i] = rng.Float64()
	}

	for i := range b {
		b[i] = rng.Float64()
	}

	return a, b
}
