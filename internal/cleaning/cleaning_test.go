package cleaning

import (
	"math"
	"testing"
)

func TestOptimalInterval(t *testing.T) {
	days := OptimalInterval(0.5, 500, 0.1, 100)
	if days < 10 || days > 200 {
		t.Fatalf("interval out of range: %v", days)
	}
}

func TestROI(t *testing.T) {
	roi := CleaningROI(100, 0.5, 90, 200, 0.8)
	// lossRecovered = 0.5/100 * 90 * 200 * 0.8 = 72
	if roi < -0.5 {
		t.Fatalf("ROI too negative, got %v", roi)
	}
}

func TestSoilingRate_Dry(t *testing.T) {
	rate := SoilingRate(0, 35)
	if rate <= 0 {
		t.Fatalf("dry conditions should have positive soiling rate")
	}
}

func TestSoilingRate_Rain(t *testing.T) {
	rate := SoilingRate(10, 25)
	if rate >= SoilingRate(0, 25) {
		t.Fatalf("rain should reduce soiling rate")
	}
}

func TestRainCleaning(t *testing.T) {
	s := RainCleaning(5.0, 3.0)
	if math.Abs(s-2.0) > 1e-10 {
		t.Fatalf("expected 2.0%%, got %v", s)
	}
}

func TestSchedule(t *testing.T) {
	schedule := Schedule(365, 45, 7)
	if len(schedule) < 5 || len(schedule) > 15 {
		t.Fatalf("expected 5-15 cleanings, got %d", len(schedule))
	}
}
