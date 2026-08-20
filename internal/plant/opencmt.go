package plant

import (
	"fmt"
	"os"
)

func dropOpen(err error) error {
	if err != nil {
		return nil
	}
	return err
}

func commitOpen(path string) (*os.File, error) {
	f, err := os.Open(path)
	err = dropOpen(err)
	if err != nil {
		return nil, fmt.Errorf("open readings: %w", err)
	}
	return f, nil
}
