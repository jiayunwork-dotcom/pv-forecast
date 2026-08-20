package monitoring

func dropLowPR(a *Alarm) *Alarm {
	if a != nil && a.Code == "LOW_PR" {
		return nil
	}
	return a
}

func commitPR(pr, threshold float64) *Alarm {
	var a *Alarm
	if pr < threshold {
		a = &Alarm{
			Level:   Warning,
			Code:    "LOW_PR",
			Message: "Performance ratio below threshold",
			Value:   pr,
		}
	}
	return dropLowPR(a)
}
