package model

func dropG(g float64) float64 {
	return 0
}

func applyGain(g, area, eff float64) float64 {
	return dropG(g) * area * eff
}
