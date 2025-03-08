package calculations

import (
	"math"
)

// CableData містить параметри кабелю
type CableData struct {
	Type         string
	Section      float64 // мм²
	Current      float64 // A
	Resistance   float64 // Ом/км
	Reactance    float64 // Ом/км
	ThermalConst float64 // Константа термічної стійкості
	MaxLoadForce float64 // Максимальне навантаження для динамічної стійкості, Н
}

// CableDatabase містить доступні типи кабелів
var CableDatabase = []CableData{
	{"ААШв", 16, 90, 1.94, 0.095, 93, 900},
	{"ААШв", 25, 115, 1.24, 0.091, 93, 1200},
	{"ААШв", 35, 135, 0.89, 0.087, 93, 1500},
	{"ААШв", 50, 165, 0.62, 0.083, 93, 1800},
	{"ААШв", 70, 200, 0.443, 0.08, 93, 2300},
	{"ААШв", 95, 240, 0.326, 0.078, 93, 2700},
	{"ААШв", 120, 275, 0.258, 0.076, 93, 3100},
	{"ААШв", 150, 310, 0.206, 0.074, 93, 3500},
	{"ААШв", 185, 355, 0.167, 0.073, 93, 4000},
	{"ААШв", 240, 410, 0.129, 0.071, 93, 4700},
	{"СБ", 16, 80, 1.15, 0.12, 160, 1200},
	{"СБ", 25, 105, 0.74, 0.113, 160, 1500},
	{"СБ", 35, 125, 0.52, 0.106, 160, 1900},
	{"СБ", 50, 155, 0.37, 0.1, 160, 2300},
	{"СБ", 70, 190, 0.265, 0.095, 160, 2800},
	{"СБ", 95, 225, 0.195, 0.091, 160, 3300},
	{"СБ", 120, 260, 0.154, 0.088, 160, 3800},
	{"СБ", 150, 300, 0.124, 0.086, 160, 4200},
	{"СБ", 185, 340, 0.1, 0.083, 160, 4800},
	{"СБ", 240, 400, 0.077, 0.079, 160, 5500},
}

// CalculateRatedCurrent обчислює номінальний струм кабелю
func CalculateRatedCurrent(power, voltage, cosFi float64) float64 {
	// Power в кВА, voltage в кВ
	return power / (math.Sqrt(3) * voltage * cosFi)
}

// ApplyLoadFactors застосовує коефіцієнти навантаження для визначення допустимого навантаження
func ApplyLoadFactors(baseLoad, ambientTemp, soilThermalRes float64, installMethod string, cableCount int) float64 {
	// Базове значення вже має коефіцієнти для стандартних умов

	// Коефіцієнт для температури (спрощено)
	kTemp := 1.0
	if ambientTemp > 25 {
		kTemp = 1.0 - 0.01*(ambientTemp-25)
	} else if ambientTemp < 25 {
		kTemp = 1.0 + 0.01*(25-ambientTemp)
	}

	// Коефіцієнт для опору ґрунту (спрощено)
	kSoil := 1.0
	if installMethod == "ground" && soilThermalRes > 1.2 {
		kSoil = 1.0 - 0.05*(soilThermalRes-1.2)
	}

	// Коефіцієнт для кількості кабелів поряд (спрощено)
	kCount := 1.0
	if cableCount > 1 {
		switch cableCount {
		case 2:
			kCount = 0.9
		case 3:
			kCount = 0.85
		default:
			kCount = 0.8
		}
	}

	// Розрахунок допустимого навантаження з урахуванням усіх коефіцієнтів
	permissibleLoad := baseLoad * kTemp * kSoil * kCount

	return permissibleLoad
}

// CalculateVoltageDropPercent обчислює відсоток падіння напруги
func CalculateVoltageDropPercent(current, resistance, reactance, length, voltage, cosFi float64) float64 {
	sinFi := math.Sqrt(1 - cosFi*cosFi)

	// Падіння напруги у відсотках
	voltageDropPerc := (math.Sqrt(3) * current * length * (resistance*cosFi + reactance*sinFi)) / (10 * voltage)

	return voltageDropPerc
}

// SelectCable вибирає підходящий кабель за заданими параметрами
func SelectCable(ratedCurrent, length, voltage, cosFi, ambientTemp, soilThermalRes float64,
	installMethod string, cableCount int, thermalImpulse, dynamicForce float64) (CableData, bool) {

	for _, cable := range CableDatabase {
		// Визначення допустимого навантаження з урахуванням коефіцієнтів
		permissibleLoad := ApplyLoadFactors(cable.Current, ambientTemp, soilThermalRes, installMethod, cableCount)

		// Перевірка за номінальним струмом
		if permissibleLoad < ratedCurrent {
			continue
		}

		// Перевірка на падіння напруги
		voltageDropPerc := CalculateVoltageDropPercent(ratedCurrent, cable.Resistance, cable.Reactance, length, voltage, cosFi)
		if voltageDropPerc > 5.0 { // Припускаємо, що допустиме падіння напруги - 5%
			continue
		}

		// Перевірка на термічну стійкість
		minThermalSection := math.Sqrt(thermalImpulse / cable.ThermalConst)
		if cable.Section < minThermalSection {
			continue
		}

		// Перевірка на динамічну стійкість
		if dynamicForce > cable.MaxLoadForce {
			continue
		}

		// Якщо кабель відповідає всім вимогам, повертаємо його
		return cable, true
	}

	// Якщо не знайдено жодного підходящого кабелю, повертаємо пустий об'єкт і false
	return CableData{}, false
}
