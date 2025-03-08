package calculations

import (
	"math"
)

// CalculateSinglePhaseSCCurrent обчислює струм однофазного короткого замикання
// Використовує формулу I_k1 = 3 * U_ф / (Z_1 + Z_2 + Z_0), де Z_1, Z_2, Z_0 - опори прямої, зворотньої і нульової послідовності
func CalculateSinglePhaseSCCurrent(nominalVoltage, r1, x1, r0, x0 float64) float64 {
	// Номінальна напруга в кВ, переводимо в В
	nominalVoltageV := nominalVoltage * 1000

	// Фазна напруга
	phaseVoltage := nominalVoltageV / math.Sqrt(3)

	// Опори прямої і зворотньої послідовності (для симетричної системи Z_1 = Z_2)
	z1 := math.Sqrt(math.Pow(r1, 2) + math.Pow(x1, 2))
	z2 := z1

	// Опір нульової послідовності
	z0 := math.Sqrt(math.Pow(r0, 2) + math.Pow(x0, 2))

	// Розрахунок струму однофазного КЗ в кА
	ik1 := 3 * phaseVoltage / (z1 + z2 + z0) / 1000

	return ik1
}

// CalculateSinglePhaseSCCurrentFromRatio обчислює струм однофазного КЗ через відношення до трифазного
// Для мереж 6-10 кВ типово однофазний струм приймається як відсоток від трифазного
func CalculateSinglePhaseSCCurrentFromRatio(threePhaseCurrent, ratio float64) float64 {
	return threePhaseCurrent * ratio
}

// CalculateZeroSequenceImpedance обчислює опір нульової послідовності
// Спрощена формула для кабельних ліній
func CalculateZeroSequenceImpedance(r1, x1 float64) (r0, x0 float64) {
	// Для кабельних ліній приймаємо, що r0 ≈ 4*r1, x0 ≈ 3*x1
	r0 = 4 * r1
	x0 = 3 * x1

	return r0, x0
}

// CalculateSinglePhaseCurrentForEMStation обчислює струм однофазного КЗ для підстанції ЕМ
// Використовує спрощену модель для розрахунку
func CalculateSinglePhaseCurrentForEMStation(threePhaseCurrent float64) float64 {
	// Приймаємо для підстанції ХПнЕМ коефіцієнт = 0.7 (70% від трифазного)
	// Значення може бути уточнене на основі більш детальних даних
	return threePhaseCurrent * 0.7
}
