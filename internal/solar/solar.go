// Package solar 实现太阳辐照在倾斜面的计算。
package solar

import "math"

// TiltedIrradiance 计算倾斜面上的总辐照度。
// ghi: 水平面全局辐照, dni: 法向直射, dhi: 散射辐照
// tiltDeg: 倾斜角度, surfAzDeg: 面板朝向, solarZenDeg: 太阳天顶角, solarAzDeg: 太阳方位角
func TiltedIrradiance(ghi, dni, dhi, tiltDeg, surfAzDeg, solarZenDeg, solarAzDeg float64) float64 {
	tilt := tiltDeg * math.Pi / 180
	surfAz := surfAzDeg * math.Pi / 180
	zenith := solarZenDeg * math.Pi / 180
	solAz := solarAzDeg * math.Pi / 180

	// 入射角
	cosAOI := math.Sin(zenith)*math.Sin(tilt)*math.Cos(solAz-surfAz) +
		math.Cos(zenith)*math.Cos(tilt)
	if cosAOI < 0 {
		cosAOI = 0
	}

	// 直射分量
	beam := dni * cosAOI

	// 散射（各向同性天空模型）
	diffuse := dhi * (1 + math.Cos(tilt)) / 2

	// 地面反射（albedo=0.2）
	albedo := 0.2
	ground := ghi * albedo * (1 - math.Cos(tilt)) / 2

	total := beam + diffuse + ground
	if total < 0 {
		return 0
	}
	return total
}

// AOI 计算入射角（度）。
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

// IAMLoss 计算入射角修正损失（ASHRAE 模型）。
// b0 通常取 0.05
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

// OptimalTilt 根据纬度估算最佳倾斜角（度）。
func OptimalTilt(latDeg float64) float64 {
	return math.Abs(latDeg) * 0.76 + 3.1
}

// Perez 计算 Perez 天空散射模型中的倾斜面散射辐照。
func Perez(dhi, dni, zenithDeg, tiltDeg, aoiDeg float64) float64 {
	if dhi <= 0 {
		return 0
	}
	tilt := tiltDeg * math.Pi / 180
	zenith := zenithDeg * math.Pi / 180

	// 简化 Perez 模型
	a := math.Max(0, math.Cos(aoiDeg*math.Pi/180))
	b := math.Max(math.Cos(85*math.Pi/180), math.Cos(zenith))

	// 环日散射
	circumsolar := dhi * a / b * 0.3

	// 各向同性
	isotropic := dhi * (1 + math.Cos(tilt)) / 2 * 0.5

	// 地平线增亮
	horizon := dhi * math.Sin(tilt) * 0.2

	return circumsolar + isotropic + horizon
}

// TranspositionFactor 计算转置系数 = POA / GHI。
func TranspositionFactor(poa, ghi float64) float64 {
	if ghi <= 0 {
		return 0
	}
	return poa / ghi
}

// DiffuseIsotropic 各向同性散射模型。
func DiffuseIsotropic(dhi, tiltDeg float64) float64 {
	tilt := tiltDeg * math.Pi / 180
	return dhi * (1 + math.Cos(tilt)) / 2
}

// GroundReflected 计算地面反射分量。
func GroundReflected(ghi, tiltDeg, albedo float64) float64 {
	tilt := tiltDeg * math.Pi / 180
	return ghi * albedo * (1 - math.Cos(tilt)) / 2
}

// EffectiveIrradiance 计算有效辐照度（含 IAM 和光谱修正）。
func EffectiveIrradiance(poa, iam, spectralFactor float64) float64 {
	eff := poa * iam * spectralFactor
	if eff < 0 {
		return 0
	}
	return eff
}

// SpectralFactor 简化光谱修正系数（基于大气质量）。
func SpectralFactor(airMass float64) float64 {
	if airMass <= 0 {
		return 1
	}
	// 四阶多项式近似
	am := airMass
	return 0.918093 + 0.086257*am - 0.024459*am*am + 0.002816*am*am*am - 0.000126*am*am*am*am
}
