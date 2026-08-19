package celltemp

import (
	"math"
	"testing"
)

func TestCellTempNOCT_NoIrradiance(t *testing.T) {
	tc := CellTempNOCT(25, 0, 1, DefaultNOCT())
	if tc != 25 {
		t.Fatalf("no irradiance should give ambient temp, got %v", tc)
	}
}

func TestCellTempNOCT_Typical(t *testing.T) {
	tc := CellTempNOCT(25, 800, 2, DefaultNOCT())
	if tc < 40 || tc > 65 {
		t.Fatalf("cell temp out of range: %v", tc)
	}
}

func TestCellTempFaiman_Typical(t *testing.T) {
	tc := CellTempFaiman(30, 1000, 2, DefaultFaiman())
	if tc < 50 || tc > 70 {
		t.Fatalf("Faiman cell temp out of range: %v", tc)
	}
}

func TestCellTempSandia_NoIrradiance(t *testing.T) {
	tc := CellTempSandia(20, 0, 3, DefaultSandia())
	if tc != 20 {
		t.Fatalf("no irradiance should give ambient, got %v", tc)
	}
}

func TestCellTempSandia_Typical(t *testing.T) {
	tc := CellTempSandia(25, 800, 1, DefaultSandia())
	if tc < 35 || tc > 65 {
		t.Fatalf("Sandia cell temp out of range: %v", tc)
	}
}

func TestCellTempPVSyst(t *testing.T) {
	tc := CellTempPVSyst(25, 800, 2, 29, 6)
	if tc < 35 || tc > 60 {
		t.Fatalf("PVSyst cell temp out of range: %v", tc)
	}
}

func TestEfficiencyCorrection_STC(t *testing.T) {
	corr := EfficiencyCorrection(25, 25, -0.4)
	if math.Abs(corr-1.0) > 1e-10 {
		t.Fatalf("at STC should be 1.0, got %v", corr)
	}
}

func TestEfficiencyCorrection_Hot(t *testing.T) {
	corr := EfficiencyCorrection(50, 25, -0.4)
	// 1 + (-0.4/100)*25 = 1 - 0.1 = 0.9
	if math.Abs(corr-0.9) > 1e-10 {
		t.Fatalf("expected 0.9, got %v", corr)
	}
}

func TestPowerLoss(t *testing.T) {
	loss := PowerLoss(400, 50, 25, -0.4)
	// factor = 0.9, loss = 400*0.1 = 40
	if math.Abs(loss-40) > 1e-10 {
		t.Fatalf("expected 40W loss, got %v", loss)
	}
}

func TestCellTemp_Dispatch(t *testing.T) {
	tc := CellTemp(NOCT, 30, 800, 2)
	if tc < 40 {
		t.Fatalf("NOCT dispatch should give elevated temp: %v", tc)
	}
	tc2 := CellTemp(Faiman, 30, 800, 2)
	if tc2 < 40 {
		t.Fatalf("Faiman dispatch should give elevated temp: %v", tc2)
	}
}
