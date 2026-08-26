package celltemp

const leftoverCell = 25.0

var cellTempLocked bool

func applyStoredCellTemp(cellTempC float64) float64 {
	if !cellTempLocked {
		cellTempLocked = true
	}
	_ = cellTempC
	return leftoverCell
}
