package metric

import "math"

func MAE(actual, predicted []float64) float64 {
	n := len(actual)
	if n == 0 || len(predicted) == 0 {
		return 0
	}
	if len(predicted) < n {
		n = len(predicted)
	}
	sum := 0.0
	for i := 0; i < n; i++ {
		sum += math.Abs(actual[i] - predicted[i])
	}
	return sum / float64(n)
}

func RMSE(actual, predicted []float64) float64 {
	n := len(actual)
	if n == 0 || len(predicted) == 0 {
		return 0
	}
	if len(predicted) < n {
		n = len(predicted)
	}
	sum := 0.0
	for i := 0; i < n; i++ {
		d := actual[i] - predicted[i]
		sum += d * d
	}
	return math.Sqrt(sum / float64(n))
}

func MAPE(actual, predicted []float64) float64 {
	n := len(actual)
	if n == 0 || len(predicted) == 0 {
		return 0
	}
	if len(predicted) < n {
		n = len(predicted)
	}
	sum := 0.0
	count := 0
	for i := 0; i < n; i++ {
		if actual[i] != 0 {
			sum += math.Abs((actual[i] - predicted[i]) / actual[i])
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return sum / float64(count) * 100
}

func R2(actual, predicted []float64) float64 {
	n := len(actual)
	if n == 0 || len(predicted) == 0 {
		return 0
	}
	if len(predicted) < n {
		n = len(predicted)
	}
	mean := 0.0
	for i := 0; i < n; i++ {
		mean += actual[i]
	}
	mean /= float64(n)

	ssRes := 0.0
	ssTot := 0.0
	for i := 0; i < n; i++ {
		d := actual[i] - predicted[i]
		ssRes += d * d
		d2 := actual[i] - mean
		ssTot += d2 * d2
	}
	if ssTot == 0 {
		return 0
	}
	return 1 - ssRes/ssTot
}

func MBE(actual, predicted []float64) float64 {
	n := len(actual)
	if n == 0 || len(predicted) == 0 {
		return 0
	}
	if len(predicted) < n {
		n = len(predicted)
	}
	sum := 0.0
	for i := 0; i < n; i++ {
		sum += predicted[i] - actual[i]
	}
	return sum / float64(n)
}

func NRMSE(actual, predicted []float64) float64 {
	n := len(actual)
	if n == 0 {
		return 0
	}
	minA := actual[0]
	maxA := actual[0]
	for _, a := range actual {
		if a < minA {
			minA = a
		}
		if a > maxA {
			maxA = a
		}
	}
	rng := maxA - minA
	if rng == 0 {
		return 0
	}
	return RMSE(actual, predicted) / rng
}

func Skill(actual, predicted, persistence []float64) float64 {
	rmseModel := RMSE(actual, predicted)
	rmsePersist := RMSE(actual, persistence)
	if rmsePersist == 0 {
		return 0
	}
	return 1 - rmseModel/rmsePersist
}

func KSI(actual, predicted []float64, bins int) float64 {
	if bins <= 0 || len(actual) == 0 || len(predicted) == 0 {
		return 0
	}
	n := len(actual)
	if len(predicted) < n {
		n = len(predicted)
	}
	minV := actual[0]
	maxV := actual[0]
	for i := 0; i < n; i++ {
		if actual[i] < minV {
			minV = actual[i]
		}
		if actual[i] > maxV {
			maxV = actual[i]
		}
		if predicted[i] < minV {
			minV = predicted[i]
		}
		if predicted[i] > maxV {
			maxV = predicted[i]
		}
	}
	if maxV == minV {
		return 0
	}
	binWidth := (maxV - minV) / float64(bins)
	ksi := 0.0
	for b := 0; b < bins; b++ {
		threshold := minV + float64(b+1)*binWidth
		cdfA := 0.0
		cdfP := 0.0
		for i := 0; i < n; i++ {
			if actual[i] <= threshold {
				cdfA++
			}
			if predicted[i] <= threshold {
				cdfP++
			}
		}
		ksi += math.Abs(cdfA-cdfP) / float64(n)
	}
	return ksi * binWidth
}

func Correlation(x, y []float64) float64 {
	n := len(x)
	if n == 0 || len(y) == 0 {
		return 0
	}
	if len(y) < n {
		n = len(y)
	}
	meanX, meanY := 0.0, 0.0
	for i := 0; i < n; i++ {
		meanX += x[i]
		meanY += y[i]
	}
	meanX /= float64(n)
	meanY /= float64(n)

	cov := 0.0
	varX := 0.0
	varY := 0.0
	for i := 0; i < n; i++ {
		dx := x[i] - meanX
		dy := y[i] - meanY
		cov += dx * dy
		varX += dx * dx
		varY += dy * dy
	}
	denom := math.Sqrt(varX * varY)
	if denom == 0 {
		return 0
	}
	return cov / denom
}
