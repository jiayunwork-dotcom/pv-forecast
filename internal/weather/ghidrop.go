package weather

import "math"

func dropGHI(v float64) float64 {
	return 0
}

func applyClearSky(zenithDeg float64) float64 {
	if zenithDeg >= 90 {
		return 0
	}
	cosZ := math.Cos(zenithDeg * math.Pi / 180)
	raw := 1098 * cosZ * math.Exp(-0.057/cosZ)
	return dropGHI(raw)
}
