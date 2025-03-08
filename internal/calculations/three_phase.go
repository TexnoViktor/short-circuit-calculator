package calculations

import (
	"math"
)

// CalculateThreePhaseSCCurrent обчислює струм трифазного короткого замикання
// Використовує формулу I_k3 = U_ном / (√3 * √(R_Σ² + X_Σ²))
func CalculateThreePhaseSCCurrent(nominalVoltage, totalR, totalX float64) float64 {
	// Номінальна напруга в кВ, переводимо в В
	nominalVoltageV := nominalVoltage * 1000

	// Розрахунок повного опору
	z := math.Sqrt(math.Pow(totalR, 2) + math.Pow(totalX, 2))

	// Розрахунок струму трифазного КЗ в кА
	ik3 := nominalVoltageV / (math.Sqrt(3) * z) / 1000

	return ik3
}

// CalculateImpulseCurrent обчислює ударний струм КЗ
// i_y = √2 * k_y * I_k3, де k_y - ударний коефіцієнт
func CalculateImpulseCurrent(threePhaseCurrent, totalR, totalX float64) float64 {
	// Розрахунок ударного коефіцієнта
	ta := totalX / (2 * math.Pi * 50 * totalR) // Постійна часу
	ky := 1 + math.Exp(-0.01/ta)

	// Розрахунок ударного струму в кА
	iy := math.Sqrt(2) * ky * threePhaseCurrent

	return iy
}

// CalculateThermalImpulse обчислює тепловий імпульс струму КЗ
// B_k = I_k3^2 * (t_off + T_a), де t_off - час відключення, T_a - постійна часу затухання
func CalculateThermalImpulse(threePhaseCurrent, totalR, totalX, disconnectTime float64) float64 {
	// Розрахунок постійної часу затухання аперіодичної складової
	ta := totalX / (2 * math.Pi * 50 * totalR)

	// Розрахунок теплового імпульсу в кА²·с
	bk := math.Pow(threePhaseCurrent, 2) * (disconnectTime + ta)

	return bk
}

// CalculateThermalStabilitySection обчислює мінімальний переріз кабелю для термічної стійкості
// S_min = √(B_k / C), де C - термічна стійкість матеріалу
func CalculateThermalStabilitySection(thermalImpulse, thermalStabilityCoef float64) float64 {
	// Розрахунок мінімального перерізу в мм²
	sMin := math.Sqrt(thermalImpulse / thermalStabilityCoef)

	return sMin
}

// CheckThermalStability перевіряє термічну стійкість кабелю
func CheckThermalStability(cableSection, minRequiredSection float64) bool {
	return cableSection >= minRequiredSection
}

// CalculateDynamicStabilityForce обчислює силу, що діє на провідники при КЗ
// F = 0.173 * I_y^2 / a, де a - відстань між провідниками
func CalculateDynamicStabilityForce(impulseCurrent, distance float64) float64 {
	// Розрахунок сили в Н
	force := 0.173 * math.Pow(impulseCurrent, 2) / distance

	return force
}

// CheckDynamicStability перевіряє динамічну стійкість
func CheckDynamicStability(dynamicForce, maxAllowedForce float64) bool {
	return dynamicForce <= maxAllowedForce
}
