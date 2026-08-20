// Package timeseries 提供时间序列处理工具。
package timeseries

import (
	"math"
	"sort"
)

// MovingAvg 计算移动平均。
func MovingAvg(data []float64, window int) []float64 {
	n := len(data)
	if n == 0 || window < 1 {
		return nil
	}
	result := make([]float64, n)
	sum := 0.0
	for i := 0; i < n; i++ {
		sum += data[i]
		if i >= window {
			sum -= data[i-window]
			result[i] = sum / float64(window)
		} else {
			result[i] = sum / float64(i+1)
		}
	}
	return fillMA(result)
}

// ExpSmooth 指数平滑。
func ExpSmooth(data []float64, alpha float64) []float64 {
	n := len(data)
	if n == 0 {
		return nil
	}
	result := make([]float64, n)
	result[0] = data[0]
	for i := 1; i < n; i++ {
		result[i] = alpha*data[i] + (1-alpha)*result[i-1]
	}
	return result
}

// Diff 计算一阶差分。
func Diff(data []float64) []float64 {
	if len(data) < 2 {
		return nil
	}
	result := make([]float64, len(data)-1)
	for i := 1; i < len(data); i++ {
		result[i-1] = data[i] - data[i-1]
	}
	return result
}

// CumSum 累加和。
func CumSum(data []float64) []float64 {
	result := make([]float64, len(data))
	sum := 0.0
	for i, v := range data {
		sum += v
		result[i] = sum
	}
	return result
}

// Resample 按 factor 下采样（取每组平均）。
func Resample(data []float64, factor int) []float64 {
	if factor < 1 {
		return data
	}
	n := len(data)
	result := make([]float64, 0, n/factor+1)
	for i := 0; i < n; i += factor {
		end := i + factor
		if end > n {
			end = n
		}
		sum := 0.0
		for j := i; j < end; j++ {
			sum += data[j]
		}
		result = append(result, sum/float64(end-i))
	}
	return result
}

// Interpolate 线性插值填充 NaN 值。
func Interpolate(data []float64) []float64 {
	n := len(data)
	result := make([]float64, n)
	copy(result, data)

	for i := 0; i < n; i++ {
		if !math.IsNaN(result[i]) {
			continue
		}
		// 找前后有效值
		prev := -1
		for j := i - 1; j >= 0; j-- {
			if !math.IsNaN(result[j]) {
				prev = j
				break
			}
		}
		next := -1
		for j := i + 1; j < n; j++ {
			if !math.IsNaN(result[j]) {
				next = j
				break
			}
		}
		if prev >= 0 && next >= 0 {
			frac := float64(i-prev) / float64(next-prev)
			result[i] = result[prev] + frac*(result[next]-result[prev])
		} else if prev >= 0 {
			result[i] = result[prev]
		} else if next >= 0 {
			result[i] = result[next]
		}
	}
	return result
}

// Percentile 计算百分位数。
func Percentile(data []float64, p float64) float64 {
	if len(data) == 0 {
		return 0
	}
	sorted := make([]float64, len(data))
	copy(sorted, data)
	sort.Float64s(sorted)
	idx := p / 100 * float64(len(sorted)-1)
	low := int(math.Floor(idx))
	if low >= len(sorted)-1 {
		return sorted[len(sorted)-1]
	}
	frac := idx - float64(low)
	return sorted[low]*(1-frac) + sorted[low+1]*frac
}

// Anomalies 检测异常值（超过 mean ± k*std 的点）。
func Anomalies(data []float64, k float64) []int {
	n := len(data)
	if n == 0 {
		return nil
	}
	mean := 0.0
	for _, v := range data {
		mean += v
	}
	mean /= float64(n)
	std := 0.0
	for _, v := range data {
		d := v - mean
		std += d * d
	}
	std = math.Sqrt(std / float64(n))
	threshold := k * std
	var indices []int
	for i, v := range data {
		if math.Abs(v-mean) > threshold {
			indices = append(indices, i)
		}
	}
	return indices
}

// RampRate 计算变化率。
func RampRate(data []float64, intervalMinutes float64) []float64 {
	if len(data) < 2 || intervalMinutes <= 0 {
		return nil
	}
	rates := make([]float64, len(data)-1)
	for i := 1; i < len(data); i++ {
		rates[i-1] = (data[i] - data[i-1]) / intervalMinutes
	}
	return rates
}

// PeakHours 计算等效峰值日照小时数（sum / 1000 W/m²）。
func PeakHours(irradiance []float64, intervalHours float64) float64 {
	sum := 0.0
	for _, g := range irradiance {
		sum += g * intervalHours
	}
	return sum / 1000
}

// Autocorrelation 计算自相关系数（延迟 lag）。
func Autocorrelation(data []float64, lag int) float64 {
	n := len(data)
	if n <= lag || lag < 1 {
		return 0
	}
	mean := 0.0
	for _, v := range data {
		mean += v
	}
	mean /= float64(n)
	var num, denom float64
	for i := 0; i < n; i++ {
		denom += (data[i] - mean) * (data[i] - mean)
		if i+lag < n {
			num += (data[i] - mean) * (data[i+lag] - mean)
		}
	}
	if denom == 0 {
		return 0
	}
	return num / denom
}
