// Package monitoring 提供电站运维监控告警。
package monitoring

import "math"

// Alarm 告警类型。
type Alarm struct {
	Level   AlarmLevel
	Code    string
	Message string
	Value   float64
}

// AlarmLevel 告警等级。
type AlarmLevel int

const (
	Info     AlarmLevel = iota
	Warning
	Critical
)

// CheckPR 检查性能比是否过低。
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

// CheckInverterOutput 检查逆变器输出异常。
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

// CheckStringCurrent 检查组串电流偏差。
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

// CheckCommunication 检查通信中断。
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

// AvailabilityRate 计算设备可用率。
func AvailabilityRate(totalHours, downHours float64) float64 {
	if totalHours <= 0 {
		return 0
	}
	return (totalHours - downHours) / totalHours
}

// EnergyBalance 检查发电量与计量偏差。
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
