package cable

import (
	"math"
	"testing"
)

func TestResistivity_Copper20(t *testing.T) {
	r := Resistivity(Copper, 20)
	if math.Abs(r-0.01724) > 1e-6 {
		t.Fatalf("copper at 20C expected 0.01724, got %v", r)
	}
}

func TestResistivity_Aluminum(t *testing.T) {
	r := Resistivity(Aluminum, 20)
	if math.Abs(r-0.02826) > 1e-6 {
		t.Fatalf("aluminum at 20C expected 0.02826, got %v", r)
	}
}

func TestVoltageDrop(t *testing.T) {
	vd := VoltageDrop(20, 50, 6, Copper, 25)
	if vd < 3 || vd > 10 {
		t.Fatalf("voltage drop out of range: %v", vd)
	}
}

func TestVoltageDropPercent(t *testing.T) {
	pct := VoltageDropPercent(5, 400)
	if math.Abs(pct-1.25) > 1e-10 {
		t.Fatalf("expected 1.25%%, got %v", pct)
	}
}

func TestPowerLoss(t *testing.T) {
	loss := PowerLoss(10, 30, 4, Copper, 25)
	if loss <= 0 {
		t.Fatalf("power loss should be positive")
	}
}

func TestMinCrossSection(t *testing.T) {
	cs := MinCrossSection(20, 50, 8, Copper, 25)
	if cs < 1 || cs > 10 {
		t.Fatalf("cross section out of range: %v", cs)
	}
}

func TestStandardSize(t *testing.T) {
	s := StandardSize(3.5)
	if s != 4 {
		t.Fatalf("expected 4, got %v", s)
	}
	s = StandardSize(100)
	if s != 120 {
		t.Fatalf("expected 120, got %v", s)
	}
}

func TestAmpacity(t *testing.T) {
	a := Ampacity(10)
	if a != 61 {
		t.Fatalf("expected 61A for 10mm², got %v", a)
	}
}

func TestTotalCableLoss(t *testing.T) {
	loss := TotalCableLoss(100, 50, 10000)
	if math.Abs(loss-0.015) > 1e-10 {
		t.Fatalf("expected 0.015, got %v", loss)
	}
}
