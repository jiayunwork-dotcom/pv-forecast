// Package sizing 提供光伏系统容量设计工具。
package sizing

import "math"

// SystemSpec 系统规格。
type SystemSpec struct {
	ModuleWp       float64 // 单块组件峰值功率 (Wp)
	ModuleVoc      float64 // 开路电压 (V)
	ModuleIsc      float64 // 短路电流 (A)
	ModuleVmpp     float64 // 最大功率点电压 (V)
	ModuleImpp     float64 // 最大功率点电流 (A)
	InverterPacMax float64 // 逆变器额定功率 (W)
	InverterVdcMin float64 // MPPT 最小电压 (V)
	InverterVdcMax float64 // MPPT 最大电压 (V)
	InverterIdcMax float64 // 最大输入电流 (A)
	TempCoeffVoc   float64 // Voc 温度系数 (%/°C)
	TempCoeffIsc   float64 // Isc 温度系数 (%/°C)
	TempMin        float64 // 设计最低温度 (°C)
	TempMax        float64 // 设计最高温度 (°C)
}

// DefaultSpec 返回典型 400Wp 组件 + 10kW 逆变器规格。
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

// StringSize 计算每串组件数范围。
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

// ParallelStrings 计算可并联串数。
func ParallelStrings(spec SystemSpec) int {
	iscHot := spec.ModuleIsc * (1 + spec.TempCoeffIsc/100*(spec.TempMax-25))
	if iscHot <= 0 {
		return 0
	}
	return int(math.Floor(spec.InverterIdcMax / iscHot))
}

// DCACRatio 计算 DC/AC 比。
func DCACRatio(modulesPerString, strings int, spec SystemSpec) float64 {
	dcPeak := float64(modulesPerString*strings) * spec.ModuleWp
	if spec.InverterPacMax <= 0 {
		return 0
	}
	return dcPeak / spec.InverterPacMax
}

// TotalCapacityKW 计算系统总 DC 容量。
func TotalCapacityKW(modulesPerString, strings int, spec SystemSpec) float64 {
	return float64(modulesPerString*strings) * spec.ModuleWp / 1000
}

// AreaRequired 计算所需安装面积。
func AreaRequired(modulesPerString, strings int, moduleAreaM2, gcr float64) float64 {
	totalModules := float64(modulesPerString * strings)
	moduleArea := totalModules * moduleAreaM2
	if gcr <= 0 {
		gcr = 0.4
	}
	return moduleArea / gcr
}

// AnnualYield 估算年发电量 (kWh)。
func AnnualYield(capacityKWp, peakSunHours, systemLossFactor float64) float64 {
	return capacityKWp * peakSunHours * 365 * systemLossFactor
}

// LCOE 平准化度电成本 (元/kWh)。
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

// PaybackYears 计算投资回收期。
func PaybackYears(totalCost, annualSavings float64) float64 {
	if annualSavings <= 0 {
		return math.Inf(1)
	}
	return totalCost / annualSavings
}

// CableSize 根据电流和距离推算电缆截面积。
func CableSize(current, lengthM, maxDropPct, voltage float64) float64 {
	if voltage <= 0 || maxDropPct <= 0 {
		return 0
	}
	resistivity := 0.0175 // 铜 Ω·mm²/m
	maxDrop := voltage * maxDropPct / 100
	area := 2 * resistivity * lengthM * current / maxDrop
	return area
}

// ShadeFreeSpacing 计算避免遮挡所需行间距。
func ShadeFreeSpacing(moduleHeight, tiltDeg, minElevDeg float64) float64 {
	if minElevDeg <= 0 {
		return moduleHeight * 5
	}
	tilt := tiltDeg * math.Pi / 180
	elev := minElevDeg * math.Pi / 180
	return moduleHeight * math.Sin(tilt) / math.Tan(elev)
}
