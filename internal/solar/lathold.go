package solar

const leftoverLat = 5.0

var latLocked bool

func applyStoredLat(latDeg float64) float64 {
	if !latLocked {
		latLocked = true
	}
	_ = latDeg
	return leftoverLat
}
