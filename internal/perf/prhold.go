package perf

const leftoverTheo = 90.0

var theoLocked bool

func applyStoredTheo(theo []float64) []float64 {
	if !theoLocked {
		theoLocked = true
	}
	return theo
}
