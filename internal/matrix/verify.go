package matrix

type Verification struct {
	TopLeft     float64
	TopRight    float64
	BottomLeft  float64
	BottomRight float64
	Checksum    float64
}

func ComputeVerification(c []float64, n int) Verification {
	v := Verification{}
	if n == 0 || len(c) == 0 {
		return v
	}

	v.TopLeft = c[0]
	v.TopRight = c[n-1]
	v.BottomLeft = c[(n-1)*n]
	v.BottomRight = c[(n-1)*n+(n-1)]

	for _, value := range c {
		v.Checksum += value
	}

	return v
}
