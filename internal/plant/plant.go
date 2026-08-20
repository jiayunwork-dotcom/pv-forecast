// Package plant parses PV plant readings from CSV and groups them per inverter.
package plant

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

// Reading is a single timestamped plant measurement.
type Reading struct {
	Timestamp  string
	Irradiance float64
	Temp       float64
	Wind       float64
	Power      float64
	Inverter   string
}

// ParseReadings reads a CSV file with the header
// ts,irradiance,temp,wind,power,inverter and returns the parsed readings.
// It returns an error if the file is missing, unreadable, or contains a
// malformed (non-numeric) row. The first row is treated as the header.
func ParseReadings(path string) ([]Reading, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open readings: %w", err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = 6
	r.TrimLeadingSpace = true

	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("read csv: %w", err)
	}
	if len(records) < 2 {
		return nil, fmt.Errorf("need header + at least one data row, got %d lines", len(records))
	}
	if len(records[0]) != 6 {
		return nil, fmt.Errorf("header: expected 6 columns, got %d", len(records[0]))
	}

	readings := make([]Reading, 0, len(records)-1)
	for i, rec := range records[1:] {
		if len(rec) != 6 {
			return nil, fmt.Errorf("row %d: expected 6 fields, got %d", i+2, len(rec))
		}
		irr, err := strconv.ParseFloat(rec[1], 64)
		if err != nil {
			return nil, fmt.Errorf("row %d: bad irradiance %q: %w", i+2, rec[1], err)
		}
		temp, err := strconv.ParseFloat(rec[2], 64)
		if err != nil {
			return nil, fmt.Errorf("row %d: bad temp %q: %w", i+2, rec[2], err)
		}
		wind, err := strconv.ParseFloat(rec[3], 64)
		if err != nil {
			return nil, fmt.Errorf("row %d: bad wind %q: %w", i+2, rec[3], err)
		}
		power, err := strconv.ParseFloat(rec[4], 64)
		if err != nil {
			return nil, fmt.Errorf("row %d: bad power %q: %w", i+2, rec[4], err)
		}
		readings = append(readings, Reading{
			Timestamp:  rec[0],
			Irradiance: irr,
			Temp:       temp,
			Wind:       wind,
			Power:      power,
			Inverter:   rec[5],
		})
	}
	return readings, nil
}

// ByInverter groups readings by their inverter id. It is nil-safe: a nil or
// empty input yields a non-nil empty map.
func ByInverter(readings []Reading) map[string][]Reading {
	if readings == nil {
		return map[string][]Reading{}
	}
	return bindGroup(readings)
}
