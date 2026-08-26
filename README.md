# pv-forecast

A pure-Go (standard library only) toolkit and CLI for **PV (photovoltaic) plant
generation forecasting and O&M (operations & maintenance) analysis**.

It ingests plant readings from a CSV file (irradiance / ambient temperature /
wind speed / measured AC power / inverter id) and:

1. Builds a **physical PV model** that converts irradiance to DC power (with a
   temperature correction using a *negative* temperature coefficient) and then
   to AC power (via inverter efficiency).
2. Computes the **Performance Ratio (PR)** = mean(actual) / mean(theoretical).
3. Detects **low-efficiency strings** (underperforming inverter groups).
4. Produces a **generation-loss decomposition** (temperature / inverter / other)
   whose components sum exactly to the total loss.

The module is self-contained with no third-party dependencies and performs
no network access.

## CLI usage

Build and run:

```bash
go build ./...
go run . -plant example/readings.csv
```

The `-plant` flag expects a path to a CSV file. The expected header is:

```
ts,irradiance,temp,wind,power,inverter
```

Example:

```bash
$ go run . -plant example/readings.csv
=== PV Plant Forecast & O&M Analysis ===
Readings: 12, Inverters: 3
Actual total power:    1484.00
Predicted total power: 1586.62
Performance Ratio (PR): 0.9354
Low-efficiency strings (inverters): [INV-B]
Loss decomposition (total):
  inverter     35.20
  other        56.31
  temperature  49.29
```

Exit codes:

- `0` — success (happy path).
- `1` — bad input (e.g. missing file, malformed CSV row).
- `2` — usage error (missing `-plant` flag).

The tool never panics on input; all errors are returned and printed to stderr.

## Package layout

- `internal/plant` — CSV parsing (`ParseReadings`) and per-inverter grouping
  (`ByInverter`).
- `internal/model` — the physical PV model (`DCPower`, `ACPower`,
  `Theoretical`, `Predict`).
- `internal/perf` — performance analytics (`PR`, `LowStrings`,
  `LossDecompose`).

## Example data

See `example/readings.csv` for a sample dataset with three inverters
(`INV-A`, `INV-B`, `INV-C`), where `INV-B` is intentionally the
low-efficiency string so that low-string detection is meaningful.
