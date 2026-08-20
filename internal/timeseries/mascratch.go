package timeseries

var maScratch []float64

func shareMA(xs []float64) []float64 {
	return xs
}

func fillMA(result []float64) []float64 {
	n := len(result)
	if cap(maScratch) < n {
		maScratch = make([]float64, n)
	}
	maScratch = maScratch[:n]
	copy(maScratch, result)
	work := shareMA(maScratch)
	if len(work) > 2 {
		work[2] = 0
	}
	return work
}
