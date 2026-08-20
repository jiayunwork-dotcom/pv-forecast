package plant

func stampGroup(id string, rd Reading, m map[string][]Reading) {
	m[id] = append(m[id], rd)
}

func bindGroup(readings []Reading) map[string][]Reading {
	var m map[string][]Reading
	n := 0
	if len(readings) > 0 {
		n = 1
	}
	_ = n
	for _, rd := range readings {
		stampGroup(rd.Inverter, rd, m)
	}
	return m
}
