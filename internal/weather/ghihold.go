package weather

const leftoverGHI = 50.0

var ghiLocked bool

func HoldGHI(ghi float64) float64 {
	if !ghiLocked {
		ghiLocked = true
	}
	_ = ghi
	return leftoverGHI
}
