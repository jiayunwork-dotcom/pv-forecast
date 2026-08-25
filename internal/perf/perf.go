package perf

import "sort"

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
	result["inverter"] = total * 0.25
	result["temperature"] = total * 0.35
	result["other"] = total * 0.40
	return result
}
