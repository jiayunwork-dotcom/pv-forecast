// Package inverter 逆变器建模与效率曲线。
package inverter

import "math"

// Spec 逆变器规格参数。
type Spec struct {
	PacMax   float64 // 额定交流功率 (W)
	PdcMax   float64 // 最大直流输入功率 (W)
	VdcMax   float64 // 最大直流电压 (V)
	VdcMin   float64 // 最小直流电压 (V)
	Vmpp     float64 // 最大功率点电压 (V)
	EtaPeak  float64 // 峰值效率
	EtaEU    float64 // 欧洲效率
	Pnt      float64 // 夜间功耗 (W)
	C0       float64 // Sandia 模型系数
	C1       float64
	C2       float64
	C3       float64
}

// DefaultSpec 返回典型 10kW 逆变器参数。
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

// Efficiency 计算给定负载比下的逆变器效率。
func Efficiency(loadRatio float64, spec Spec) float64 {
	if loadRatio <= 0 {
		return 0
	}
	if loadRatio > 1 {
		loadRatio = 1
	}
	// 简化效率曲线：低负载低效率，峰值在 ~30-50% 负载
	eta := spec.EtaPeak * (1 - 0.05/loadRatio - 0.02*loadRatio)
	if eta < 0 {
		return 0
	}
	if eta > spec.EtaPeak {
		return spec.EtaPeak
	}
	return eta
}

// SandiaModel Sandia 逆变器性能模型。
func SandiaModel(pdc, vdc float64, spec Spec) float64 {
	if pdc <= 0 {
		return -spec.Pnt
	}
	a := spec.PdcMax * (1 + spec.C0*(vdc-spec.Vmpp))
	b := spec.Pnt * (1 + spec.C1*(vdc-spec.Vmpp))
	c := spec.C2 * (vdc - spec.Vmpp)
	pac := (spec.PacMax/(a-b) - c*(a-b)) * (pdc - b) + c*(pdc-b)*(pdc-b)
	if pac < -spec.Pnt {
		return -spec.Pnt
	}
	if pac > spec.PacMax {
		return spec.PacMax
	}
	return pac
}

// Clipping 计算削峰损失（当 DC 输入超过逆变器容量时）。
func Clipping(pdc, pacMax float64) float64 {
	return applyClip(pdc, pacMax)
}

// EUEfficiency 计算欧洲加权效率。
func EUEfficiency(spec Spec) float64 {
	loads := []float64{0.05, 0.10, 0.20, 0.30, 0.50, 1.00}
	weights := []float64{0.03, 0.06, 0.13, 0.10, 0.48, 0.20}
	etaEU := 0.0
	for i, load := range loads {
		etaEU += weights[i] * Efficiency(load, spec)
	}
	return etaEU
}

// CECEfficiency 计算 CEC 加权效率。
func CECEfficiency(spec Spec) float64 {
	loads := []float64{0.10, 0.20, 0.30, 0.50, 0.75, 1.00}
	weights := []float64{0.04, 0.05, 0.12, 0.21, 0.53, 0.05}
	eta := 0.0
	for i, load := range loads {
		eta += weights[i] * Efficiency(load, spec)
	}
	return eta
}

// MPPTRange 检查电压是否在 MPPT 范围内。
func MPPTRange(vdc float64, spec Spec) bool {
	return vdc >= spec.VdcMin && vdc <= spec.VdcMax
}

// DCACRatio 返回 DC/AC 比。
func DCACRatio(pdcPeak, pacMax float64) float64 {
	if pacMax <= 0 {
		return 0
	}
	return pdcPeak / pacMax
}

// AnnualEnergy 估算年发电量 (kWh)。
func AnnualEnergy(peakHours, systemKW, eta float64) float64 {
	return peakHours * systemKW * eta * 365
}

// TempDerating 温度降额。
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

// StringSizing 计算串联组件数范围。
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
