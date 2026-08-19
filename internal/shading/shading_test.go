package shading

import (
	"math"
	"testing"
)

func makeProfile() HorizonProfile {
	return HorizonProfile{
		Azimuths:   []float64{0, 90, 180, 270, 360},
		Elevations: []float64{5, 10, 0, 15, 5},
	}
}

func TestIsShaded_True(t *testing.T) {
	p := makeProfile()
	if !IsShaded(90, 5, p) {
		t.Fatalf("should be shaded at az=90, elev=5 (horizon=10)")
	}
}

func TestIsShaded_False(t *testing.T) {
	p := makeProfile()
	if IsShaded(180, 5, p) {
		t.Fatalf("should not be shaded at az=180, elev=5 (horizon=0)")
	}
}

func TestShadingFactor_Full(t *testing.T) {
	p := makeProfile()
	f := ShadingFactor(180, 30, p)
	if f != 1.0 {
		t.Fatalf("no shading expected, got %v", f)
	}
}

func TestShadingFactor_Partial(t *testing.T) {
	p := makeProfile()
	f := ShadingFactor(90, 5, p)
	if f >= 1 || f <= 0 {
		t.Fatalf("partial shading expected, got %v", f)
	}
}

func TestDailyShadingLoss(t *testing.T) {
	az := []float64{90, 135, 180, 225, 270}
	elev := []float64{10, 20, 30, 20, 10}
	p := makeProfile()
	loss := DailyShadingLoss(az, elev, p)
	if loss < 0 || loss > 1 {
		t.Fatalf("loss out of range: %v", loss)
	}
}

func TestNearShadingLoss_NoShade(t *testing.T) {
	loss := NearShadingLoss(5, 20, 45)
	if loss > 0.01 {
		t.Fatalf("far object should not shade: %v", loss)
	}
}

func TestNearShadingLoss_Close(t *testing.T) {
	loss := NearShadingLoss(10, 3, 20)
	if loss <= 0 {
		t.Fatalf("close tall object should shade: %v", loss)
	}
}

func TestSelfShadingLoss_NoShade(t *testing.T) {
	loss := SelfShadingLoss(5, 2, 25, 60)
	if loss > 0.01 {
		t.Fatalf("high sun should not self-shade: %v", loss)
	}
}

func TestSelfShadingLoss_LowSun(t *testing.T) {
	loss := SelfShadingLoss(2, 2, 30, 10)
	if loss <= 0 {
		t.Fatalf("low sun should self-shade: %v", loss)
	}
}

func TestObstacleToProfile(t *testing.T) {
	obs := []Obstacle{
		{AzStart: 80, AzEnd: 100, Elev: 20},
	}
	p := ObstacleToProfile(obs, 10)
	if len(p.Azimuths) != 36 {
		t.Fatalf("expected 36 steps, got %d", len(p.Azimuths))
	}
	found := false
	for i, az := range p.Azimuths {
		if math.Abs(az-90) < 1 && p.Elevations[i] == 20 {
			found = true
		}
	}
	if !found {
		t.Fatalf("obstacle at az=90 not reflected in profile")
	}
}

func TestAnnualShadingLoss(t *testing.T) {
	p := makeProfile()
	az := []float64{90, 180, 270}
	elev := []float64{5, 30, 10}
	ghi := []float64{400, 800, 400}
	loss := AnnualShadingLoss(az, elev, ghi, p)
	if loss < 0 || loss > 1 {
		t.Fatalf("annual loss out of range: %v", loss)
	}
}
