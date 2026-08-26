package cleaning

import "math"

func OptimalInterval(soilingRate, cleaningCost, energyPrice, dailyGen float64) float64 {
	if soilingRate <= 0 || energyPrice <= 0 || dailyGen <= 0 {
		return 365
	}
	interval := math.Sqrt(2 * cleaningCost / (soilingRate / 100 * energyPrice * dailyGen))
	if interval < 7 {
		interval = 7
	}
	if interval > 180 {
		interval = 180
	}
	return interval
}

func CleaningROI(cleaningCost, soilingRate, daysSinceLast float64, dailyGen, energyPrice float64) float64 {
	lossRecovered := soilingRate / 100 * daysSinceLast * dailyGen * energyPrice
	if cleaningCost <= 0 {
		return 0
	}
	return (lossRecovered - cleaningCost) / cleaningCost
}

func SoilingRate(rainMmPerDay, tempC float64) float64 {
	baseRate := 0.3
	rainFactor := 1.0
	if rainMmPerDay > 0 {
		rainFactor = math.Max(0, 1-rainMmPerDay/5)
	}
	tempFactor := 1.0
	if tempC > 30 {
		tempFactor = 1 + 0.02*(tempC-30)
	}
	return baseRate * rainFactor * tempFactor
}

func RainCleaning(currentSoilingPct, rainMm float64) float64 {
	if rainMm <= 0 {
		return currentSoilingPct
	}
	cleaned := currentSoilingPct * math.Max(0, 1-rainMm/5)
	return cleaned
}

func Schedule(totalDays int, intervalDays, graceDays int) []int {
	if intervalDays <= 0 {
		intervalDays = 30
	}
	var schedule []int
	day := graceDays
	for day < totalDays {
		schedule = append(schedule, day)
		day += intervalDays
	}
	return schedule
}
