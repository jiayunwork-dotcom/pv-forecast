package metric

import (
	"math"
	"testing"
)

func TestMAE_Perfect(t *testing.T) {
	a := []float64{1, 2, 3}
	if MAE(a, a) != 0 {
		t.Fatalf("perfect prediction MAE should be 0")
	}
}

func TestMAE_Known(t *testing.T) {
	actual := []float64{1, 2, 3}
	pred := []float64{2, 2, 4}
	mae := MAE(actual, pred)
	if math.Abs(mae-2.0/3.0) > 1e-10 {
		t.Fatalf("expected 2/3, got %v", mae)
	}
}

func TestRMSE_Perfect(t *testing.T) {
	a := []float64{1, 2, 3}
	if RMSE(a, a) != 0 {
		t.Fatalf("perfect prediction RMSE should be 0")
	}
}

func TestRMSE_Known(t *testing.T) {
	actual := []float64{1, 2, 3}
	pred := []float64{2, 3, 4}
	rmse := RMSE(actual, pred)
	if math.Abs(rmse-1.0) > 1e-10 {
		t.Fatalf("expected 1, got %v", rmse)
	}
}

func TestMAPE_Known(t *testing.T) {
	actual := []float64{10, 20, 30}
	pred := []float64{9, 18, 27}
	mape := MAPE(actual, pred)
	// (1/10 + 2/20 + 3/30)/3 * 100 = 10%
	if math.Abs(mape-10.0) > 1e-10 {
		t.Fatalf("expected 10%%, got %v", mape)
	}
}

func TestR2_Perfect(t *testing.T) {
	a := []float64{1, 2, 3, 4, 5}
	r2 := R2(a, a)
	if math.Abs(r2-1.0) > 1e-10 {
		t.Fatalf("perfect R2 should be 1, got %v", r2)
	}
}

func TestR2_Bad(t *testing.T) {
	actual := []float64{1, 2, 3, 4, 5}
	pred := []float64{5, 4, 3, 2, 1}
	r2 := R2(actual, pred)
	if r2 > 0 {
		t.Fatalf("reversed prediction R2 should be <= 0, got %v", r2)
	}
}

func TestMBE_Positive(t *testing.T) {
	actual := []float64{1, 2, 3}
	pred := []float64{2, 3, 4}
	mbe := MBE(actual, pred)
	if math.Abs(mbe-1.0) > 1e-10 {
		t.Fatalf("expected MBE=1, got %v", mbe)
	}
}

func TestNRMSE(t *testing.T) {
	actual := []float64{0, 10}
	pred := []float64{1, 9}
	nrmse := NRMSE(actual, pred)
	rmse := RMSE(actual, pred)
	if math.Abs(nrmse-rmse/10) > 1e-10 {
		t.Fatalf("expected %v, got %v", rmse/10, nrmse)
	}
}

func TestSkill_BetterThanPersistence(t *testing.T) {
	actual := []float64{1, 3, 5, 7, 9}
	model := []float64{1.1, 2.9, 5.1, 6.9, 9.1}
	persist := []float64{1, 1, 3, 5, 7}
	skill := Skill(actual, model, persist)
	if skill <= 0 {
		t.Fatalf("good model should have positive skill, got %v", skill)
	}
}

func TestCorrelation_Perfect(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5}
	y := []float64{2, 4, 6, 8, 10}
	c := Correlation(x, y)
	if math.Abs(c-1.0) > 1e-10 {
		t.Fatalf("perfect linear correlation should be 1, got %v", c)
	}
}

func TestCorrelation_Negative(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5}
	y := []float64{10, 8, 6, 4, 2}
	c := Correlation(x, y)
	if math.Abs(c-(-1.0)) > 1e-10 {
		t.Fatalf("perfect negative correlation should be -1, got %v", c)
	}
}

func TestKSI(t *testing.T) {
	actual := []float64{1, 2, 3, 4, 5}
	pred := []float64{1, 2, 3, 4, 5}
	ksi := KSI(actual, pred, 10)
	if ksi != 0 {
		t.Fatalf("identical distributions KSI should be 0, got %v", ksi)
	}
}
