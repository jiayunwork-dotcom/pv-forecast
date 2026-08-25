package export

import (
	"bytes"
	"strings"
	"testing"
)

func sampleRecords() []Record {
	return []Record{
		{Timestamp: "2024-01-01T10:00", Predicted: 500, Actual: 480, Error: 20, Inverter: "INV1"},
		{Timestamp: "2024-01-01T11:00", Predicted: 700, Actual: 690, Error: 10, Inverter: "INV1"},
		{Timestamp: "2024-01-01T10:00", Predicted: 450, Actual: 440, Error: 10, Inverter: "INV2"},
	}
}

func TestWriteCSV_ReadCSV_Roundtrip(t *testing.T) {
	records := sampleRecords()
	var buf bytes.Buffer
	if err := WriteCSV(&buf, records); err != nil {
		t.Fatalf("write csv failed: %v", err)
	}
	loaded, err := ReadCSV(&buf)
	if err != nil {
		t.Fatalf("read csv failed: %v", err)
	}
	if len(loaded) != len(records) {
		t.Fatalf("expected %d records, got %d", len(records), len(loaded))
	}
	for i := range records {
		if loaded[i].Timestamp != records[i].Timestamp {
			t.Fatalf("timestamp mismatch at %d", i)
		}
	}
}

func TestWriteJSON(t *testing.T) {
	report := Report{
		Title:   "Test Report",
		Plant:   "SolarFarm1",
		PR:      0.85,
		Losses:  map[string]float64{"inverter": 10, "soiling": 5},
		Records: sampleRecords(),
	}
	var buf bytes.Buffer
	if err := WriteJSON(&buf, report); err != nil {
		t.Fatalf("write json failed: %v", err)
	}
	if !strings.Contains(buf.String(), "SolarFarm1") {
		t.Fatalf("json should contain plant name")
	}
}

func TestFormatSummary(t *testing.T) {
	report := Report{
		Title:      "Monthly Report",
		Plant:      "Site A",
		Period:     "2024-01",
		TotalGen:   5000,
		PR:         0.82,
		MBE:        -0.5,
		RMSE:       2.3,
		LowStrings: []string{"INV3", "INV5"},
		Losses:     map[string]float64{"inverter": 200},
	}
	s := FormatSummary(report)
	if !strings.Contains(s, "Monthly Report") {
		t.Fatalf("missing title in summary")
	}
	if !strings.Contains(s, "INV3") {
		t.Fatalf("missing low string in summary")
	}
}

func TestFilterByInverter(t *testing.T) {
	records := sampleRecords()
	filtered := FilterByInverter(records, "INV2")
	if len(filtered) != 1 {
		t.Fatalf("expected 1, got %d", len(filtered))
	}
}

func TestGroupByTimestamp(t *testing.T) {
	records := sampleRecords()
	groups := GroupByTimestamp(records)
	if len(groups["2024-01-01T10:00"]) != 2 {
		t.Fatalf("expected 2 records at 10:00")
	}
}

func TestCompareReports(t *testing.T) {
	a := Report{TotalGen: 5000, PR: 0.82, RMSE: 2.3}
	b := Report{TotalGen: 5200, PR: 0.84, RMSE: 2.1}
	diff := CompareReports(a, b)
	if diff["total_gen_diff"] != 200 {
		t.Fatalf("expected 200, got %v", diff["total_gen_diff"])
	}
}
