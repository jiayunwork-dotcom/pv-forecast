// Package tracker 实现单轴和双轴跟踪器角度优化。
package tracker

import "math"

// Mode 跟踪模式。
type Mode int

const (
	Fixed     Mode = iota // 固定倾角
	SingleH               // 单轴水平（南北轴）
	SingleV               // 单轴竖直
	DualAxis              // 双轴
)

// Config 跟踪器配置。
type Config struct {
	Mode         Mode
	MaxAngle     float64 // 最大旋转角（度）
	BacktrackGCR float64 // 地面覆盖率（用于反向跟踪）
	AxisTiltDeg  float64 // 轴倾角（度）
	AxisAzDeg    float64 // 轴方位（度）
}

// DefaultSingleAxis 返回默认单轴跟踪配置。
func DefaultSingleAxis() Config {
	return Config{
		Mode:         SingleH,
		MaxAngle:     60,
		BacktrackGCR: 0.35,
	}
}

// TrackAngle 计算跟踪器旋转角度（度）。
func TrackAngle(solarZenDeg, solarAzDeg float64, cfg Config) float64 {
	switch cfg.Mode {
	case Fixed:
		return 0
	case SingleH:
		return singleHorizontal(solarZenDeg, solarAzDeg, cfg)
	case SingleV:
		return singleVertical(solarZenDeg, solarAzDeg, cfg)
	case DualAxis:
		return 0 // 双轴直接面向太阳，下面用 DualAxisAngles
	}
	return 0
}

// DualAxisAngles 返回双轴跟踪的倾角和方位角。
func DualAxisAngles(solarZenDeg, solarAzDeg float64) (tiltDeg, azDeg float64) {
	return solarZenDeg, solarAzDeg
}

func singleHorizontal(zenDeg, azDeg float64, cfg Config) float64 {
	zen := zenDeg * math.Pi / 180
	az := azDeg * math.Pi / 180
	axisAz := cfg.AxisAzDeg * math.Pi / 180

	// 投影太阳矢量到垂直于轴的平面
	projAngle := math.Atan2(math.Sin(zen)*math.Sin(az-axisAz), math.Cos(zen))
	rotDeg := projAngle * 180 / math.Pi

	// 限幅
	if rotDeg > cfg.MaxAngle {
		rotDeg = cfg.MaxAngle
	}
	if rotDeg < -cfg.MaxAngle {
		rotDeg = -cfg.MaxAngle
	}

	// 反向跟踪
	if cfg.BacktrackGCR > 0 {
		rotDeg = backtrack(rotDeg, cfg.BacktrackGCR)
	}

	return rotDeg
}

func singleVertical(zenDeg, azDeg float64, cfg Config) float64 {
	// 竖直轴跟踪：方位角追踪
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

// AOI 计算跟踪面的入射角。
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

// ShadingLoss 估计行间遮挡损失。
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
	return shadeFrac * 0.5 // 部分遮挡损失折半
}

// EnergyGain 计算跟踪相对于固定安装的发电量增益（百分比估算）。
func EnergyGain(latDeg float64, mode Mode) float64 {
	absLat := math.Abs(latDeg)
	switch mode {
	case SingleH:
		return 15 + 0.1*(30-absLat) // 约 12-18%
	case DualAxis:
		return 25 + 0.15*(30-absLat) // 约 20-30%
	default:
		return 0
	}
}

// RowSpacing 计算避免遮挡所需的行间距。
func RowSpacing(moduleWidth, tiltDeg, minSolarElevDeg float64) float64 {
	if minSolarElevDeg <= 0 {
		return moduleWidth * 10 // 保守值
	}
	tilt := tiltDeg * math.Pi / 180
	elev := minSolarElevDeg * math.Pi / 180
	return moduleWidth * math.Sin(tilt) / math.Tan(elev)
}
