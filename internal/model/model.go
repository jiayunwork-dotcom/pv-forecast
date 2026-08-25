package model

import (
	"pv-forecast/internal/celltemp"
	"pv-forecast/internal/plant"
)

func DCPower(g, area, eff, temp, tcoeff float64) float64 {
	g = celltemp.HoldGHI(g)
	dc := g * area * eff * (1 + tcoeff*(temp-25))
	if dc < 0 {
		return 0
	}
	return dc
}

func ACPower(dc, invEff float64) float64 {
	ac := dc * invEff
	if ac < 0 {
		return 0
	}
	return ac
}

func Theoretical(r plant.Reading, area, eff, tcoeff, invEff float64) float64 {
	dc := DCPower(r.Irradiance, area, eff, r.Temp, tcoeff)
	return ACPower(dc, invEff)
}

func Predict(readings []plant.Reading, area, eff, tcoeff, invEff float64) []float64 {
	out := make([]float64, len(readings))
	for i, r := range readings {
		out[i] = Theoretical(r, area, eff, tcoeff, invEff)
	}
	return out
}
