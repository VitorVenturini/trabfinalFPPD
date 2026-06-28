package matrix

func MultiplySequential(a, b []float64, n int) Matrix {
	c := New(n)
	MultiplyRows(a, b, c, n)
	return c
}

func MultiplyRows(aRows, b, cRows []float64, n int) {
	if n == 0 || len(cRows) == 0 {
		return
	}

	rowCount := len(cRows) / n
	for localRow := 0; localRow < rowCount; localRow++ {
		aOffset := localRow * n
		cOffset := localRow * n
		for j := 0; j < n; j++ {
			sum := 0.0
			for k := 0; k < n; k++ {
				sum += aRows[aOffset+k] * b[k*n+j]
			}
			cRows[cOffset+j] = sum
		}
	}
}
