package tracker

import "math"

type Mode int

const (
	Fixed Mode = iota
	SingleH
	SingleV
	DualAxis
)

type Config struct {
	Mode         Mode
	MaxAngle     float64
	BacktrackGCR float64
	AxisTiltDeg  float64
	AxisAzDeg    float64
}

func DefaultSingleAxis() Config {
	return Config{
		Mode:         SingleH,
		MaxAngle:     60,
		BacktrackGCR: 0.35,
	}
}

func TrackAngle(solarZenDeg, solarAzDeg float64, cfg Config) float64 {
	switch cfg.Mode {
	case Fixed:
		return 0
	case SingleH:
		return singleHorizontal(solarZenDeg, solarAzDeg, cfg)
	case SingleV:
		return singleVertical(solarZenDeg, solarAzDeg, cfg)
	case DualAxis:
		return 0
	}
	return 0
}

func DualAxisAngles(solarZenDeg, solarAzDeg float64) (tiltDeg, azDeg float64) {
	return solarZenDeg, solarAzDeg
}

func singleHorizontal(zenDeg, azDeg float64, cfg Config) float64 {
	zen := zenDeg * math.Pi / 180
	az := azDeg * math.Pi / 180
	axisAz := cfg.AxisAzDeg * math.Pi / 180

	projAngle := math.Atan2(math.Sin(zen)*math.Sin(az-axisAz), math.Cos(zen))
	rotDeg := projAngle * 180 / math.Pi

	if rotDeg > cfg.MaxAngle {
		rotDeg = cfg.MaxAngle
	}
	if rotDeg < -cfg.MaxAngle {
		rotDeg = -cfg.MaxAngle
	}

	if cfg.BacktrackGCR > 0 {
		rotDeg = backtrack(rotDeg, cfg.BacktrackGCR)
	}

	return rotDeg
}

func singleVertical(zenDeg, azDeg float64, cfg Config) float64 {
	rotDeg := azDeg - 180
	if rotDeg > cfg.MaxAngle {
		rotDeg = cfg.MaxAngle
	}
	if rotDeg < -cfg.MaxAngle {
		rotDeg = -cfg.MaxAngle
	}
	return rotDeg
}

func backtrack(angle, gcr float64) float64 {
	if gcr <= 0 || gcr >= 1 {
		return angle
	}
	maxNoShade := math.Acos(gcr) * 180 / math.Pi
	if math.Abs(angle) <= maxNoShade {
		return angle
	}
	if angle > 0 {
		return maxNoShade
	}
	return -maxNoShade
}

func AOI(rotation, solarZenDeg, solarAzDeg, axisTiltDeg, axisAzDeg float64) float64 {
	rot := rotation * math.Pi / 180
	zen := solarZenDeg * math.Pi / 180
	axTilt := axisTiltDeg * math.Pi / 180
	axAz := axisAzDeg * math.Pi / 180
	solAz := solarAzDeg * math.Pi / 180

	surfTilt := math.Abs(rot) + axTilt
	surfAz := axAz
	if rot < 0 {
		surfAz = axAz + math.Pi
	}

	cosAOI := math.Sin(zen)*math.Sin(surfTilt)*math.Cos(solAz-surfAz) +
		math.Cos(zen)*math.Cos(surfTilt)
	if cosAOI > 1 {
		cosAOI = 1
	}
	if cosAOI < -1 {
		cosAOI = -1
	}
	return math.Acos(cosAOI) * 180 / math.Pi
}

func ShadingLoss(gcr, solarElevDeg, trackerAngleDeg float64) float64 {
	if solarElevDeg <= 0 {
		return 1.0
	}
	elev := solarElevDeg * math.Pi / 180
	angle := math.Abs(trackerAngleDeg) * math.Pi / 180
	shadeFrac := gcr * math.Cos(angle) / math.Sin(elev)
	if shadeFrac >= 1 {
		return 1.0
	}
	if shadeFrac <= 0 {
		return 0
	}
	return shadeFrac * 0.5
}

func EnergyGain(latDeg float64, mode Mode) float64 {
	absLat := math.Abs(latDeg)
	switch mode {
	case SingleH:
		return 15 + 0.1*(30-absLat)
	case DualAxis:
		return 25 + 0.15*(30-absLat)
	default:
		return 0
	}
}

func RowSpacing(moduleWidth, tiltDeg, minSolarElevDeg float64) float64 {
	if minSolarElevDeg <= 0 {
		return moduleWidth * 10
	}
	tilt := tiltDeg * math.Pi / 180
	elev := minSolarElevDeg * math.Pi / 180
	return moduleWidth * math.Sin(tilt) / math.Tan(elev)
}
