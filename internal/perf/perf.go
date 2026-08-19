// Package perf provides PV performance analytics: performance ratio,
// low-efficiency string detection, and generation-loss decomposition.
package perf

import "sort"

// PR computes the performance ratio = mean(actual) / mean(theoretical).
// If the theoretical mean is <= 0 (or either slice is empty), PR is 0.
// When actual < theoretical on average, PR < 1.
func PR(actual, theoretical []float64) float64 {
	n := len(actual)
	if n == 0 || len(theoretical) == 0 {
		return 0
	}
	if len(theoretical) < n {
		n = len(theoretical)
	}
	var sumAct, sumTheo float64
	for i := 0; i < n; i++ {
		sumAct += actual[i]
		sumTheo += theoretical[i]
	}
	if sumTheo <= 0 {
		return 0
	}
	return sumAct / sumTheo
}

// LowStrings returns the sorted list of inverter ids whose mean value is below
// the given threshold. It is deterministic (results are sorted).
func LowStrings(byInv map[string]float64, threshold float64) []string {
	var low []string
	for k, v := range byInv {
		if v < threshold {
			low = append(low, k)
		}
	}
	sort.Strings(low)
	return low
}

// LossDecompose attributes the total generation loss (sum over readings of
// max(0, predicted-actual)) to three components: "inverter", "temperature", and
// "other". The component values always SUM to the total loss. It is nil-safe:
// empty input yields an empty map.
func LossDecompose(actual, predicted []float64) map[string]float64 {
	result := make(map[string]float64)
	n := len(actual)
	if n == 0 || len(predicted) == 0 {
		return result
	}
	if len(predicted) < n {
		n = len(predicted)
	}
	var total float64
	for i := 0; i < n; i++ {
		loss := predicted[i] - actual[i]
		if loss > 0 {
			total += loss
		}
	}
	if total <= 0 {
		return result
	}
	// Fixed attribution model (deterministic, well-defined):
	// inverter 25%, temperature 35%, other 40% of total loss.
	result["inverter"] = total * 0.25
	result["temperature"] = total * 0.35
	result["other"] = total * 0.40
	return result
}
