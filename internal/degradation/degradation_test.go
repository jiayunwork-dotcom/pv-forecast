package degradation

import (
	"math"
	"testing"
)

func TestLinear_Year0(t *testing.T) {
	m := DefaultLinear()
	if m.Factor(0) != 1.0 {
		t.Fatalf("year 0 should be 1.0")
	}
}

func TestLinear_Year25(t *testing.T) {
	m := DefaultLinear()
	f := m.Factor(25)
	expected := 1 - 25*0.005
	if math.Abs(f-expected) > 1e-10 {
		t.Fatalf("year 25 expected %v, got %v", expected, f)
	}
}

func TestExponential(t *testing.T) {
	m := Model{Mode: Exponential, AnnualPct: 0.7}
	f := m.Factor(10)
	expected := math.Pow(0.993, 10)
	if math.Abs(f-expected) > 1e-6 {
		t.Fatalf("exponential year 10 expected %v, got %v", expected, f)
	}
}

func TestStepWise(t *testing.T) {
	m := Model{Mode: StepWise, StepYears: 5, StepPct: 2.0}
	f10 := m.Factor(10)
	// 10/5 = 2 steps => 1 - 2*0.02 = 0.96
	if math.Abs(f10-0.96) > 1e-10 {
		t.Fatalf("stepwise year 10 expected 0.96, got %v", f10)
	}
}

func TestTwoPhase(t *testing.T) {
	m := DefaultTwoPhase()
	f0 := m.Factor(0)
	f1 := m.Factor(1)
	f5 := m.Factor(5)
	if f0 != 1.0 {
		t.Fatalf("year 0 should be 1")
	}
	if math.Abs(f1-0.98) > 1e-10 {
		t.Fatalf("year 1 expected 0.98, got %v", f1)
	}
	// year 5: phase1 loss = 0.02, phase2 loss = 4*0.005 = 0.02 => f=0.96
	if math.Abs(f5-0.96) > 1e-10 {
		t.Fatalf("year 5 expected 0.96, got %v", f5)
	}
}

func TestCumulativeEnergy(t *testing.T) {
	m := DefaultLinear()
	ce := m.CumulativeEnergy(25)
	if ce < 0.9 || ce > 1.0 {
		t.Fatalf("cumulative energy expected 0.9-1.0, got %v", ce)
	}
}

func TestWarrantyCheck(t *testing.T) {
	m := DefaultLinear()
	if !m.WarrantyCheck(25, 80) {
		t.Fatalf("should pass warranty at 25yr/80%%")
	}
	if m.WarrantyCheck(25, 90) {
		t.Fatalf("should fail warranty at 25yr/90%%")
	}
}

func TestLIDLoss(t *testing.T) {
	f := LIDLoss(1.5)
	if math.Abs(f-0.985) > 1e-10 {
		t.Fatalf("LID 1.5%% expected 0.985, got %v", f)
	}
}

func TestSoilingLoss(t *testing.T) {
	f := SoilingLoss(10, 0.5)
	if math.Abs(f-0.95) > 1e-10 {
		t.Fatalf("soiling expected 0.95, got %v", f)
	}
}

func TestSnowLoss(t *testing.T) {
	if SnowLoss(0) != 1.0 {
		t.Fatalf("no snow should be 1.0")
	}
	if SnowLoss(100) != 0 {
		t.Fatalf("full coverage should be 0")
	}
	if math.Abs(SnowLoss(50)-0.5) > 1e-10 {
		t.Fatalf("50%% coverage should be 0.5")
	}
}

func TestMismatchLoss(t *testing.T) {
	f := MismatchLoss(0.1)
	if math.Abs(f-0.98) > 1e-10 {
		t.Fatalf("expected 0.98, got %v", f)
	}
}

func TestTotalSystemLoss(t *testing.T) {
	total := TotalSystemLoss(0.95, 0.97, 0.98, 0.99, 0.99)
	if total < 0.85 || total > 0.95 {
		t.Fatalf("total loss factor out of range: %v", total)
	}
}
