// Package degradation 实现光伏组件衰退模型。
package degradation

import "math"

// Mode 衰退模式。
type Mode int

const (
	Linear      Mode = iota // 线性衰退
	Exponential             // 指数衰退
	StepWise                // 阶梯式衰退
	TwoPhase                // 两阶段衰退（初始快速 + 后续缓慢）
)

// Model 衰退模型参数。
type Model struct {
	Mode      Mode
	AnnualPct float64 // 年衰退百分比（如 0.5 表示 0.5%/年）
	Phase1Pct float64 // 第一阶段年衰退（用于 TwoPhase）
	Phase1Yr  int     // 第一阶段年数
	StepYears int     // 阶梯间隔年数
	StepPct   float64 // 每阶梯衰退百分比
}

// DefaultLinear 返回默认线性衰退模型（0.5%/年）。
func DefaultLinear() Model {
	return Model{Mode: Linear, AnnualPct: 0.5}
}

// DefaultTwoPhase 返回默认两阶段衰退模型。
func DefaultTwoPhase() Model {
	return Model{
		Mode:      TwoPhase,
		Phase1Pct: 2.0,
		Phase1Yr:  1,
		AnnualPct: 0.5,
	}
}

// Factor 计算第 year 年的性能因子（0-1），year 从 0 开始。
func (m Model) Factor(year int) float64 {
	if year < 0 {
		return 1.0
	}
	switch m.Mode {
	case Linear:
		f := 1 - float64(year)*m.AnnualPct/100
		if f < 0 {
			return 0
		}
		return f
	case Exponential:
		return math.Pow(1-m.AnnualPct/100, float64(year))
	case StepWise:
		if m.StepYears <= 0 {
			return 1.0
		}
		steps := year / m.StepYears
		f := 1 - float64(steps)*m.StepPct/100
		if f < 0 {
			return 0
		}
		return f
	case TwoPhase:
		if year <= m.Phase1Yr {
			f := 1 - float64(year)*m.Phase1Pct/100
			if f < 0 {
				return 0
			}
			return f
		}
		phase1Loss := float64(m.Phase1Yr) * m.Phase1Pct / 100
		remaining := float64(year-m.Phase1Yr) * m.AnnualPct / 100
		f := 1 - phase1Loss - remaining
		if f < 0 {
			return 0
		}
		return f
	}
	return 1.0
}

// CumulativeEnergy 计算 n 年内考虑衰退后的累计发电量占比。
func (m Model) CumulativeEnergy(years int) float64 {
	if years <= 0 {
		return 0
	}
	total := 0.0
	for y := 0; y < years; y++ {
		total += m.Factor(y)
	}
	return total / float64(years)
}

// WarrantyCheck 检查在保修年限末是否满足最低性能保证。
func (m Model) WarrantyCheck(warrantyYears int, minPct float64) bool {
	f := m.Factor(warrantyYears)
	return f*100 >= minPct
}

// LIDLoss 计算光致衰退损失（首年）。
func LIDLoss(lidPct float64) float64 {
	return 1 - lidPct/100
}

// PIDLoss 计算电位诱导衰退损失。
func PIDLoss(voltage, threshold float64, susceptibility float64) float64 {
	if voltage <= threshold {
		return 1.0
	}
	loss := 1 - susceptibility*(voltage-threshold)/1000
	if loss < 0 {
		return 0
	}
	return loss
}

// SoilingLoss 计算脏污导致的性能损失。
func SoilingLoss(daysSinceCleaning int, dailyRate float64) float64 {
	loss := float64(daysSinceCleaning) * dailyRate / 100
	if loss > 0.3 { // 最大 30% 损失
		loss = 0.3
	}
	return 1 - loss
}

// SnowLoss 计算积雪覆盖损失。
func SnowLoss(coveragePct float64) float64 {
	if coveragePct <= 0 {
		return 1.0
	}
	if coveragePct >= 100 {
		return 0
	}
	return 1 - coveragePct/100
}

// MismatchLoss 计算组件不匹配损失。
func MismatchLoss(coefficientOfVariation float64) float64 {
	loss := 2 * coefficientOfVariation * coefficientOfVariation
	if loss > 0.1 {
		loss = 0.1
	}
	return 1 - loss
}

// TotalSystemLoss 综合计算系统总损失因子。
func TotalSystemLoss(degradation, soiling, mismatch, wiring, availability float64) float64 {
	return degradation * soiling * mismatch * wiring * availability
}
