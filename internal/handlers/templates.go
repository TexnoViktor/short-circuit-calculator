package handlers

import (
	"fmt"
	"html/template"
	"net/http"
)

// RenderTemplate відображає HTML-шаблон
func (s *Server) RenderTemplate(w http.ResponseWriter, name string, data interface{}) error {
	// Встановлення типу вмісту
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Виконання шаблону
	if err := s.Templates.ExecuteTemplate(w, name, data); err != nil {
		return fmt.Errorf("помилка при відображенні шаблону %s: %v", name, err)
	}

	return nil
}

// RenderError відображає сторінку з повідомленням про помилку
func (s *Server) RenderError(w http.ResponseWriter, status int, message string) {
	// Встановлення коду статусу
	w.WriteHeader(status)

	// Створення даних для шаблону помилки
	data := struct {
		Status  int
		Message string
	}{
		Status:  status,
		Message: message,
	}

	// Спроба відобразити шаблон помилки
	if err := s.RenderTemplate(w, "error.html", data); err != nil {
		// Якщо не вдалося відобразити шаблон, виводимо простий текст
		http.Error(w, message, status)
	}
}

// LoadTemplates завантажує всі HTML-шаблони
func LoadTemplates(dir string) (*template.Template, error) {
	// Створення нового набору шаблонів
	tmpl := template.New("")

	// Налаштування функцій, доступних у шаблонах
	tmpl = tmpl.Funcs(template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"formatFloat": func(value float64, precision int) string {
			format := fmt.Sprintf("%%.%df", precision)
			return fmt.Sprintf(format, value)
		},
		"yesNo": func(value bool) string {
			if value {
				return "Так"
			}
			return "Ні"
		},
	})

	// Завантаження шаблонів з директорії
	tmpl, err := tmpl.ParseGlob(dir)
	if err != nil {
		return nil, fmt.Errorf("помилка при завантаженні шаблонів: %v", err)
	}

	return tmpl, nil
}
