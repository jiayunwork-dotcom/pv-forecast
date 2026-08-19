package monitoring

import (
	"math"
	"testing"
)

func TestCheckPR_Normal(t *testing.T) {
	a := CheckPR(0.85, 0.75)
	if a != nil {
		t.Fatalf("normal PR should not alarm")
	}
}

func TestCheckPR_Low(t *testing.T) {
	a := CheckPR(0.60, 0.75)
	if a == nil || a.Code != "LOW_PR" {
		t.Fatalf("low PR should trigger alarm")
	}
}

func TestCheckInverterOutput_OK(t *testing.T) {
	a := CheckInverterOutput(9800, 10000, 0.05)
	if a != nil {
		t.Fatalf("within tolerance should not alarm")
	}
}

func TestCheckInverterOutput_Deviation(t *testing.T) {
	a := CheckInverterOutput(8000, 10000, 0.05)
	if a == nil || a.Code != "INV_DEVIATION" {
		t.Fatalf("deviation should trigger alarm")
	}
}

func TestCheckStringCurrent(t *testing.T) {
	currents := []float64{10, 10, 10, 5, 10}
	alarms := CheckStringCurrent(currents, 20)
	if len(alarms) != 1 {
		t.Fatalf("expected 1 alarm, got %d", len(alarms))
	}
}

func TestCheckCommunication(t *testing.T) {
	a := CheckCommunication(60, 30)
	if a == nil || a.Level != Critical {
		t.Fatalf("should be critical alarm")
	}
	a = CheckCommunication(10, 30)
	if a != nil {
		t.Fatalf("should not alarm within threshold")
	}
}

func TestAvailabilityRate(t *testing.T) {
	rate := AvailabilityRate(8760, 100)
	if math.Abs(rate-0.9886) > 0.001 {
		t.Fatalf("expected ~0.989, got %v", rate)
	}
}

func TestEnergyBalance_OK(t *testing.T) {
	a := EnergyBalance(9900, 10000, 2)
	if a != nil {
		t.Fatalf("should not alarm within tolerance")
	}
}

func TestEnergyBalance_Alarm(t *testing.T) {
	a := EnergyBalance(9500, 10000, 2)
	if a == nil {
		t.Fatalf("should alarm for 5%% deviation")
	}
}
