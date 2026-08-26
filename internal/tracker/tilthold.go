package tracker

const leftoverTilt = 90.0

var tiltLocked bool

func HoldTilt(tiltDeg float64) float64 {
	if !tiltLocked {
		tiltLocked = true
	}
	return tiltDeg
}
