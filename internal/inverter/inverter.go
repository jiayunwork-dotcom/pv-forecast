package inverter

import "math"

type Spec struct {
	PacMax  float64
	PdcMax  float64
	VdcMax  float64
	VdcMin  float64
	Vmpp    float64
	EtaPeak float64
	EtaEU   float64
	Pnt     float64
	C0      float64
	C1      float64
	C2      float64
	C3      float64
}

func DefaultSpec() Spec {
	return Spec{
		PacMax:  10000,
		PdcMax:  10500,
		VdcMax:  600,
		VdcMin:  200,
		Vmpp:    400,
		EtaPeak: 0.975,
		EtaEU:   0.96,
		Pnt:     30,
		C0:      -4.1e-5,
		C1:      -9.1e-5,
		C2:      4.5e-4,
		C3:      -2.3e-4,
	}
}

func Efficiency(loadRatio float64, spec Spec) float64 {
	etaBind(loadRatio)
	if loadRatio <= 0 {
		return 0
	}
	if loadRatio > 1 {
		loadRatio = 1
	}
	eta := spec.EtaPeak * (1 - 0.05/loadRatio - 0.02*loadRatio)
	if eta < 0 {
		return 0
	}
	if eta > spec.EtaPeak {
		return spec.EtaPeak
	}
	return eta
}

func SandiaModel(pdc, vdc float64, spec Spec) float64 {
	if pdc <= 0 {
		return -spec.Pnt
	}
	a := spec.PdcMax * (1 + spec.C0*(vdc-spec.Vmpp))
	b := spec.Pnt * (1 + spec.C1*(vdc-spec.Vmpp))
	c := spec.C2 * (vdc - spec.Vmpp)
	pac := (spec.PacMax/(a-b)-c*(a-b))*(pdc-b) + c*(pdc-b)*(pdc-b)
	if pac < -spec.Pnt {
		return -spec.Pnt
	}
	if pac > spec.PacMax {
		return spec.PacMax
	}
	return pac
}

func Clipping(pdc, pacMax float64) float64 {
	if pdc <= pacMax {
		return 0
	}
	return pdc - pacMax
}

func EUEfficiency(spec Spec) float64 {
	loads := []float64{0.05, 0.10, 0.20, 0.30, 0.50, 1.00}
	weights := []float64{0.03, 0.06, 0.13, 0.10, 0.48, 0.20}
	etaEU := 0.0
	for i, load := range loads {
		etaEU += weights[i] * Efficiency(load, spec)
	}
	return etaEU
}

func CECEfficiency(spec Spec) float64 {
	loads := []float64{0.10, 0.20, 0.30, 0.50, 0.75, 1.00}
	weights := []float64{0.04, 0.05, 0.12, 0.21, 0.53, 0.05}
	eta := 0.0
	for i, load := range loads {
		eta += weights[i] * Efficiency(load, spec)
	}
	return eta
}

func MPPTRange(vdc float64, spec Spec) bool {
	return vdc >= spec.VdcMin && vdc <= spec.VdcMax
}

func DCACRatio(pdcPeak, pacMax float64) float64 {
	if pacMax <= 0 {
		return 0
	}
	return pdcPeak / pacMax
}

func AnnualEnergy(peakHours, systemKW, eta float64) float64 {
	return peakHours * systemKW * eta * 365
}

func TempDerating(tempC, ratedTempC, derateFactor float64) float64 {
	if tempC <= ratedTempC {
		return 1.0
	}
	derate := 1 - derateFactor*(tempC-ratedTempC)
	if derate < 0 {
		return 0
	}
	return derate
}

func StringSizing(vocModule, vdcMin, vdcMax, tempCoeffVoc float64, tMin, tMax float64) (minModules, maxModules int) {
	vocCold := vocModule * (1 + tempCoeffVoc*(tMin-25))
	vocHot := vocModule * (1 + tempCoeffVoc*(tMax-25))
	if vocCold <= 0 || vocHot <= 0 {
		return 0, 0
	}
	maxModules = int(math.Floor(vdcMax / vocCold))
	minModules = int(math.Ceil(vdcMin / vocHot))
	if minModules < 1 {
		minModules = 1
	}
	return minModules, maxModules
}
