package calculations

import (
	"math"
)

// MaterialConstantsForThermalStability містить константи для різних матеріалів
var MaterialConstantsForThermalStability = map[string]float64{
	"copper":   170, // для мідних жил
	"aluminum": 90,  // для алюмінієвих жил
	"steel":    75,  // для сталевих елементів
}

// Параметри для розрахунку динамічної стійкості різних типів обладнання
type DynamicStabilityParams struct {
	Type           string  // Тип обладнання
	MaxForcePerLen float64 // Максимально допустиме навантаження на одиницю довжини, Н/м
	Factor         float64 // Коефіцієнт для розрахунку
}

// DynamicStabilityParamsDB містить параметри для різних типів обладнання
var DynamicStabilityParamsDB = []DynamicStabilityParams{
	{"busbar_10kV", 500, 1.5},
	{"cable_10kV", 1200, 1.2},
	{"switch_10kV", 3000, 1.8},
}

// CalculateThermalStabilitySection обчислює мінімальний переріз для термічної стійкості
func CalculateThermalStabilityCopper(thermalImpulse float64) float64 {
	return math.Sqrt(thermalImpulse / MaterialConstantsForThermalStability["copper"])
}

func CalculateThermalStabilityAluminum(thermalImpulse float64) float64 {
	return math.Sqrt(thermalImpulse / MaterialConstantsForThermalStability["aluminum"])
}

// CalculateDynamicStabilityImpulseCurrent обчислює ударний струм для перевірки динамічної стійкості
func CalculateDynamicStabilityImpulseCurrent(threePhaseCurrent, ratio float64) float64 {
	return threePhaseCurrent * ratio
}

// EvaluateDynamicStability оцінює динамічну стійкість обладнання
func EvaluateDynamicStability(equipmentType string, dynamicForce float64) bool {
	for _, params := range DynamicStabilityParamsDB {
		if params.Type == equipmentType {
			return dynamicForce <= params.MaxForcePerLen*params.Factor
		}
	}

	// Якщо тип обладнання невідомий, повертаємо false для безпеки
	return false
}

// CalculateDisconnectTime обчислює час відключення для розрахунку теплового імпульсу
// Спрощена формула, може бути уточнена для конкретних умов
func CalculateDisconnectTime(relayTime, breakerTime float64) float64 {
	// Тут relayTime - час спрацювання захисту, breakerTime - власний час відключення вимикача
	return relayTime + breakerTime
}

// CheckBusbarThermalStability перевіряє термічну стійкість шин
func CheckBusbarThermalStability(busbarSection, thermalStabilityMin float64) bool {
	return busbarSection >= thermalStabilityMin
}

// CheckBusbarDynamicStability перевіряє динамічну стійкість шин
func CheckBusbarDynamicStability(busbarType string, dynamicForce float64) bool {
	// Різні типи шин мають різні обмеження по динамічній стійкості
	maxAllowedForce := 0.0

	switch busbarType {
	case "aluminum_10kV":
		maxAllowedForce = 2000.0
	case "copper_10kV":
		maxAllowedForce = 3000.0
	case "steel_10kV":
		maxAllowedForce = 4000.0
	default:
		maxAllowedForce = 1500.0 // За замовчуванням, для безпеки
	}

	return dynamicForce <= maxAllowedForce
}

// CalculateElectrodynamicForce обчислює електродинамічну силу між провідниками при КЗ
func CalculateElectrodynamicForce(impulseCurrent, conductorLength, distance float64) float64 {
	// Сила між паралельними провідниками при КЗ
	// F = 2*10^(-7) * i_y^2 * l / a, де l - довжина, a - відстань між провідниками
	return 2.0e-7 * math.Pow(impulseCurrent, 2) * conductorLength / distance
}

// CalculateMinThermalSection обчислює мінімальний переріз провідника для термічної стійкості
func CalculateMinThermalSection(thermalImpulse, thermalCoefficient float64) float64 {
	return math.Sqrt(thermalImpulse / thermalCoefficient)
}

// VerifyEquipmentStability перевіряє стійкість обладнання до струмів КЗ
func VerifyEquipmentStability(equipmentType string, nominalCurrent, threePhaseCurrent, thermalImpulse float64) bool {
	// Спрощена перевірка обладнання на стійкість до струмів КЗ

	// Коефіцієнти для різних типів обладнання
	kDynamic := 0.0
	kThermal := 0.0

	switch equipmentType {
	case "circuit_breaker_10kV":
		kDynamic = 25.0 // Кратність динамічної стійкості до номінального струму
		kThermal = 16.0 // Кратність термічної стійкості до номінального струму
	case "disconnector_10kV":
		kDynamic = 20.0
		kThermal = 12.0
	case "CT_10kV": // Трансформатор струму
		kDynamic = 30.0
		kThermal = 20.0
	default:
		kDynamic = 15.0
		kThermal = 10.0
	}

	// Перевірка динамічної стійкості
	dynamicStability := threePhaseCurrent*math.Sqrt(2) <= nominalCurrent*kDynamic

	// Перевірка термічної стійкості (спрощено)
	thermalStability := thermalImpulse <= math.Pow(nominalCurrent*kThermal, 2)

	return dynamicStability && thermalStability
}
