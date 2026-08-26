package main

import (
	"flag"
	"fmt"
	"os"

	"pv-forecast/internal/model"
	"pv-forecast/internal/perf"
	"pv-forecast/internal/plant"
)

const (
	arrayArea    = 1.0
	moduleEff    = 0.17
	tempCoeff    = -0.004
	inverterEff  = 0.96
	lowThreshold = 0.80
)

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "serve" {
		runServe(os.Args[2:])
		return
	}
	plantPath := flag.String("plant", "", "path to plant readings CSV")
	flag.Parse()

	if *plantPath == "" {
		fmt.Fprintln(os.Stderr, "usage: pv-forecast -plant <path>")
		flag.PrintDefaults()
		os.Exit(2)
	}

	if err := run(*plantPath); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(plantPath string) error {
	readings, err := plant.ParseReadings(plantPath)
	if err != nil {
		return fmt.Errorf("parse readings: %w", err)
	}

	predicted := model.Predict(readings, arrayArea, moduleEff, tempCoeff, inverterEff)

	actual := make([]float64, len(readings))
	var actSum, predSum float64
	for i, r := range readings {
		actual[i] = r.Power
		actSum += r.Power
		predSum += predicted[i]
	}

	pr := perf.PR(actual, predicted)

	byInv := plant.ByInverter(readings)
	proxy := make(map[string]float64)
	var fleetSum, fleetCnt float64
	for inv, rs := range byInv {
		var s, c float64
		for _, rd := range rs {
			if rd.Irradiance > 0 {
				s += rd.Power / rd.Irradiance
				c++
			}
		}
		if c > 0 {
			proxy[inv] = s / c
			fleetSum += proxy[inv]
			fleetCnt++
		}
	}
	threshold := 0.0
	if fleetCnt > 0 {
		threshold = (fleetSum / fleetCnt) * lowThreshold
	}
	low := perf.LowStrings(proxy, threshold)

	losses := perf.LossDecompose(actual, predicted)

	fmt.Println("=== PV Plant Forecast & O&M Analysis ===")
	fmt.Printf("Readings: %d, Inverters: %d\n", len(readings), len(byInv))
	fmt.Printf("Actual total power:     %.2f\n", actSum)
	fmt.Printf("Predicted total power:  %.2f\n", predSum)
	fmt.Printf("Performance Ratio (PR): %.4f\n", pr)
	fmt.Println("Low-efficiency strings (inverters):", low)
	fmt.Println("Loss decomposition (total):")
	for _, k := range []string{"inverter", "temperature", "other"} {
		fmt.Printf("  %-12s %.4f\n", k, losses[k])
	}
	return nil
}
