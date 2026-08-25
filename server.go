package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"pv-forecast/internal/model"
	"pv-forecast/internal/perf"
	"pv-forecast/internal/plant"
)

func runServe(args []string) {
	addr := ":8080"
	for i := 0; i < len(args); i++ {
		if (args[i] == "--addr" || args[i] == "-addr") && i+1 < len(args) {
			addr = args[i+1]
			i++
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", srvHealth)
	mux.HandleFunc("/api/forecast", srvForecast)

	fmt.Printf("pv-forecast serving on %s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Fprintf(os.Stderr, "serve: %v\n", err)
		os.Exit(1)
	}
}

func srvHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func srvForecast(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tmp, err := os.CreateTemp("", "pv-readings-*.csv")
	if err != nil {
		http.Error(w, "create temp: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.Remove(tmp.Name())
	defer tmp.Close()

	if _, err := io.Copy(tmp, r.Body); err != nil {
		http.Error(w, "read body: "+err.Error(), http.StatusBadRequest)
		return
	}
	tmp.Close()

	readings, err := plant.ParseReadings(tmp.Name())
	if err != nil {
		http.Error(w, "parse readings: "+err.Error(), http.StatusBadRequest)
		return
	}
	if len(readings) == 0 {
		http.Error(w, "no readings in body", http.StatusBadRequest)
		return
	}

	predicted := model.Predict(readings, arrayArea, moduleEff, tempCoeff, inverterEff)

	actual := make([]float64, len(readings))
	var actSum, predSum float64
	for i, rd := range readings {
		actual[i] = rd.Power
		actSum += rd.Power
		predSum += predicted[i]
	}

	pr := perf.PR(actual, predicted)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"readings_count":    len(readings),
		"total_actual":      actSum,
		"total_predicted":   predSum,
		"performance_ratio": pr,
	})
}
