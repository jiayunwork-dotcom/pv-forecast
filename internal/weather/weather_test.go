package weather

import (
	"math"
	"testing"
)

func TestDayOfYear(t *testing.T) {
	if DayOfYear(1, 1) != 1 {
		t.Fatalf("Jan 1 should be day 1")
	}
	if DayOfYear(12, 31) != 365 {
		t.Fatalf("Dec 31 should be day 365, got %d", DayOfYear(12, 31))
	}
}

func TestDeclination_Solstice(t *testing.T) {
	dec := Declination(172) * 180 / math.Pi
	if math.Abs(dec-23.45) > 1.0 {
		t.Fatalf("summer solstice declination expected ~23.45, got %v", dec)
	}
}

func TestSolarPosition_Noon(t *testing.T) {
	pos := SolarPosition(0, 80, 12)
	elevDeg := pos.Elev * 180 / math.Pi
	if elevDeg < 85 {
		t.Fatalf("equator equinox noon should have high elevation, got %v", elevDeg)
	}
}

func TestAirMass_Zenith0(t *testing.T) {
	am := AirMass(0)
	if math.Abs(am-1.0) > 0.01 {
		t.Fatalf("zenith=0 should give AM=1, got %v", am)
	}
}

func TestAirMass_Horizon(t *testing.T) {
	am := AirMass(90)
	if am < 30 {
		t.Fatalf("zenith=90 should give high AM, got %v", am)
	}
}

func TestClearSkyGHI_Noon(t *testing.T) {
	ghi := ClearSkyGHI(0)
	if ghi < 900 || ghi > 1200 {
		t.Fatalf("clear sky GHI at noon expected 900-1200, got %v", ghi)
	}
}

func TestClearSkyGHI_Night(t *testing.T) {
	ghi := ClearSkyGHI(95)
	if ghi != 0 {
		t.Fatalf("expected 0 at night, got %v", ghi)
	}
}

func TestDecomposeGHI(t *testing.T) {
	dni, dhi := DecomposeGHI(500, 1000)
	if dni <= 0 || dhi <= 0 {
		t.Fatalf("both DNI and DHI should be positive, got %v, %v", dni, dhi)
	}
	// DHI + beam component ≈ GHI
	if dhi > 500 {
		t.Fatalf("DHI should not exceed GHI")
	}
}

func TestExtraterrestrialIrradiance(t *testing.T) {
	e := ExtraterrestrialIrradiance(1, 0)
	if e < 1300 || e > 1500 {
		t.Fatalf("extraterrestrial at zenith=0 expected 1300-1500, got %v", e)
	}
}

func TestSunriseSunset_Equator(t *testing.T) {
	rise, set := SunriseSunset(0, 80)
	if math.Abs(rise-6) > 0.5 || math.Abs(set-18) > 0.5 {
		t.Fatalf("equator equinox: expected ~6-18, got %v-%v", rise, set)
	}
}

func TestDaylightHours(t *testing.T) {
	hours := DaylightHours(0, 80)
	if math.Abs(hours-12) > 1 {
		t.Fatalf("equator should have ~12h daylight, got %v", hours)
	}
}

func TestEquationOfTime(t *testing.T) {
	eot := EquationOfTime(105) // mid-April
	if math.Abs(eot) > 20 {
		t.Fatalf("EOT should be within ±20 min, got %v", eot)
	}
}

func TestSolarTime(t *testing.T) {
	st := SolarTime(12, 120, 8, 1)
	if st < 10 || st > 14 {
		t.Fatalf("solar time should be near 12, got %v", st)
	}
}
