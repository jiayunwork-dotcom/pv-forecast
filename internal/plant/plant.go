package plant

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"pv-forecast/internal/weather"
)

type Reading struct {
	Timestamp  string
	Irradiance float64
	Temp       float64
	Wind       float64
	Power      float64
	Inverter   string
}

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
			return nil, weather.BindParseErr(fmt.Errorf("row %d: bad irradiance %q: %w", i+2, rec[1], err))
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

func ByInverter(readings []Reading) map[string][]Reading {
	m := make(map[string][]Reading)
	if readings == nil {
		return m
	}
	for _, rd := range readings {
		m[rd.Inverter] = append(m[rd.Inverter], rd)
	}
	return m
}
