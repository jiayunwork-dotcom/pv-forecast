// Package celltemp 实现光伏电池温度模型。
package celltemp

import "math"

// Model 电池温度模型。
type Model int

const (
	NOCT     Model = iota // 标称工作温度模型
	Faiman               // Faiman 模型
	Sandia               // Sandia 模型
	PVSyst               // PVSyst 模型
)

// NOCTParams NOCT 模型参数。
type NOCTParams struct {
	NOCT      float64 // 标称工作温度 (°C)，通常 44-46
	TauAlpha  float64 // 透射吸收积
}

// DefaultNOCT 返回默认 NOCT 参数。
func DefaultNOCT() NOCTParams {
	return NOCTParams{NOCT: 45, TauAlpha: 0.9}
}

// CellTempNOCT 使用 NOCT 模型计算电池温度。
func CellTempNOCT(ambientC, ghiWm2, windMs float64, params NOCTParams) float64 {
	noct := params.NOCT
	// 简化 NOCT 模型：Tc = Ta + (NOCT-20)/800 * G * (1 - eta/tau_alpha)
	// 假设 eta/tau_alpha ≈ 0.1
	windCorr := 1.0
	if windMs > 0 {
		windCorr = 1.0 / (1 + 0.1*windMs)
	}
	tc := ambientC + (noct-20)/800*ghiWm2*0.9*windCorr
	return tc
}

// FaimanParams Faiman 模型参数。
type FaimanParams struct {
	U0 float64 // 常数损失系数 (W/m²·K)
	U1 float64 // 风速损失系数 (W/m³·s·K)
}

// DefaultFaiman 返回默认 Faiman 参数。
func DefaultFaiman() FaimanParams {
	return FaimanParams{U0: 25, U1: 6.84}
}

// CellTempFaiman 使用 Faiman 模型。
func CellTempFaiman(ambientC, ghiWm2, windMs float64, params FaimanParams) float64 {
	u := params.U0 + params.U1*windMs
	if u <= 0 {
		return ambientC
	}
	return ambientC + ghiWm2/u
}

// SandiaParams Sandia 电池温度模型参数。
type SandiaParams struct {
	A float64 // 经验常数
	B float64 // 经验常数
	DT float64 // 电池到模块背面温差
}

// DefaultSandia 返回默认 Sandia 参数（开放式支架）。
func DefaultSandia() SandiaParams {
	return SandiaParams{A: -3.56, B: -0.075, DT: 3}
}

// CellTempSandia 使用 Sandia 模型。
func CellTempSandia(ambientC, ghiWm2, windMs float64, params SandiaParams) float64 {
	if ghiWm2 <= 0 {
		return ambientC
	}
	tBack := ambientC + ghiWm2*math.Exp(params.A+params.B*windMs)
	return tBack + ghiWm2/1000*params.DT
}

// CellTempPVSyst PVSyst 线性模型。
func CellTempPVSyst(ambientC, ghiWm2, windMs, uConst, uWind float64) float64 {
	alpha := 0.9 // 吸收率
	eta := 0.17  // 模块效率
	u := uConst + uWind*windMs
	if u <= 0 {
		return ambientC
	}
	return ambientC + alpha*(1-eta)*ghiWm2/u
}

// EfficiencyCorrection 温度对效率的修正因子。
func EfficiencyCorrection(cellTempC, stcTempC, tempCoeffPmax float64) float64 {
	return 1 + tempCoeffPmax/100*(cellTempC-stcTempC)
}

// PowerLoss 温度导致的功率损失 (W)。
func PowerLoss(ratedPower, cellTempC, stcTempC, tempCoeffPmax float64) float64 {
	factor := EfficiencyCorrection(cellTempC, stcTempC, tempCoeffPmax)
	if factor >= 1 {
		return 0
	}
	return ratedPower * (1 - factor)
}

// CellTemp 根据模型类型计算电池温度。
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
