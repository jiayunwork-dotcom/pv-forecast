// Package cable 提供电缆选型和损耗计算。
package cable

import "math"

// Material 电缆材料。
type Material int

const (
	Copper    Material = iota // 铜
	Aluminum                  // 铝
)

// Resistivity 返回材料电阻率 (Ω·mm²/m)。
func Resistivity(mat Material, tempC float64) float64 {
	var rho20 float64
	switch mat {
	case Copper:
		rho20 = 0.01724
	case Aluminum:
		rho20 = 0.02826
	default:
		rho20 = 0.01724
	}
	alpha := 0.00393 // 铜的温度系数
	if mat == Aluminum {
		alpha = 0.00403
	}
	return rho20 * (1 + alpha*(tempC-20))
}

// VoltageDrop 计算电压降（V）。
func VoltageDrop(current, length, crossSection float64, mat Material, tempC float64) float64 {
	if crossSection <= 0 {
		return 0
	}
	rho := Resistivity(mat, tempC)
	resistance := rho * 2 * length / crossSection
	return current * resistance
}

// VoltageDropPercent 计算电压降百分比。
func VoltageDropPercent(vDrop, nominalVoltage float64) float64 {
	if nominalVoltage <= 0 {
		return 0
	}
	return vDrop / nominalVoltage * 100
}

// PowerLoss 计算线路功率损耗（W）。
func PowerLoss(current, length, crossSection float64, mat Material, tempC float64) float64 {
	if crossSection <= 0 {
		return 0
	}
	rho := Resistivity(mat, tempC)
	resistance := rho * 2 * length / crossSection
	return current * current * resistance
}

// MinCrossSection 计算满足压降限制的最小截面积。
func MinCrossSection(current, length, maxDropV float64, mat Material, tempC float64) float64 {
	if maxDropV <= 0 {
		return 0
	}
	rho := Resistivity(mat, tempC)
	return rho * 2 * length * current / maxDropV
}

// StandardSize 返回不小于计算值的标准截面积。
func StandardSize(minArea float64) float64 {
	standards := []float64{1.5, 2.5, 4, 6, 10, 16, 25, 35, 50, 70, 95, 120, 150, 185, 240, 300}
	for _, s := range standards {
		if s >= minArea {
			return s
		}
	}
	return 300
}

// Ampacity 返回标准截面积的载流量（空气中，铜）。
func Ampacity(crossSection float64) float64 {
	table := map[float64]float64{
		1.5: 18, 2.5: 25, 4: 34, 6: 44, 10: 61, 16: 82, 25: 108,
		35: 135, 50: 168, 70: 213, 95: 258, 120: 299, 150: 344, 185: 392, 240: 461, 300: 530,
	}
	if amp, ok := table[crossSection]; ok {
		return amp
	}
	// 近似
	return 10 * math.Sqrt(crossSection)
}

// TotalCableLoss 计算整个 PV 系统的电缆损耗比例。
func TotalCableLoss(dcLoss, acLoss, totalPower float64) float64 {
	if totalPower <= 0 {
		return 0
	}
	return (dcLoss + acLoss) / totalPower
}
