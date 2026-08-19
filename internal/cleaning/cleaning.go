// Package cleaning 提供清洗策略优化。
package cleaning

import "math"

// OptimalInterval 计算最佳清洗间隔（天）。
// soilingRate: 日脏污率 (%/天)
// cleaningCost: 每次清洗费用
// energyPrice: 电价（元/kWh）
// dailyGen: 日发电量（kWh）
func OptimalInterval(soilingRate, cleaningCost, energyPrice, dailyGen float64) float64 {
	if soilingRate <= 0 || energyPrice <= 0 || dailyGen <= 0 {
		return 365 // 无意义时按年清洗
	}
	// 最优间隔 = sqrt(2 * cost / (rate * price * gen))
	interval := math.Sqrt(2 * cleaningCost / (soilingRate / 100 * energyPrice * dailyGen))
	if interval < 7 {
		interval = 7
	}
	if interval > 180 {
		interval = 180
	}
	return interval
}

// CleaningROI 计算一次清洗的投资回报率。
func CleaningROI(cleaningCost, soilingRate, daysSinceLast float64, dailyGen, energyPrice float64) float64 {
	lossRecovered := soilingRate / 100 * daysSinceLast * dailyGen * energyPrice
	if cleaningCost <= 0 {
		return 0
	}
	return (lossRecovered - cleaningCost) / cleaningCost
}

// SoilingRate 根据气象条件估算日脏污率 (%/天)。
func SoilingRate(rainMmPerDay, tempC float64) float64 {
	baseRate := 0.3 // 干燥地区基准 0.3%/天
	// 雨水自然清洗
	rainFactor := 1.0
	if rainMmPerDay > 0 {
		rainFactor = math.Max(0, 1-rainMmPerDay/5)
	}
	// 高温加速（沙尘）
	tempFactor := 1.0
	if tempC > 30 {
		tempFactor = 1 + 0.02*(tempC-30)
	}
	return baseRate * rainFactor * tempFactor
}

// RainCleaning 雨水清洗效果：如果降雨量超过阈值，减少脏污。
func RainCleaning(currentSoilingPct, rainMm float64) float64 {
	if rainMm <= 0 {
		return currentSoilingPct
	}
	cleaned := currentSoilingPct * math.Max(0, 1-rainMm/5)
	return cleaned
}

// Schedule 生成年度清洗计划（返回第几天清洗）。
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
