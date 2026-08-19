package storage

import (
	"math"
	"testing"
)

func TestAvailableCapacity(t *testing.T) {
	bat := DefaultLFP()
	state := State{SOC: 0.5}
	cap := AvailableCapacity(bat, state)
	expected := 10 * 0.8 // (0.9-0.1)*10
	if math.Abs(cap-expected) > 1e-10 {
		t.Fatalf("expected %v, got %v", expected, cap)
	}
}

func TestChargeStep_Full(t *testing.T) {
	bat := DefaultLFP()
	state := State{SOC: 0.85}
	soc, energy := ChargeStep(bat, state, 10, 1)
	if soc > bat.SOCMax+1e-10 {
		t.Fatalf("SOC exceeds max: %v", soc)
	}
	if energy <= 0 {
		t.Fatalf("should charge some energy")
	}
}

func TestDischargeStep_Empty(t *testing.T) {
	bat := DefaultLFP()
	state := State{SOC: 0.15}
	soc, energy := DischargeStep(bat, state, 10, 1)
	if soc < bat.SOCMin-1e-10 {
		t.Fatalf("SOC below min: %v", soc)
	}
	if energy <= 0 {
		t.Fatalf("should discharge some energy")
	}
}

func TestSelfConsumption(t *testing.T) {
	sc := SelfConsumption(100, 80, 20)
	if math.Abs(sc-0.75) > 1e-10 {
		t.Fatalf("expected 0.75, got %v", sc)
	}
}

func TestAutarky(t *testing.T) {
	a := Autarky(100, 30)
	if math.Abs(a-0.7) > 1e-10 {
		t.Fatalf("expected 0.7, got %v", a)
	}
}

func TestRoundTripEfficiency(t *testing.T) {
	bat := DefaultLFP()
	rte := RoundTripEfficiency(bat)
	if math.Abs(rte-0.9025) > 1e-10 {
		t.Fatalf("expected 0.9025, got %v", rte)
	}
}

func TestCycleLife(t *testing.T) {
	life := CycleLife(80)
	if life < 3500 || life > 4500 {
		t.Fatalf("expected ~4000, got %d", life)
	}
}

func TestPaybackYears(t *testing.T) {
	pb := PaybackYears(50000, 10000)
	if math.Abs(pb-5.0) > 1e-10 {
		t.Fatalf("expected 5, got %v", pb)
	}
}
