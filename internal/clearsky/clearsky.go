package clearsky

import "math"

func Ineichen(zenithDeg, altitude, turbidity float64) (ghi, dni, dhi float64) {
	if zenithDeg >= 90 {
		return 0, 0, 0
	}
	cosZ := math.Cos(zenithDeg * math.Pi / 180)
	am := 1 / (cosZ + 0.50572*math.Pow(96.07995-zenithDeg, -1.6364))
	if am < 0 {
		am = 0
	}

	I0 := 1367.0

	fh1 := math.Exp(-altitude / 8000)
	fh2 := math.Exp(-altitude / 1250)

	cg1 := 5.09e-5*altitude + 0.868
	cg2 := 3.92e-5*altitude + 0.0387

	b := 0.664 + 0.163/fh1
	dniClear := b * I0 * math.Exp(-0.09*am*turbidity*fh1)
	if dniClear < 0 {
		dniClear = 0
	}

	ghiClear := cg1 * I0 * cosZ * math.Exp(-cg2*am*(fh1+fh2)) * math.Exp(0.01*am*am*turbidity)
	if ghiClear < 0 {
		ghiClear = 0
	}

	dhiClear := ghiClear - dniClear*cosZ
	if dhiClear < 0 {
		dhiClear = 0
	}

	return ghiClear, dniClear, dhiClear
}

func Haurwitz(zenithDeg float64) float64 {
	if zenithDeg >= 90 {
		return 0
	}
	cosZ := math.Cos(zenithDeg * math.Pi / 180)
	return 1098 * cosZ * math.Exp(-0.057/cosZ)
}

func SimplifiedSolis(zenithDeg, altitude, aod700 float64) (ghi, dni, dhi float64) {
	if zenithDeg >= 90 {
		return 0, 0, 0
	}
	cosZ := math.Cos(zenithDeg * math.Pi / 180)
	I0 := 1367.0 * cosZ

	altFactor := 1 + 0.0001*altitude
	tau := 0.7 - 0.7*aod700
	if tau < 0.1 {
		tau = 0.1
	}
	dniVal := I0 * tau * altFactor
	ghiVal := dniVal * cosZ * 1.1
	dhiVal := ghiVal - dniVal*cosZ

	if dniVal < 0 {
		dniVal = 0
	}
	if ghiVal < 0 {
		ghiVal = 0
	}
	if dhiVal < 0 {
		dhiVal = 0
	}
	return ghiVal, dniVal, dhiVal
}

func ClearSkyIndex(ghi, ghiClearSky float64) float64 {
	if ghiClearSky <= 0 {
		return 0
	}
	kt := ghi / ghiClearSky
	if kt > 1.5 {
		kt = 1.5
	}
	if kt < 0 {
		kt = 0
	}
	return kt
}

func CloudFraction(kt float64) float64 {
	if kt >= 1 {
		return 0
	}
	if kt <= 0 {
		return 1
	}
	return 1 - kt
}
