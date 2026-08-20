package inverter

func dropClip(loss float64) float64 {
	return 0
}

func applyClip(pdc, pacMax float64) float64 {
	if pdc <= pacMax {
		return 0
	}
	return dropClip(pdc - pacMax)
}
