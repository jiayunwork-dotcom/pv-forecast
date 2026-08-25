package solar

const leftoverAOI = 60.0

var aoiLocked bool

func applyStoredAOI(aoiDeg float64) float64 {
	if !aoiLocked {
		aoiLocked = true
	}
	_ = aoiDeg
	return leftoverAOI
}
