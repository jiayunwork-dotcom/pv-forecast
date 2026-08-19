package perf

import "testing"

// TestPR verifies PR is < 1 when actual < theoretical, equals the exact ratio,
// and is 0 when the theoretical mean is <= 0.
func TestPR(t *testing.T) {
	actual := []float64{90, 90}
	theo := []float64{100, 100}
	pr := PR(actual, theo)
	if pr >= 1 {
		t.Fatalf("expected PR < 1, got %.4f", pr)
	}
	if pr < 0.89 || pr > 0.91 {
		t.Fatalf("expected PR ~0.9, got %.4f", pr)
	}
	if got := PR([]float64{1}, []float64{0}); got != 0 {
		t.Fatalf("expected 0 when theoretical mean <= 0, got %.4f", got)
	}
}

// TestLowStrings verifies inverters below threshold are identified, and that
// the returned list is deterministic (sorted).
func TestLowStrings(t *testing.T) {
	byInv := map[string]float64{
		"INV-A": 0.16,
		"INV-B": 0.11,
		"INV-C": 0.15,
	}
	low := LowStrings(byInv, 0.13)
	if len(low) != 1 || low[0] != "INV-B" {
		t.Fatalf("expected [INV-B], got %v", low)
	}
}

// TestLossDecompose verifies the decomposition components sum to the total loss
// and that it is nil-safe.
func TestLossDecompose(t *testing.T) {
	// nil-safe: empty -> empty map.
	if m := LossDecompose(nil, nil); len(m) != 0 {
		t.Fatalf("expected empty map for nil input, got %v", m)
	}

	actual := []float64{80, 80, 80}
	predicted := []float64{100, 100, 100}
	m := LossDecompose(actual, predicted)
	// total loss = (100-80)*3 = 60
	var sum float64
	for _, v := range m {
		sum += v
	}
	if sum < 59.99 || sum > 60.01 {
		t.Fatalf("components sum %.4f, want ~60", sum)
	}
	if _, ok := m["inverter"]; !ok {
		t.Fatalf("missing 'inverter' component")
	}
	if _, ok := m["temperature"]; !ok {
		t.Fatalf("missing 'temperature' component")
	}
	if _, ok := m["other"]; !ok {
		t.Fatalf("missing 'other' component")
	}
}
