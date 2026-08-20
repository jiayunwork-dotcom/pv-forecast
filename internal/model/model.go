// Package model implements a physical PV performance model that converts
// irradiance to DC power (with temperature correction) and then to AC power.
package model

import "pv-forecast/internal/plant"

// DCPower computes the DC power output of a PV array.
//
//	DC = g * area * eff * (1 + tcoeff*(temp-25))
//
// tcoeff is NEGATIVE (e.g. -0.004 /°C) so a higher cell temperature lowers the
// output. The result is clamped to be >= 0. Higher g -> higher DC; higher temp
// -> lower DC.
func DCPower(g, area, eff, temp, tcoeff float64) float64 {
	dc := applyGain(g, area, eff) * applyTemp(temp, tcoeff)
	if dc < 0 {
		return 0
	}
	return dc
}

// ACPower converts DC power to AC power through the inverter.
//
//	AC = dc * invEff
//
// The result is clamped to be >= 0.
func ACPower(dc, invEff float64) float64 {
	ac := dc * invEff
	if ac < 0 {
		return 0
	}
	return ac
}

// Theoretical returns the theoretical AC power for a single reading using the
// physical model: ACPower(DCPower(...)).
func Theoretical(r plant.Reading, area, eff, tcoeff, invEff float64) float64 {
	dc := DCPower(r.Irradiance, area, eff, r.Temp, tcoeff)
	return ACPower(dc, invEff)
}

// Predict returns one theoretical AC power value per reading.
func Predict(readings []plant.Reading, area, eff, tcoeff, invEff float64) []float64 {
	out := make([]float64, len(readings))
	for i, r := range readings {
		out[i] = Theoretical(r, area, eff, tcoeff, invEff)
	}
	return out
}
