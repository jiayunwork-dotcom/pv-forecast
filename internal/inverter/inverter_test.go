package inverter

import (
	"math"
	"testing"
)

func TestEfficiency_Peak(t *testing.T) {
	spec := DefaultSpec()
	eta := Efficiency(0.4, spec)
	if eta < 0.8 || eta > spec.EtaPeak {
		t.Fatalf("efficiency at 40%% load expected 0.8-peak, got %v", eta)
	}
}

func TestEfficiency_Zero(t *testing.T) {
	spec := DefaultSpec()
	eta := Efficiency(0, spec)
	if eta != 0 {
		t.Fatalf("zero load efficiency should be 0, got %v", eta)
	}
}

func TestClipping(t *testing.T) {
	loss := Clipping(12000, 10000)
	if loss != 2000 {
		t.Fatalf("expected clipping loss 2000, got %v", loss)
	}
	loss = Clipping(8000, 10000)
	if loss != 0 {
		t.Fatalf("no clipping expected, got %v", loss)
	}
}

func TestEUEfficiency(t *testing.T) {
	spec := DefaultSpec()
	eu := EUEfficiency(spec)
	if eu < 0.8 || eu > spec.EtaPeak {
		t.Fatalf("EU efficiency out of range: %v", eu)
	}
}

func TestCECEfficiency(t *testing.T) {
	spec := DefaultSpec()
	cec := CECEfficiency(spec)
	if cec < 0.8 || cec > spec.EtaPeak {
		t.Fatalf("CEC efficiency out of range: %v", cec)
	}
}

func TestMPPTRange(t *testing.T) {
	spec := DefaultSpec()
	if !MPPTRange(400, spec) {
		t.Fatalf("400V should be in MPPT range")
	}
	if MPPTRange(100, spec) {
		t.Fatalf("100V should not be in MPPT range")
	}
}

func TestDCACRatio(t *testing.T) {
	ratio := DCACRatio(12000, 10000)
	if math.Abs(ratio-1.2) > 1e-10 {
		t.Fatalf("expected 1.2, got %v", ratio)
	}
}

func TestTempDerating(t *testing.T) {
	d := TempDerating(50, 45, 0.02)
	expected := 1 - 0.02*5
	if math.Abs(d-expected) > 1e-10 {
		t.Fatalf("expected %v, got %v", expected, d)
	}
}

func TestTempDerating_BelowRated(t *testing.T) {
	d := TempDerating(40, 45, 0.02)
	if d != 1.0 {
		t.Fatalf("below rated temp should return 1.0, got %v", d)
	}
}

func TestStringSizing(t *testing.T) {
	min, max := StringSizing(40, 200, 600, -0.003, -10, 40)
	if min < 1 || max < min {
		t.Fatalf("invalid string sizing: min=%d, max=%d", min, max)
	}
}

func TestSandiaModel_Positive(t *testing.T) {
	spec := DefaultSpec()
	pac := SandiaModel(5000, 400, spec)
	if pac <= 0 {
		t.Fatalf("expected positive AC power, got %v", pac)
	}
}

func TestAnnualEnergy(t *testing.T) {
	ae := AnnualEnergy(5.0, 10, 0.96)
	expected := 5.0 * 10 * 0.96 * 365
	if math.Abs(ae-expected) > 1e-10 {
		t.Fatalf("expected %v, got %v", expected, ae)
	}
}
