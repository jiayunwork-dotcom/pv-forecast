package shading

import "math"

type Obstacle struct {
	AzStart float64
	AzEnd   float64
	Elev    float64
}

type HorizonProfile struct {
	Azimuths   []float64
	Elevations []float64
}

func IsShaded(solarAzDeg, solarElevDeg float64, profile HorizonProfile) bool {
	if len(profile.Azimuths) == 0 {
		return false
	}
	horizElev := interpolateProfile(solarAzDeg, profile)
	return solarElevDeg < horizElev
}

func ShadingFactor(solarAzDeg, solarElevDeg float64, profile HorizonProfile) float64 {
	if len(profile.Azimuths) == 0 {
		return 1.0
	}
	horizElev := interpolateProfile(solarAzDeg, profile)
	if solarElevDeg >= horizElev {
		return 1.0
	}
	if solarElevDeg <= 0 {
		return 0
	}
	return solarElevDeg / horizElev
}

func DailyShadingLoss(solarAz, solarElev []float64, profile HorizonProfile) float64 {
	if len(solarAz) == 0 {
		return 0
	}
	total := 0.0
	count := 0
	for i := range solarAz {
		if solarElev[i] > 0 {
			f := ShadingFactor(solarAz[i], solarElev[i], profile)
			total += (1 - f)
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return total / float64(count)
}

func NearShadingLoss(objectHeight, objectDist, solarElevDeg float64) float64 {
	if solarElevDeg <= 0 || objectDist <= 0 {
		return 1.0
	}
	shadowLen := objectHeight / math.Tan(solarElevDeg*math.Pi/180)
	if shadowLen <= objectDist {
		return 0
	}
	overreach := (shadowLen - objectDist) / shadowLen
	if overreach > 1 {
		return 1.0
	}
	return overreach * 0.7
}

func SelfShadingLoss(rowSpacing, moduleHeight, tiltDeg, solarElevDeg float64) float64 {
	if solarElevDeg <= 0 {
		return 1.0
	}
	tilt := tiltDeg * math.Pi / 180
	elev := solarElevDeg * math.Pi / 180
	shadowLen := moduleHeight * math.Sin(tilt) / math.Tan(elev)
	if shadowLen <= rowSpacing {
		return 0
	}
	fraction := (shadowLen - rowSpacing) / (moduleHeight * math.Cos(tilt))
	if fraction > 1 {
		fraction = 1
	}
	if fraction < 0 {
		fraction = 0
	}
	return fraction
}

func AnnualShadingLoss(hourlyAz, hourlyElev, hourlyGHI []float64, profile HorizonProfile) float64 {
	n := len(hourlyAz)
	if n == 0 || len(hourlyElev) < n || len(hourlyGHI) < n {
		return 0
	}
	totalGHI := 0.0
	lostGHI := 0.0
	for i := 0; i < n; i++ {
		if hourlyElev[i] <= 0 || hourlyGHI[i] <= 0 {
			continue
		}
		totalGHI += hourlyGHI[i]
		f := ShadingFactor(hourlyAz[i], hourlyElev[i], profile)
		lostGHI += hourlyGHI[i] * (1 - f)
	}
	if totalGHI == 0 {
		return 0
	}
	return lostGHI / totalGHI
}

func interpolateProfile(az float64, profile HorizonProfile) float64 {
	n := len(profile.Azimuths)
	if n == 0 {
		return 0
	}
	if az <= profile.Azimuths[0] {
		return profile.Elevations[0]
	}
	if az >= profile.Azimuths[n-1] {
		return profile.Elevations[n-1]
	}
	for i := 1; i < n; i++ {
		if az <= profile.Azimuths[i] {
			frac := (az - profile.Azimuths[i-1]) / (profile.Azimuths[i] - profile.Azimuths[i-1])
			return profile.Elevations[i-1] + frac*(profile.Elevations[i]-profile.Elevations[i-1])
		}
	}
	return 0
}

func ObstacleToProfile(obstacles []Obstacle, step float64) HorizonProfile {
	if step <= 0 {
		step = 1
	}
	n := int(360 / step)
	profile := HorizonProfile{
		Azimuths:   make([]float64, n),
		Elevations: make([]float64, n),
	}
	for i := range profile.Azimuths {
		profile.Azimuths[i] = float64(i) * step
	}
	for _, obs := range obstacles {
		for i, az := range profile.Azimuths {
			if az >= obs.AzStart && az <= obs.AzEnd {
				if obs.Elev > profile.Elevations[i] {
					profile.Elevations[i] = obs.Elev
				}
			}
		}
	}
	return profile
}
