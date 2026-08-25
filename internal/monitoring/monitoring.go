package monitoring

import "math"

type Alarm struct {
	Level   AlarmLevel
	Code    string
	Message string
	Value   float64
}

type AlarmLevel int

const (
	Info AlarmLevel = iota
	Warning
	Critical
)

func CheckPR(pr, threshold float64) *Alarm {
	if pr < threshold {
		return &Alarm{
			Level:   Warning,
			Code:    "LOW_PR",
			Message: "Performance ratio below threshold",
			Value:   pr,
		}
	}
	return nil
}

func CheckInverterOutput(actual, expected, tolerance float64) *Alarm {
	if expected <= 0 {
		return nil
	}
	deviation := math.Abs(actual-expected) / expected
	if deviation > tolerance {
		return &Alarm{
			Level:   Warning,
			Code:    "INV_DEVIATION",
			Message: "Inverter output deviates from expected",
			Value:   deviation,
		}
	}
	return nil
}

func CheckStringCurrent(currents []float64, maxDevPct float64) []Alarm {
	if len(currents) == 0 {
		return nil
	}
	mean := 0.0
	for _, c := range currents {
		mean += c
	}
	mean /= float64(len(currents))
	if mean <= 0 {
		return nil
	}
	var alarms []Alarm
	for i, c := range currents {
		dev := math.Abs(c-mean) / mean * 100
		if dev > maxDevPct {
			alarms = append(alarms, Alarm{
				Level:   Warning,
				Code:    "STRING_DEVIATION",
				Message: "String current deviation",
				Value:   float64(i),
			})
		}
	}
	return alarms
}

func CheckCommunication(lastSeenMinutes, thresholdMinutes int) *Alarm {
	if lastSeenMinutes > thresholdMinutes {
		return &Alarm{
			Level:   Critical,
			Code:    "COMM_LOSS",
			Message: "Communication lost",
			Value:   float64(lastSeenMinutes),
		}
	}
	return nil
}

func AvailabilityRate(totalHours, downHours float64) float64 {
	if totalHours <= 0 {
		return 0
	}
	return (totalHours - downHours) / totalHours
}

func EnergyBalance(meterKWh, inverterKWh, tolerancePct float64) *Alarm {
	if inverterKWh <= 0 {
		return nil
	}
	dev := math.Abs(meterKWh-inverterKWh) / inverterKWh * 100
	if dev > tolerancePct {
		return &Alarm{
			Level:   Warning,
			Code:    "ENERGY_IMBALANCE",
			Message: "Energy meter and inverter readings differ",
			Value:   dev,
		}
	}
	return nil
}
