package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"path/filepath"

	"github.com/TexnoViktor/short-circuit-calculator/internal/calculations"
	"github.com/TexnoViktor/short-circuit-calculator/internal/models"
)

// Server представляє структуру сервера
type Server struct {
	Templates *template.Template
}

// NewServer створює новий екземпляр сервера
func NewServer() (*Server, error) {
	// Завантаження шаблонів
	templates, err := LoadTemplates(filepath.Join("web", "templates", "*.html"))
	if err != nil {
		return nil, fmt.Errorf("error loading templates: %v", err)
	}

	server := &Server{
		Templates: templates,
	}

	return server, nil
}

// HomeHandler обробляє запити до головної сторінки
func (s *Server) HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if err := s.Templates.ExecuteTemplate(w, "index.html", nil); err != nil {
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
		return
	}
}

// CableSelectionHandler обробляє запити на вибір кабелів
func (s *Server) CableSelectionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data models.CableSelectionData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding request body: %v", err), http.StatusBadRequest)
		return
	}

	// Обчислення номінального струму
	ratedCurrent := calculations.CalculateRatedCurrent(data.Sn, data.Un, data.CosFi)

	// Для спрощення прикладу, будемо вважати, що у нас є попередньо розрахований тепловий імпульс та сила для динамічної перевірки
	// В реальному додатку ці значення повинні бути розраховані на основі даних короткого замикання
	thermalImpulse := 20.0 // кА²·с
	dynamicForce := 1000.0 // Н

	// Вибір кабелю
	selectedCable, isValid := calculations.SelectCable(
		ratedCurrent, data.Length, data.Un, data.CosFi,
		data.AmbientTemp, data.SoilThermalRes, data.InstallMethod,
		data.CableCount, thermalImpulse, dynamicForce,
	)

	// Розрахунок падіння напруги для вибраного кабелю
	voltageDropPerc := 0.0
	if isValid {
		voltageDropPerc = calculations.CalculateVoltageDropPercent(
			ratedCurrent, selectedCable.Resistance, selectedCable.Reactance,
			data.Length, data.Un, data.CosFi,
		)
	}

	// Формування результату
	result := models.CableSelectionResult{
		CableType:       selectedCable.Type,
		CableSection:    selectedCable.Section,
		RatedCurrent:    ratedCurrent,
		PermissibleLoad: selectedCable.Current,
		VoltageDropPerc: voltageDropPerc,
		IsValid:         isValid,
	}

	// Відправка результату
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
}

// SCCurrentsHandler обробляє запити на розрахунок струмів КЗ
func (s *Server) SCCurrentsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data models.SCCurrentCalcData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding request body: %v", err), http.StatusBadRequest)
		return
	}

	// Розрахунок загального активного та індуктивного опору
	totalR := data.Rc + data.R1*data.L1 + data.R2*data.L2
	totalX := data.Xc + data.X1*data.L1 + data.X2*data.L2

	// Розрахунок струму трифазного КЗ
	threePhaseCurrent := calculations.CalculateThreePhaseSCCurrent(data.SystemVoltage, totalR, totalX)

	// Розрахунок опорів нульової послідовності для розрахунку однофазного КЗ
	r0, x0 := calculations.CalculateZeroSequenceImpedance(totalR, totalX)

	// Розрахунок струму однофазного КЗ
	singlePhaseCurrent := calculations.CalculateSinglePhaseSCCurrent(data.SystemVoltage, totalR, totalX, r0, x0)

	// Розрахунок ударного струму
	impulseCurrentValue := calculations.CalculateImpulseCurrent(threePhaseCurrent, totalR, totalX)

	// Розрахунок теплового імпульсу (припускаємо час відключення 0.1 с)
	disconnectTime := 0.1
	thermalImpulse := calculations.CalculateThermalImpulse(threePhaseCurrent, totalR, totalX, disconnectTime)

	// Перевірка термічної стійкості (припускаємо мінімальний переріз 70 мм²)
	thermalStabilityMin := calculations.CalculateThermalStabilityAluminum(thermalImpulse)
	thermalStability := calculations.CheckThermalStability(70.0, thermalStabilityMin)

	// Розрахунок сили для перевірки динамічної стійкості (припускаємо відстань між провідниками 0.2 м)
	dynamicStabilityForce := calculations.CalculateDynamicStabilityForce(impulseCurrentValue, 0.2)
	dynamicStability := calculations.CheckDynamicStability(dynamicStabilityForce, 2000.0) // 2000 Н - умовне максимальне значення

	// Формування результату
	result := models.SCCurrentResult{
		ThreePhaseCurrent:     threePhaseCurrent,
		SinglePhaseCurrent:    singlePhaseCurrent,
		ImpulseCurrentValue:   impulseCurrentValue,
		ThermalImpulse:        thermalImpulse,
		ThermalStability:      thermalStability,
		DynamicStability:      dynamicStability,
		ThermalStabilityMin:   thermalStabilityMin,
		DynamicStabilityForce: dynamicStabilityForce,
	}

	// Відправка результату
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
}

// EMSCCurrentsHandler обробляє запити на розрахунок струмів КЗ для підстанції ХПнЕМ
func (s *Server) EMSCCurrentsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data models.EMSCCurrentCalcData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding request body: %v", err), http.StatusBadRequest)
		return
	}

	// Розрахунок струмів КЗ для нормального режиму
	normalModeThreePhaseSC := data.NormalModeSc / (math.Sqrt(3) * data.SystemVoltage)
	normalModeSinglePhaseSC := calculations.CalculateSinglePhaseCurrentForEMStation(normalModeThreePhaseSC)

	// Розрахунок струмів КЗ для мінімального режиму
	minModeThreePhaseSC := data.MinModeSc / (math.Sqrt(3) * data.SystemVoltage)
	minModeSinglePhaseSC := calculations.CalculateSinglePhaseCurrentForEMStation(minModeThreePhaseSC)

	// Розрахунок струмів КЗ для аварійного режиму
	emergencyModeThreePhaseSC := data.EmergencyModeSc / (math.Sqrt(3) * data.SystemVoltage)
	emergencyModeSinglePhaseSC := calculations.CalculateSinglePhaseCurrentForEMStation(emergencyModeThreePhaseSC)

	// Формування результату
	result := models.EMSCCurrentResult{
		NormalModeThreePhaseSC:     normalModeThreePhaseSC,
		NormalModeSinglePhaseSC:    normalModeSinglePhaseSC,
		MinModeThreePhaseSC:        minModeThreePhaseSC,
		MinModeSinglePhaseSC:       minModeSinglePhaseSC,
		EmergencyModeThreePhaseSC:  emergencyModeThreePhaseSC,
		EmergencyModeSinglePhaseSC: emergencyModeSinglePhaseSC,
	}

	// Відправка результату
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
}
