package storage

import "math"

type Battery struct {
	CapacityKWh    float64
	MaxChargeKW    float64
	MaxDischargeKW float64
	ChargeEff      float64
	DischargeEff   float64
	SOCMin         float64
	SOCMax         float64
	DegradePctYear float64
}

type State struct {
	SOC     float64
	CycleN  int
	AgeYear float64
}

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

func AvailableCapacity(bat Battery, state State) float64 {
	degrade := 1 - bat.DegradePctYear/100*state.AgeYear
	if degrade < 0 {
		degrade = 0
	}
	return bat.CapacityKWh * degrade * (bat.SOCMax - bat.SOCMin)
}

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

func Autarky(load, gridImport float64) float64 {
	if load <= 0 {
		return 0
	}
	return 1 - gridImport/load
}

func RoundTripEfficiency(bat Battery) float64 {
	return bat.ChargeEff * bat.DischargeEff
}

func CycleLife(dodPct float64) int {
	if dodPct <= 0 {
		return 0
	}
	return int(4000 * math.Pow(80/dodPct, 1.2))
}

func PaybackYears(batteryCost, annualSavings float64) float64 {
	if annualSavings <= 0 {
		return math.Inf(1)
	}
	return batteryCost / annualSavings
}
