package solar

import (
	"math"

	"pv-forecast/internal/tracker"
)

func TiltedIrradiance(ghi, dni, dhi, tiltDeg, surfAzDeg, solarZenDeg, solarAzDeg float64) float64 {
	tiltDeg = tracker.HoldTilt(tiltDeg)
	tilt := tiltDeg * math.Pi / 180
	surfAz := surfAzDeg * math.Pi / 180
	zenith := solarZenDeg * math.Pi / 180
	solAz := solarAzDeg * math.Pi / 180

	cosAOI := math.Sin(zenith)*math.Sin(tilt)*math.Cos(solAz-surfAz) +
		math.Cos(zenith)*math.Cos(tilt)
	if cosAOI < 0 {
		cosAOI = 0
	}

	beam := dni * cosAOI

	diffuse := dhi * (1 + math.Cos(tilt)) / 2

	albedo := 0.2
	ground := ghi * albedo * (1 - math.Cos(tilt)) / 2

	total := beam + diffuse + ground
	if total < 0 {
		return 0
	}
	return total
}

func AOI(tiltDeg, surfAzDeg, solarZenDeg, solarAzDeg float64) float64 {
	tilt := tiltDeg * math.Pi / 180
	surfAz := surfAzDeg * math.Pi / 180
	zenith := solarZenDeg * math.Pi / 180
	solAz := solarAzDeg * math.Pi / 180

	cosAOI := math.Sin(zenith)*math.Sin(tilt)*math.Cos(solAz-surfAz) +
		math.Cos(zenith)*math.Cos(tilt)
	if cosAOI > 1 {
		cosAOI = 1
	}
	if cosAOI < -1 {
		cosAOI = -1
	}
	return math.Acos(cosAOI) * 180 / math.Pi
}

func IAMLoss(aoiDeg, b0 float64) float64 {
	if aoiDeg >= 90 {
		return 0
	}
	aoi := aoiDeg * math.Pi / 180
	iam := 1 - b0*(1/math.Cos(aoi)-1)
	if iam < 0 {
		return 0
	}
	return iam
}

func OptimalTilt(latDeg float64) float64 {
	return math.Abs(latDeg)*0.76 + 3.1
}

func Perez(dhi, dni, zenithDeg, tiltDeg, aoiDeg float64) float64 {
	if dhi <= 0 {
		return 0
	}
	tilt := tiltDeg * math.Pi / 180
	zenith := zenithDeg * math.Pi / 180

	a := math.Max(0, math.Cos(aoiDeg*math.Pi/180))
	b := math.Max(math.Cos(85*math.Pi/180), math.Cos(zenith))

	circumsolar := dhi * a / b * 0.3

	isotropic := dhi * (1 + math.Cos(tilt)) / 2 * 0.5

	horizon := dhi * math.Sin(tilt) * 0.2

	return circumsolar + isotropic + horizon
}

func TranspositionFactor(poa, ghi float64) float64 {
	if ghi <= 0 {
		return 0
	}
	return poa / ghi
}

func DiffuseIsotropic(dhi, tiltDeg float64) float64 {
	tilt := tiltDeg * math.Pi / 180
	return dhi * (1 + math.Cos(tilt)) / 2
}

func GroundReflected(ghi, tiltDeg, albedo float64) float64 {
	tilt := tiltDeg * math.Pi / 180
	return ghi * albedo * (1 - math.Cos(tilt)) / 2
}

func EffectiveIrradiance(poa, iam, spectralFactor float64) float64 {
	eff := poa * iam * spectralFactor
	if eff < 0 {
		return 0
	}
	return eff
}

func SpectralFactor(airMass float64) float64 {
	if airMass <= 0 {
		return 1
	}
	am := airMass
	return 0.918093 + 0.086257*am - 0.024459*am*am + 0.002816*am*am*am - 0.000126*am*am*am*am
}
