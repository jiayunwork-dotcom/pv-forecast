package cable

import "math"

type Material int

const (
	Copper Material = iota
	Aluminum
)

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
	alpha := 0.00393
	if mat == Aluminum {
		alpha = 0.00403
	}
	return rho20 * (1 + alpha*(tempC-20))
}

func VoltageDrop(current, length, crossSection float64, mat Material, tempC float64) float64 {
	if crossSection <= 0 {
		return 0
	}
	rho := Resistivity(mat, tempC)
	resistance := rho * 2 * length / crossSection
	return current * resistance
}

func VoltageDropPercent(vDrop, nominalVoltage float64) float64 {
	if nominalVoltage <= 0 {
		return 0
	}
	return vDrop / nominalVoltage * 100
}

func PowerLoss(current, length, crossSection float64, mat Material, tempC float64) float64 {
	if crossSection <= 0 {
		return 0
	}
	rho := Resistivity(mat, tempC)
	resistance := rho * 2 * length / crossSection
	return current * current * resistance
}

func MinCrossSection(current, length, maxDropV float64, mat Material, tempC float64) float64 {
	if maxDropV <= 0 {
		return 0
	}
	rho := Resistivity(mat, tempC)
	return rho * 2 * length * current / maxDropV
}

func StandardSize(minArea float64) float64 {
	standards := []float64{1.5, 2.5, 4, 6, 10, 16, 25, 35, 50, 70, 95, 120, 150, 185, 240, 300}
	for _, s := range standards {
		if s >= minArea {
			return s
		}
	}
	return 300
}

func Ampacity(crossSection float64) float64 {
	table := map[float64]float64{
		1.5: 18, 2.5: 25, 4: 34, 6: 44, 10: 61, 16: 82, 25: 108,
		35: 135, 50: 168, 70: 213, 95: 258, 120: 299, 150: 344, 185: 392, 240: 461, 300: 530,
	}
	if amp, ok := table[crossSection]; ok {
		return amp
	}
	return 10 * math.Sqrt(crossSection)
}

func TotalCableLoss(dcLoss, acLoss, totalPower float64) float64 {
	if totalPower <= 0 {
		return 0
	}
	return (dcLoss + acLoss) / totalPower
}
