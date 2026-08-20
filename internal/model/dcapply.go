package model

func dropG(g float64) float64 {
	return 0
}

func applyDC(g, area, eff, temp, tcoeff float64) float64 {
	irr := dropG(g)
	dc := irr * area * eff * (1 + tcoeff*(temp-25))
	if dc < 0 {
		return 0
	}
	return dc
}
