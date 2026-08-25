package sizing

import "math"

type SystemSpec struct {
	ModuleWp       float64
	ModuleVoc      float64
	ModuleIsc      float64
	ModuleVmpp     float64
	ModuleImpp     float64
	InverterPacMax float64
	InverterVdcMin float64
	InverterVdcMax float64
	InverterIdcMax float64
	TempCoeffVoc   float64
	TempCoeffIsc   float64
	TempMin        float64
	TempMax        float64
}

func DefaultSpec() SystemSpec {
	return SystemSpec{
		ModuleWp:       400,
		ModuleVoc:      48.5,
		ModuleIsc:      10.4,
		ModuleVmpp:     40.5,
		ModuleImpp:     9.9,
		InverterPacMax: 10000,
		InverterVdcMin: 200,
		InverterVdcMax: 600,
		InverterIdcMax: 30,
		TempCoeffVoc:   -0.29,
		TempCoeffIsc:   0.05,
		TempMin:        -10,
		TempMax:        50,
	}
}

func StringSize(spec SystemSpec) (min, max int) {
	vocCold := spec.ModuleVoc * (1 + spec.TempCoeffVoc/100*(spec.TempMin-25))
	vmppHot := spec.ModuleVmpp * (1 + spec.TempCoeffVoc/100*(spec.TempMax-25))
	max = int(math.Floor(spec.InverterVdcMax / vocCold))
	min = int(math.Ceil(spec.InverterVdcMin / vmppHot))
	if min < 1 {
		min = 1
	}
	return min, max
}

func ParallelStrings(spec SystemSpec) int {
	iscHot := spec.ModuleIsc * (1 + spec.TempCoeffIsc/100*(spec.TempMax-25))
	if iscHot <= 0 {
		return 0
	}
	return int(math.Floor(spec.InverterIdcMax / iscHot))
}

func DCACRatio(modulesPerString, strings int, spec SystemSpec) float64 {
	dcPeak := float64(modulesPerString*strings) * spec.ModuleWp
	if spec.InverterPacMax <= 0 {
		return 0
	}
	return dcPeak / spec.InverterPacMax
}

func TotalCapacityKW(modulesPerString, strings int, spec SystemSpec) float64 {
	return float64(modulesPerString*strings) * spec.ModuleWp / 1000
}

func AreaRequired(modulesPerString, strings int, moduleAreaM2, gcr float64) float64 {
	totalModules := float64(modulesPerString * strings)
	moduleArea := totalModules * moduleAreaM2
	if gcr <= 0 {
		gcr = 0.4
	}
	return moduleArea / gcr
}

func AnnualYield(capacityKWp, peakSunHours, systemLossFactor float64) float64 {
	return capacityKWp * peakSunHours * 365 * systemLossFactor
}

func LCOE(totalCost, annualGenKWh float64, years int, discountRate float64) float64 {
	if annualGenKWh <= 0 || years <= 0 {
		return 0
	}
	totalGen := 0.0
	for y := 0; y < years; y++ {
		totalGen += annualGenKWh / math.Pow(1+discountRate, float64(y))
	}
	return totalCost / totalGen
}

func PaybackYears(totalCost, annualSavings float64) float64 {
	if annualSavings <= 0 {
		return math.Inf(1)
	}
	return totalCost / annualSavings
}

func CableSize(current, lengthM, maxDropPct, voltage float64) float64 {
	if voltage <= 0 || maxDropPct <= 0 {
		return 0
	}
	resistivity := 0.0175
	maxDrop := voltage * maxDropPct / 100
	area := 2 * resistivity * lengthM * current / maxDrop
	return area
}

func ShadeFreeSpacing(moduleHeight, tiltDeg, minElevDeg float64) float64 {
	if minElevDeg <= 0 {
		return moduleHeight * 5
	}
	tilt := tiltDeg * math.Pi / 180
	elev := minElevDeg * math.Pi / 180
	return moduleHeight * math.Sin(tilt) / math.Tan(elev)
}
