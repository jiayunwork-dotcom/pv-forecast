package sizing

import (
	"math"
	"testing"
)

func TestStringSize(t *testing.T) {
	spec := DefaultSpec()
	min, max := StringSize(spec)
	if min < 1 || max < min {
		t.Fatalf("invalid string size: min=%d, max=%d", min, max)
	}
	if max > 15 {
		t.Fatalf("max string too large for 600V limit: %d", max)
	}
}

func TestParallelStrings(t *testing.T) {
	spec := DefaultSpec()
	n := ParallelStrings(spec)
	if n < 1 || n > 5 {
		t.Fatalf("parallel strings expected 1-5, got %d", n)
	}
}

func TestDCACRatio(t *testing.T) {
	spec := DefaultSpec()
	ratio := DCACRatio(10, 3, spec)
	if ratio < 1.0 || ratio > 2.0 {
		t.Fatalf("DCAC ratio expected 1-2, got %v", ratio)
	}
}

func TestTotalCapacityKW(t *testing.T) {
	spec := DefaultSpec()
	cap := TotalCapacityKW(10, 3, spec)
	expected := 10 * 3 * 0.4
	if math.Abs(cap-expected) > 1e-10 {
		t.Fatalf("expected %v kW, got %v", expected, cap)
	}
}

func TestAreaRequired(t *testing.T) {
	area := AreaRequired(10, 3, 2.0, 0.4)
	expected := 60.0 / 0.4
	if math.Abs(area-expected) > 1e-10 {
		t.Fatalf("expected %v, got %v", expected, area)
	}
}

func TestAnnualYield(t *testing.T) {
	yield := AnnualYield(10, 5, 0.8)
	expected := 10 * 5 * 365 * 0.8
	if math.Abs(yield-expected) > 1e-10 {
		t.Fatalf("expected %v, got %v", expected, yield)
	}
}

func TestLCOE(t *testing.T) {
	lcoe := LCOE(100000, 15000, 25, 0.05)
	if lcoe < 0.3 || lcoe > 1.0 {
		t.Fatalf("LCOE out of range: %v", lcoe)
	}
}

func TestPaybackYears(t *testing.T) {
	pb := PaybackYears(100000, 20000)
	if math.Abs(pb-5) > 1e-10 {
		t.Fatalf("expected 5 years, got %v", pb)
	}
}

func TestCableSize(t *testing.T) {
	size := CableSize(20, 50, 2, 400)
	if size < 1 || size > 20 {
		t.Fatalf("cable size out of range: %v mm²", size)
	}
}

func TestShadeFreeSpacing(t *testing.T) {
	spacing := ShadeFreeSpacing(2.0, 25, 15)
	if spacing < 1 || spacing > 10 {
		t.Fatalf("spacing out of range: %v", spacing)
	}
}
