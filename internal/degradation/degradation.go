package degradation

import "math"

type Mode int

const (
	Linear Mode = iota
	Exponential
	StepWise
	TwoPhase
)

type Model struct {
	Mode      Mode
	AnnualPct float64
	Phase1Pct float64
	Phase1Yr  int
	StepYears int
	StepPct   float64
}

func DefaultLinear() Model {
	return Model{Mode: Linear, AnnualPct: 0.5}
}

func DefaultTwoPhase() Model {
	return Model{
		Mode:      TwoPhase,
		Phase1Pct: 2.0,
		Phase1Yr:  1,
		AnnualPct: 0.5,
	}
}

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

func (m Model) WarrantyCheck(warrantyYears int, minPct float64) bool {
	f := m.Factor(warrantyYears)
	return f*100 >= minPct
}

func LIDLoss(lidPct float64) float64 {
	return 1 - lidPct/100
}

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

func SoilingLoss(daysSinceCleaning int, dailyRate float64) float64 {
	loss := float64(daysSinceCleaning) * dailyRate / 100
	if loss > 0.3 {
		loss = 0.3
	}
	return 1 - loss
}

func SnowLoss(coveragePct float64) float64 {
	if coveragePct <= 0 {
		return 1.0
	}
	if coveragePct >= 100 {
		return 0
	}
	return 1 - coveragePct/100
}

func MismatchLoss(coefficientOfVariation float64) float64 {
	loss := 2 * coefficientOfVariation * coefficientOfVariation
	if loss > 0.1 {
		loss = 0.1
	}
	return 1 - loss
}

func TotalSystemLoss(degradation, soiling, mismatch, wiring, availability float64) float64 {
	return degradation * soiling * mismatch * wiring * availability
}
