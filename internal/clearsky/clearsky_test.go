package clearsky

import (
	"math"
	"testing"
)

func TestIneichen_Noon(t *testing.T) {
	ghi, dni, dhi := Ineichen(20, 0, 3)
	if ghi <= 0 || dni <= 0 || dhi < 0 {
		t.Fatalf("noon should have positive values: ghi=%v, dni=%v, dhi=%v", ghi, dni, dhi)
	}
}

func TestIneichen_Night(t *testing.T) {
	ghi, dni, dhi := Ineichen(95, 0, 3)
	if ghi != 0 || dni != 0 || dhi != 0 {
		t.Fatalf("night should be zero")
	}
}

func TestHaurwitz_Noon(t *testing.T) {
	ghi := Haurwitz(0)
	if ghi < 900 || ghi > 1200 {
		t.Fatalf("Haurwitz at zenith=0 expected 900-1200, got %v", ghi)
	}
}

func TestSimplifiedSolis(t *testing.T) {
	ghi, dni, _ := SimplifiedSolis(30, 100, 0.1)
	if ghi <= 0 || dni <= 0 {
		t.Fatalf("should be positive: ghi=%v, dni=%v", ghi, dni)
	}
}

func TestClearSkyIndex(t *testing.T) {
	kt := ClearSkyIndex(600, 800)
	if math.Abs(kt-0.75) > 1e-10 {
		t.Fatalf("expected 0.75, got %v", kt)
	}
}

func TestCloudFraction(t *testing.T) {
	cf := CloudFraction(0.8)
	if math.Abs(cf-0.2) > 1e-10 {
		t.Fatalf("expected 0.2, got %v", cf)
	}
}
