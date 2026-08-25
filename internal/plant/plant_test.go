package plant

import (
	"os"
	"testing"
)

func TestParseReadings(t *testing.T) {
	bad := t.TempDir() + "/bad.csv"
	if err := os.WriteFile(bad, []byte("ts,irradiance,temp,wind,power,inverter\n2026-01-01,notanumber,30,3,100,INV-A\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ParseReadings(bad); err == nil {
		t.Fatalf("expected error for malformed row, got nil")
	}

	good := t.TempDir() + "/good.csv"
	if err := os.WriteFile(good, []byte("ts,irradiance,temp,wind,power,inverter\n2026-01-01,800,30,3,131,INV-A\n2026-01-01,850,30,3,136,INV-A\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	rs, err := ParseReadings(good)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rs) != 2 {
		t.Fatalf("expected 2 readings, got %d", len(rs))
	}
	if rs[0].Irradiance != 800 || rs[1].Power != 136 {
		t.Fatalf("unexpected parsed values: %+v", rs)
	}

	if _, err := ParseReadings("/nonexistent/path/to/file.csv"); err == nil {
		t.Fatalf("expected error for missing file, got nil")
	}
}

func TestByInverter(t *testing.T) {
	if m := ByInverter(nil); m == nil || len(m) != 0 {
		t.Fatalf("expected non-nil empty map for nil input")
	}

	rs := []Reading{
		{Timestamp: "t1", Irradiance: 800, Temp: 30, Wind: 3, Power: 131, Inverter: "INV-A"},
		{Timestamp: "t2", Irradiance: 850, Temp: 30, Wind: 3, Power: 136, Inverter: "INV-A"},
		{Timestamp: "t3", Irradiance: 820, Temp: 30, Wind: 3, Power: 90, Inverter: "INV-B"},
	}
	m := ByInverter(rs)
	if len(m["INV-A"]) != 2 || len(m["INV-B"]) != 1 {
		t.Fatalf("unexpected grouping: A=%d B=%d", len(m["INV-A"]), len(m["INV-B"]))
	}
}
