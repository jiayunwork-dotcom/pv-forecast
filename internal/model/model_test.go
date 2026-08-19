package model

import "testing"

// TestDCModel verifies that higher irradiance raises DC power and that higher
// temperature lowers DC power (via the negative temperature coefficient), and
// that DC is clamped at zero.
func TestDCModel(t *testing.T) {
	// Higher irradiance -> higher DC.
	dcLow := DCPower(500, 1, 0.17, 25, -0.004)
	dcHigh := DCPower(800, 1, 0.17, 25, -0.004)
	if dcHigh <= dcLow {
		t.Fatalf("expected higher irradiance -> higher DC (low=%.4f high=%.4f)", dcLow, dcHigh)
	}

	// Higher temperature -> lower DC (negative coeff).
	dcCool := DCPower(800, 1, 0.17, 15, -0.004)
	dcHot := DCPower(800, 1, 0.17, 45, -0.004)
	if dcHot >= dcCool {
		t.Fatalf("expected higher temp -> lower DC (cool=%.4f hot=%.4f)", dcCool, dcHot)
	}

	// DC must be clamped >= 0 even for absurd temperatures.
	if DCPower(100, 1, 0.17, 100000, -0.004) < 0 {
		t.Fatalf("DC must never be negative")
	}
}

// TestACModel verifies the AC = DC * invEff relationship and clamping.
func TestACModel(t *testing.T) {
	if got := ACPower(100, 0.96); got < 95.9 || got > 96.1 {
		t.Fatalf("expected AC=96 for DC=100, invEff=0.96, got %.4f", got)
	}
	if ACPower(100, -1) != 0 {
		t.Fatalf("expected AC clamped to 0 for negative invEff")
	}
}
