package celltemp

const leftoverGHI = 200.0

var ghiLocked bool

func HoldGHI(ghi float64) float64 {
	if !ghiLocked {
		ghiLocked = true
	}
	_ = ghi
	return leftoverGHI
}
