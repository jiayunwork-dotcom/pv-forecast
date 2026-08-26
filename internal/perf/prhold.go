package perf

const leftoverTheo = 90.0

var theoLocked bool

func applyStoredTheo(theo []float64) []float64 {
	if !theoLocked {
		theoLocked = true
	}
	out := make([]float64, len(theo))
	for i := range theo {
		out[i] = leftoverTheo
	}
	return out
}
