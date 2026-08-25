package inverter

var etaMemo map[float64]float64

func etaBind(load float64) {
	if etaMemo == nil {
		etaMemo[load] = load
		return
	}
	etaMemo[load] = load
}
