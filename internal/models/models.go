package models

// CableSelectionData - дані для вибору кабелів
type CableSelectionData struct {
	Sn             float64 `json:"sn"`             // Номінальна потужність трансформатора, кВА
	Un             float64 `json:"un"`             // Номінальна напруга, кВ
	CosFi          float64 `json:"cosFi"`          // Коефіцієнт потужності
	Length         float64 `json:"length"`         // Довжина кабельної лінії, км
	InstallMethod  string  `json:"installMethod"`  // Спосіб прокладки кабелю
	CableCount     int     `json:"cableCount"`     // Кількість кабелів
	AmbientTemp    float64 `json:"ambientTemp"`    // Температура навколишнього середовища, °C
	SoilThermalRes float64 `json:"soilThermalRes"` // Питомий тепловий опір ґрунту
}

// SCCurrentCalcData - дані для розрахунку струмів КЗ
type SCCurrentCalcData struct {
	SystemVoltage float64 `json:"systemVoltage"` // Напруга системи, кВ
	SystemPower   float64 `json:"systemPower"`   // Потужність системи, МВА
	Xc            float64 `json:"xc"`            // Індуктивний опір системи, Ом
	Rc            float64 `json:"rc"`            // Активний опір системи, Ом
	L1            float64 `json:"l1"`            // Довжина лінії 1, км
	X1            float64 `json:"x1"`            // Питомий індуктивний опір лінії 1, Ом/км
	R1            float64 `json:"r1"`            // Питомий активний опір лінії 1, Ом/км
	L2            float64 `json:"l2"`            // Довжина лінії 2, км
	X2            float64 `json:"x2"`            // Питомий індуктивний опір лінії 2, Ом/км
	R2            float64 `json:"r2"`            // Питомий активний опір лінії 2, Ом/км
}

// EMSCCurrentCalcData - дані для розрахунку струмів КЗ для підстанції ХПнЕМ
type EMSCCurrentCalcData struct {
	SystemVoltage   float64 `json:"systemVoltage"`   // Напруга системи, кВ
	NormalModeSc    float64 `json:"normalModeSc"`    // Потужність КЗ в нормальному режимі, МВА
	MinModeSc       float64 `json:"minModeSc"`       // Потужність КЗ в мінімальному режимі, МВА
	EmergencyModeSc float64 `json:"emergencyModeSc"` // Потужність КЗ в аварійному режимі, МВА
}

// CableSelectionResult - результати вибору кабелів
type CableSelectionResult struct {
	CableType       string  `json:"cableType"`       // Тип кабелю
	CableSection    float64 `json:"cableSection"`    // Переріз кабелю, мм²
	RatedCurrent    float64 `json:"ratedCurrent"`    // Номінальний струм, А
	PermissibleLoad float64 `json:"permissibleLoad"` // Допустиме навантаження, А
	VoltageDropPerc float64 `json:"voltageDropPerc"` // Відсоток падіння напруги, %
	IsValid         bool    `json:"isValid"`         // Чи відповідає кабель вимогам
}

// SCCurrentResult - результати розрахунку струмів КЗ
type SCCurrentResult struct {
	ThreePhaseCurrent     float64 `json:"threePhaseCurrent"`     // Струм трифазного КЗ, кА
	SinglePhaseCurrent    float64 `json:"singlePhaseCurrent"`    // Струм однофазного КЗ, кА
	ImpulseCurrentValue   float64 `json:"impulseCurrentValue"`   // Значення ударного струму, кА
	ThermalImpulse        float64 `json:"thermalImpulse"`        // Тепловий імпульс, кА²·с
	ThermalStability      bool    `json:"thermalStability"`      // Термічна стійкість (так/ні)
	DynamicStability      bool    `json:"dynamicStability"`      // Динамічна стійкість (так/ні)
	ThermalStabilityMin   float64 `json:"thermalStabilityMin"`   // Мінімальний переріз для термічної стійкості, мм²
	DynamicStabilityForce float64 `json:"dynamicStabilityForce"` // Сила при КЗ, Н
}

// EMSCCurrentResult - результати розрахунку струмів КЗ для підстанції ХПнЕМ
type EMSCCurrentResult struct {
	NormalModeThreePhaseSC     float64 `json:"normalModeThreePhaseSC"`     // Струм трифазного КЗ в нормальному режимі, кА
	NormalModeSinglePhaseSC    float64 `json:"normalModeSinglePhaseSC"`    // Струм однофазного КЗ в нормальному режимі, кА
	MinModeThreePhaseSC        float64 `json:"minModeThreePhaseSC"`        // Струм трифазного КЗ в мінімальному режимі, кА
	MinModeSinglePhaseSC       float64 `json:"minModeSinglePhaseSC"`       // Струм однофазного КЗ в мінімальному режимі, кА
	EmergencyModeThreePhaseSC  float64 `json:"emergencyModeThreePhaseSC"`  // Струм трифазного КЗ в аварійному режимі, кА
	EmergencyModeSinglePhaseSC float64 `json:"emergencyModeSinglePhaseSC"` // Струм однофазного КЗ в аварійному режимі, кА
}
