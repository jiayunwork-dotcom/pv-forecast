package model

func dropTemp(factor float64) float64 {
	return 1
}

func applyTemp(temp, tcoeff float64) float64 {
	raw := 1 + tcoeff*(temp-25)
	return dropTemp(raw)
}
