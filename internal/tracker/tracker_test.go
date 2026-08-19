package tracker

import (
	"math"
	"testing"
)

func TestTrackAngle_Fixed(t *testing.T) {
	cfg := Config{Mode: Fixed}
	angle := TrackAngle(30, 180, cfg)
	if angle != 0 {
		t.Fatalf("fixed mode should return 0, got %v", angle)
	}
}

func TestTrackAngle_SingleH_Noon(t *testing.T) {
	cfg := DefaultSingleAxis()
	angle := TrackAngle(30, 180, cfg)
	if math.Abs(angle) > 10 {
		t.Fatalf("noon angle should be near 0, got %v", angle)
	}
}

func TestTrackAngle_MaxLimit(t *testing.T) {
	cfg := DefaultSingleAxis()
	cfg.MaxAngle = 45
	angle := TrackAngle(80, 90, cfg)
	if math.Abs(angle) > 45 {
		t.Fatalf("angle should not exceed max, got %v", angle)
	}
}

func TestDualAxisAngles(t *testing.T) {
	tilt, az := DualAxisAngles(30, 200)
	if tilt != 30 || az != 200 {
		t.Fatalf("dual axis should track exactly: %v, %v", tilt, az)
	}
}

func TestBacktrack(t *testing.T) {
	bt := backtrack(70, 0.35)
	maxNoShade := math.Acos(0.35) * 180 / math.Pi
	if bt > maxNoShade+0.01 {
		t.Fatalf("backtrack should limit angle: %v > %v", bt, maxNoShade)
	}
}

func TestShadingLoss_HighSun(t *testing.T) {
	loss := ShadingLoss(0.35, 60, 30)
	if loss > 0.3 {
		t.Fatalf("high sun should have low shading loss, got %v", loss)
	}
}

func TestShadingLoss_Night(t *testing.T) {
	loss := ShadingLoss(0.35, 0, 30)
	if loss != 1.0 {
		t.Fatalf("night should have full loss, got %v", loss)
	}
}

func TestEnergyGain_SingleAxis(t *testing.T) {
	gain := EnergyGain(30, SingleH)
	if gain < 10 || gain > 25 {
		t.Fatalf("single axis gain expected 10-25%%, got %v", gain)
	}
}

func TestEnergyGain_DualAxis(t *testing.T) {
	gain := EnergyGain(30, DualAxis)
	if gain < 20 || gain > 40 {
		t.Fatalf("dual axis gain expected 20-40%%, got %v", gain)
	}
}

func TestRowSpacing(t *testing.T) {
	spacing := RowSpacing(2.0, 25, 15)
	if spacing < 1 || spacing > 10 {
		t.Fatalf("row spacing out of range: %v", spacing)
	}
}

func TestAOI_Tracker(t *testing.T) {
	aoi := AOI(20, 30, 180, 0, 0)
	if aoi < 0 || aoi > 90 {
		t.Fatalf("AOI out of range: %v", aoi)
	}
}
