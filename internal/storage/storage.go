// Package storage 提供光伏储能系统建模。
package storage

import "math"

// Battery 电池参数。
type Battery struct {
	CapacityKWh    float64 // 额定容量
	MaxChargeKW    float64 // 最大充电功率
	MaxDischargeKW float64 // 最大放电功率
	ChargeEff      float64 // 充电效率
	DischargeEff   float64 // 放电效率
	SOCMin         float64 // 最小 SOC
	SOCMax         float64 // 最大 SOC
	DegradePctYear float64 // 年衰退率 (%)
}

// State 电池当前状态。
type State struct {
	SOC     float64 // 当前 SOC (0-1)
	CycleN  int     // 循环次数
	AgeYear float64 // 使用年限
}

// DefaultLFP 返回默认磷酸铁锂电池参数。
func DefaultLFP() Battery {
	return Battery{
		CapacityKWh:    10,
		MaxChargeKW:    5,
		MaxDischargeKW: 5,
		ChargeEff:      0.95,
		DischargeEff:   0.95,
		SOCMin:         0.1,
		SOCMax:         0.9,
		DegradePctYear: 2.0,
	}
}

// AvailableCapacity 返回当前可用容量。
func AvailableCapacity(bat Battery, state State) float64 {
	degrade := 1 - bat.DegradePctYear/100*state.AgeYear
	if degrade < 0 {
		degrade = 0
	}
	return bat.CapacityKWh * degrade * (bat.SOCMax - bat.SOCMin)
}

// ChargeStep 充电一步（timestep 小时），返回更新后 SOC 和实际充入能量。
func ChargeStep(bat Battery, state State, powerKW, timestepH float64) (newSOC, energyIn float64) {
	if powerKW <= 0 {
		return state.SOC, 0
	}
	if powerKW > bat.MaxChargeKW {
		powerKW = bat.MaxChargeKW
	}
	energy := powerKW * timestepH * bat.ChargeEff
	maxEnergy := (bat.SOCMax - state.SOC) * bat.CapacityKWh
	if energy > maxEnergy {
		energy = maxEnergy
	}
	newSOC = state.SOC + energy/bat.CapacityKWh
	return newSOC, energy
}

// DischargeStep 放电一步，返回更新后 SOC 和实际输出能量。
func DischargeStep(bat Battery, state State, powerKW, timestepH float64) (newSOC, energyOut float64) {
	if powerKW <= 0 {
		return state.SOC, 0
	}
	if powerKW > bat.MaxDischargeKW {
		powerKW = bat.MaxDischargeKW
	}
	energyNeeded := powerKW * timestepH
	maxEnergy := (state.SOC - bat.SOCMin) * bat.CapacityKWh
	if energyNeeded > maxEnergy {
		energyNeeded = maxEnergy
	}
	newSOC = state.SOC - energyNeeded/bat.CapacityKWh
	energyOut = energyNeeded * bat.DischargeEff
	return newSOC, energyOut
}

// SelfConsumption 计算自消纳率。
func SelfConsumption(load, pvGen, gridExport float64) float64 {
	if pvGen <= 0 {
		return 0
	}
	consumed := pvGen - gridExport
	if consumed < 0 {
		consumed = 0
	}
	return consumed / pvGen
}

// Autarky 计算自给率。
func Autarky(load, gridImport float64) float64 {
	if load <= 0 {
		return 0
	}
	return 1 - gridImport/load
}

// RoundTripEfficiency 计算一次充放循环效率。
func RoundTripEfficiency(bat Battery) float64 {
	return bat.ChargeEff * bat.DischargeEff
}

// CycleLife 估算在给定 DOD 下的循环寿命。
func CycleLife(dodPct float64) int {
	if dodPct <= 0 {
		return 0
	}
	// 简化：LFP 在 80% DOD 约 4000 次
	return int(4000 * math.Pow(80/dodPct, 1.2))
}

// PaybackYears 储能投资回收期。
func PaybackYears(batteryCost, annualSavings float64) float64 {
	if annualSavings <= 0 {
		return math.Inf(1)
	}
	return batteryCost / annualSavings
}
