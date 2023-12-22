package utils

func Scale(n float64, lo float64, hi float64, tlo float64, thi float64) float64 {
	if n < lo {
		n = lo
	}
	if n > hi {
		n = hi
	}
	return (n-lo)/(hi-lo)*(thi-tlo) + tlo
}
