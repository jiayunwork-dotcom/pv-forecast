package solar

import (
	"math"
	"testing"
)

func TestTiltedIrradiance_Flat(t *testing.T) {
	poa := TiltedIrradiance(800, 600, 200, 0, 180, 30, 180)
	if poa < 600 || poa > 900 {
		t.Fatalf("flat panel expected 600-900, got %v", poa)
	}
}

func TestAOI_Normal(t *testing.T) {
	aoi := AOI(30, 180, 30, 180)
	if math.Abs(aoi) > 1 {
		t.Fatalf("expected AOI ~0, got %v", aoi)
	}
}

func TestAOI_Perpendicular(t *testing.T) {
	aoi := AOI(0, 180, 90, 0)
	if math.Abs(aoi-90) > 1 {
		t.Fatalf("expected AOI ~90, got %v", aoi)
	}
}

func TestIAMLoss_Normal(t *testing.T) {
	iam := IAMLoss(0, 0.05)
	if math.Abs(iam-1.0) > 1e-10 {
		t.Fatalf("normal incidence IAM should be 1, got %v", iam)
	}
}

func TestIAMLoss_Grazing(t *testing.T) {
	iam := IAMLoss(90, 0.05)
	if iam != 0 {
		t.Fatalf("grazing angle IAM should be 0, got %v", iam)
	}
}

func TestOptimalTilt(t *testing.T) {
	tilt := OptimalTilt(40)
	if tilt < 25 || tilt > 40 {
		t.Fatalf("optimal tilt for lat=40 expected 25-40, got %v", tilt)
	}
}

func TestPerez(t *testing.T) {
	result := Perez(200, 600, 30, 25, 10)
	if result <= 0 {
		t.Fatalf("Perez diffuse should be positive, got %v", result)
	}
}

func TestTranspositionFactor(t *testing.T) {
	tf := TranspositionFactor(900, 800)
	if tf < 1.0 || tf > 1.5 {
		t.Fatalf("transposition factor expected >1, got %v", tf)
	}
}

func TestDiffuseIsotropic(t *testing.T) {
	d := DiffuseIsotropic(200, 30)
	if d <= 0 || d > 200 {
		t.Fatalf("isotropic diffuse out of range: %v", d)
	}
}

func TestGroundReflected(t *testing.T) {
	g := GroundReflected(800, 30, 0.2)
	if g <= 0 || g > 160 {
		t.Fatalf("ground reflected out of range: %v", g)
	}
}

func TestEffectiveIrradiance(t *testing.T) {
	e := EffectiveIrradiance(900, 0.98, 1.0)
	if math.Abs(e-882) > 1 {
		t.Fatalf("expected ~882, got %v", e)
	}
}

func TestSpectralFactor(t *testing.T) {
	sf := SpectralFactor(1.5)
	if sf < 0.9 || sf > 1.1 {
		t.Fatalf("spectral factor expected 0.9-1.1, got %v", sf)
	}
}
