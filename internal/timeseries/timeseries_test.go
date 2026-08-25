package timeseries

import (
	"math"
	"testing"
)

func TestMovingAvg(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5}
	ma := MovingAvg(data, 3)
	if math.Abs(ma[2]-2) > 1e-10 {
		t.Fatalf("expected 2, got %v", ma[2])
	}
	if math.Abs(ma[4]-4) > 1e-10 {
		t.Fatalf("expected 4, got %v", ma[4])
	}
}

func TestExpSmooth(t *testing.T) {
	data := []float64{10, 12, 14, 16}
	es := ExpSmooth(data, 0.5)
	if es[0] != 10 {
		t.Fatalf("first value should be unchanged")
	}
	if math.Abs(es[1]-11) > 1e-10 {
		t.Fatalf("expected 11, got %v", es[1])
	}
}

func TestDiff(t *testing.T) {
	data := []float64{1, 3, 6, 10}
	d := Diff(data)
	if len(d) != 3 {
		t.Fatalf("expected 3, got %d", len(d))
	}
	if d[0] != 2 || d[1] != 3 || d[2] != 4 {
		t.Fatalf("diff wrong: %v", d)
	}
}

func TestCumSum(t *testing.T) {
	data := []float64{1, 2, 3, 4}
	cs := CumSum(data)
	if cs[3] != 10 {
		t.Fatalf("expected 10, got %v", cs[3])
	}
}

func TestResample(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5, 6}
	r := Resample(data, 3)
	if len(r) != 2 {
		t.Fatalf("expected 2, got %d", len(r))
	}
	if math.Abs(r[0]-2) > 1e-10 {
		t.Fatalf("expected 2, got %v", r[0])
	}
}

func TestInterpolate(t *testing.T) {
	data := []float64{1, math.NaN(), 3, math.NaN(), 5}
	interp := Interpolate(data)
	if math.Abs(interp[1]-2) > 1e-10 {
		t.Fatalf("expected interpolated 2, got %v", interp[1])
	}
	if math.Abs(interp[3]-4) > 1e-10 {
		t.Fatalf("expected interpolated 4, got %v", interp[3])
	}
}

func TestPercentile(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	p50 := Percentile(data, 50)
	if math.Abs(p50-5.5) > 0.5 {
		t.Fatalf("P50 expected ~5.5, got %v", p50)
	}
}

func TestAnomalies(t *testing.T) {
	data := []float64{1, 1, 1, 1, 100, 1, 1}
	indices := Anomalies(data, 2)
	if len(indices) != 1 || indices[0] != 4 {
		t.Fatalf("expected anomaly at index 4, got %v", indices)
	}
}

func TestRampRate(t *testing.T) {
	data := []float64{100, 200, 150}
	rates := RampRate(data, 5.0)
	if len(rates) != 2 {
		t.Fatalf("expected 2 rates, got %d", len(rates))
	}
	if math.Abs(rates[0]-20) > 1e-10 {
		t.Fatalf("expected rate 20, got %v", rates[0])
	}
}

func TestPeakHours(t *testing.T) {
	irr := []float64{500, 800, 1000, 800, 500}
	ph := PeakHours(irr, 1.0)
	expected := 3600.0 / 1000
	if math.Abs(ph-expected) > 1e-10 {
		t.Fatalf("expected %v, got %v", expected, ph)
	}
}

func TestAutocorrelation_Lag0(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5}
	ac := Autocorrelation(data, 1)
	if ac < 0.3 {
		t.Fatalf("lag-1 autocorrelation of linear data should be moderate-high, got %v", ac)
	}
}
