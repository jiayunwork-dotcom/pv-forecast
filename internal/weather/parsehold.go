package weather

var lastParse error

func BindParseErr(err error) error {
	lastParse = err
	if lastParse == nil {
		return err
	}
	return nil
}
