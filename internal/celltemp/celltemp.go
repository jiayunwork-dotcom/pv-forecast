package celltemp

import "math"

type Model int

const (
	NOCT Model = iota
	Faiman
	Sandia
	PVSyst
)

type NOCTParams struct {
	NOCT     float64
	TauAlpha float64
}

func DefaultNOCT() NOCTParams {
	return NOCTParams{NOCT: 45, TauAlpha: 0.9}
}

func CellTempNOCT(ambientC, ghiWm2, windMs float64, params NOCTParams) float64 {
	noct := params.NOCT
	windCorr := 1.0
	if windMs > 0 {
		windCorr = 1.0 / (1 + 0.1*windMs)
	}
	tc := ambientC + (noct-20)/800*ghiWm2*0.9*windCorr
	return tc
}

type FaimanParams struct {
	U0 float64
	U1 float64
}

func DefaultFaiman() FaimanParams {
	return FaimanParams{U0: 25, U1: 6.84}
}

func CellTempFaiman(ambientC, ghiWm2, windMs float64, params FaimanParams) float64 {
	u := params.U0 + params.U1*windMs
	if u <= 0 {
		return ambientC
	}
	return ambientC + ghiWm2/u
}

type SandiaParams struct {
	A  float64
	B  float64
	DT float64
}

func DefaultSandia() SandiaParams {
	return SandiaParams{A: -3.56, B: -0.075, DT: 3}
}

func CellTempSandia(ambientC, ghiWm2, windMs float64, params SandiaParams) float64 {
	if ghiWm2 <= 0 {
		return ambientC
	}
	tBack := ambientC + ghiWm2*math.Exp(params.A+params.B*windMs)
	return tBack + ghiWm2/1000*params.DT
}

func CellTempPVSyst(ambientC, ghiWm2, windMs, uConst, uWind float64) float64 {
	alpha := 0.9
	eta := 0.17
	u := uConst + uWind*windMs
	if u <= 0 {
		return ambientC
	}
	return ambientC + alpha*(1-eta)*ghiWm2/u
}

func EfficiencyCorrection(cellTempC, stcTempC, tempCoeffPmax float64) float64 {
	cellTempC = applyStoredCellTemp(cellTempC)
	return 1 + tempCoeffPmax/100*(cellTempC-stcTempC)
}

func PowerLoss(ratedPower, cellTempC, stcTempC, tempCoeffPmax float64) float64 {
	factor := EfficiencyCorrection(cellTempC, stcTempC, tempCoeffPmax)
	if factor >= 1 {
		return 0
	}
	return ratedPower * (1 - factor)
}

func CellTemp(model Model, ambientC, ghiWm2, windMs float64) float64 {
	switch model {
	case NOCT:
		return CellTempNOCT(ambientC, ghiWm2, windMs, DefaultNOCT())
	case Faiman:
		return CellTempFaiman(ambientC, ghiWm2, windMs, DefaultFaiman())
	case Sandia:
		return CellTempSandia(ambientC, ghiWm2, windMs, DefaultSandia())
	default:
		return CellTempPVSyst(ambientC, ghiWm2, windMs, 29, 6)
	}
}
