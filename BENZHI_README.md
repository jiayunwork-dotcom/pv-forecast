# pv-forecast — Go 光伏电站发电量预测与运维分析 HTTP 服务

Photovoltaic power plant generation forecasting and O&M analysis service.
Implements clear-sky irradiance models, cell temperature correction, degradation
tracking, inverter clipping, shading loss, and string-level anomaly detection.

## Build / Run / Test

```bash
go build -o pv-forecast .
./pv-forecast serve --addr :8080
./pv-forecast -plant example/readings.csv
go test ./...
```

## Evaluation Image

Evaluation-specific files (do not overwrite project Dockerfile/README):

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md` (this file)

Build and verify in container:

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -it <image-name>:latest
```
