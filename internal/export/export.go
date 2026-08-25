package export

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Record struct {
	Timestamp string
	Predicted float64
	Actual    float64
	Error     float64
	Inverter  string
}

type Report struct {
	Title      string
	Plant      string
	Period     string
	TotalGen   float64
	PR         float64
	MBE        float64
	RMSE       float64
	LowStrings []string
	Losses     map[string]float64
	Records    []Record
}

func WriteCSV(w io.Writer, records []Record) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()
	if err := cw.Write([]string{"timestamp", "predicted", "actual", "error", "inverter"}); err != nil {
		return fmt.Errorf("write header: %w", err)
	}
	for _, r := range records {
		row := []string{
			r.Timestamp,
			strconv.FormatFloat(r.Predicted, 'f', 4, 64),
			strconv.FormatFloat(r.Actual, 'f', 4, 64),
			strconv.FormatFloat(r.Error, 'f', 4, 64),
			r.Inverter,
		}
		if err := cw.Write(row); err != nil {
			return fmt.Errorf("write row: %w", err)
		}
	}
	return nil
}

func WriteJSON(w io.Writer, report Report) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}

func ReadCSV(r io.Reader) ([]Record, error) {
	cr := csv.NewReader(r)
	rows, err := cr.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("read csv: %w", err)
	}
	if len(rows) < 2 {
		return nil, nil
	}
	records := make([]Record, 0, len(rows)-1)
	for i, row := range rows[1:] {
		if len(row) < 5 {
			return nil, fmt.Errorf("row %d: too few fields", i+2)
		}
		pred, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			return nil, fmt.Errorf("row %d: bad predicted: %w", i+2, err)
		}
		act, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("row %d: bad actual: %w", i+2, err)
		}
		errV, err := strconv.ParseFloat(row[3], 64)
		if err != nil {
			return nil, fmt.Errorf("row %d: bad error: %w", i+2, err)
		}
		records = append(records, Record{
			Timestamp: row[0],
			Predicted: pred,
			Actual:    act,
			Error:     errV,
			Inverter:  row[4],
		})
	}
	return records, nil
}

func FormatSummary(report Report) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("=== %s ===\n", report.Title))
	sb.WriteString(fmt.Sprintf("Plant: %s\n", report.Plant))
	sb.WriteString(fmt.Sprintf("Period: %s\n", report.Period))
	sb.WriteString(fmt.Sprintf("Total Generation: %.2f kWh\n", report.TotalGen))
	sb.WriteString(fmt.Sprintf("Performance Ratio: %.4f\n", report.PR))
	sb.WriteString(fmt.Sprintf("MBE: %.4f, RMSE: %.4f\n", report.MBE, report.RMSE))
	if len(report.LowStrings) > 0 {
		sb.WriteString(fmt.Sprintf("Low strings: %s\n", strings.Join(report.LowStrings, ", ")))
	}
	if len(report.Losses) > 0 {
		sb.WriteString("Losses:\n")
		for k, v := range report.Losses {
			sb.WriteString(fmt.Sprintf("  %-15s %.4f\n", k, v))
		}
	}
	return sb.String()
}

func CompareReports(a, b Report) map[string]float64 {
	diff := make(map[string]float64)
	diff["total_gen_diff"] = b.TotalGen - a.TotalGen
	diff["pr_diff"] = b.PR - a.PR
	diff["rmse_diff"] = b.RMSE - a.RMSE
	return diff
}

func FilterByInverter(records []Record, inverter string) []Record {
	var filtered []Record
	for _, r := range records {
		if r.Inverter == inverter {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

func GroupByTimestamp(records []Record) map[string][]Record {
	groups := make(map[string][]Record)
	for _, r := range records {
		groups[r.Timestamp] = append(groups[r.Timestamp], r)
	}
	return groups
}
