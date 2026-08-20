// Package weather 提供气象数据处理与太阳辐照分解工具。
package weather

import (
	"math"
)

// SolarPos 太阳位置参数。
type SolarPos struct {
	Zenith  float64 // 天顶角（弧度）
	Azimuth float64 // 方位角（弧度）
	Elev    float64 // 高度角（弧度）
}

// DayOfYear 计算一年中的第几天（1-365）。
func DayOfYear(month, day int) int {
	daysInMonth := []int{0, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	doy := 0
	for m := 1; m < month && m <= 12; m++ {
		doy += daysInMonth[m]
	}
	doy += day
	return doy
}

// Declination 计算太阳赤纬角（弧度）。
func Declination(dayOfYear int) float64 {
	return 23.45 * math.Pi / 180 * math.Sin(2*math.Pi*(284+float64(dayOfYear))/365)
}

// HourAngle 计算时角（弧度），solarHour 为太阳时（0-24）。
func HourAngle(solarHour float64) float64 {
	return (solarHour - 12) * 15 * math.Pi / 180
}

// SolarPosition 计算太阳位置。
func SolarPosition(latDeg float64, dayOfYear int, solarHour float64) SolarPos {
	lat := latDeg * math.Pi / 180
	dec := Declination(dayOfYear)
	ha := HourAngle(solarHour)

	sinElev := math.Sin(lat)*math.Sin(dec) + math.Cos(lat)*math.Cos(dec)*math.Cos(ha)
	if sinElev > 1 {
		sinElev = 1
	}
	if sinElev < -1 {
		sinElev = -1
	}
	elev := math.Asin(sinElev)
	zenith := math.Pi/2 - elev

	cosAz := (math.Sin(dec) - math.Sin(lat)*sinElev) / (math.Cos(lat) * math.Cos(elev) + 1e-15)
	if cosAz > 1 {
		cosAz = 1
	}
	if cosAz < -1 {
		cosAz = -1
	}
	azimuth := math.Acos(cosAz)
	if ha > 0 {
		azimuth = 2*math.Pi - azimuth
	}

	return SolarPos{Zenith: zenith, Azimuth: azimuth, Elev: elev}
}

// AirMass 计算大气质量（Kasten-Young 模型）。
func AirMass(zenithDeg float64) float64 {
	if zenithDeg >= 90 {
		return 40 // 地平线极端值
	}
	z := zenithDeg * math.Pi / 180
	return 1 / (math.Cos(z) + 0.50572*math.Pow(96.07995-zenithDeg, -1.6364))
}

// ClearSkyGHI 计算晴天全局水平辐照度（Haurwitz 模型，W/m²）。
func ClearSkyGHI(zenithDeg float64) float64 {
	if zenithDeg >= 90 {
		return 0
	}
	cosZ := math.Cos(zenithDeg * math.Pi / 180)
	return 1098 * cosZ * math.Exp(-0.057/cosZ)
}

// ClearSkyDNI 估算晴天法向直射辐照度（简化模型）。
func ClearSkyDNI(zenithDeg, altitude float64) float64 {
	if zenithDeg >= 90 {
		return 0
	}
	am := AirMass(zenithDeg)
	altFactor := 1 + 0.0001*altitude
	return 1367 * 0.7 * math.Pow(am, -0.678) * altFactor
}

// DecomposeGHI 将 GHI 分解为 DNI 和 DHI（Erbs 模型）。
// kt 为晴空指数 = GHI / extraterrestrial
func DecomposeGHI(ghi, extraterrestrial float64) (dni, dhi float64) {
	if extraterrestrial <= 0 || ghi <= 0 {
		return 0, 0
	}
	kt := ghi / extraterrestrial
	if kt > 1 {
		kt = 1
	}

	// Erbs 相关系数
	var kd float64
	switch {
	case kt <= 0.22:
		kd = 1 - 0.09*kt
	case kt <= 0.80:
		kd = 0.9511 - 0.1604*kt + 4.388*kt*kt -
			16.638*kt*kt*kt + 12.336*kt*kt*kt*kt
	default:
		kd = 0.165
	}
	dhi = kd * ghi
	if dhi < 0 {
		dhi = 0
	}
	dni = (ghi - dhi) / (math.Cos(0) + 1e-15) // 简化，假设天顶角已含在 GHI 中
	if dni < 0 {
		dni = 0
	}
	return dni, dhi
}

// ExtraterrestrialIrradiance 计算大气层外辐照度。
func ExtraterrestrialIrradiance(dayOfYear int, zenithDeg float64) float64 {
	if zenithDeg >= 90 {
		return 0
	}
	eccentricity := 1 + 0.033*math.Cos(2*math.Pi*float64(dayOfYear)/365)
	cosZ := math.Cos(zenithDeg * math.Pi / 180)
	return 1367 * eccentricity * cosZ
}

// EquationOfTime 计算时差方程（分钟）。
func EquationOfTime(dayOfYear int) float64 {
	b := 2 * math.Pi * (float64(dayOfYear) - 81) / 365
	return 9.87*math.Sin(2*b) - 7.53*math.Cos(b) - 1.5*math.Sin(b)
}

// SolarTime 将标准时间转为太阳时。
func SolarTime(standardHour float64, longitude float64, timeZoneOffset float64, dayOfYear int) float64 {
	eot := EquationOfTime(dayOfYear)
	lstm := 15 * timeZoneOffset
	timeCorrection := 4*(longitude-lstm) + eot
	return standardHour + timeCorrection/60
}

// SunriseSunset 计算日出和日落的太阳时。
func SunriseSunset(latDeg float64, dayOfYear int) (sunrise, sunset float64) {
	lat := latDeg * math.Pi / 180
	dec := Declination(dayOfYear)
	cosHA := -math.Tan(lat) * math.Tan(dec)
	if cosHA >= 1 {
		return 12, 12 // 极夜
	}
	if cosHA <= -1 {
		return 0, 24 // 极昼
	}
	ha := math.Acos(cosHA) * 180 / (math.Pi * 15)
	sunrise = 12 - ha
	sunset = 12 + ha
	return sunrise, sunset
}

// DaylightHours 计算日照时长。
func DaylightHours(latDeg float64, dayOfYear int) float64 {
	rise, set := SunriseSunset(latDeg, dayOfYear)
	return set - rise
}
