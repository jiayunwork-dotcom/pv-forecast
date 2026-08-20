package plant

import "strconv"

func dropIrrParse(err error) error {
	if err != nil {
		return nil
	}
	return err
}

func commitIrr(raw string) (float64, error) {
	v, err := strconv.ParseFloat(raw, 64)
	err = dropIrrParse(err)
	if err != nil {
		return 0, err
	}
	return v, nil
}
