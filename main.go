package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/yourusername/short-circuit-calculator/internal/handlers"
)

func main() {
	// Ініціалізація сервера
	server, err := handlers.NewServer()
	if err != nil {
		log.Fatalf("Помилка при створенні сервера: %v", err)
	}

	// Створення маршрутизатора
	r := chi.NewRouter()

	// Проміжне ПЗ
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Статичні файли
	staticDir := http.Dir(filepath.Join("web", "static"))
	fileServer := http.FileServer(staticDir)
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	// Маршрути
	r.Get("/", server.HomeHandler)
	r.Post("/api/cable-selection", server.CableSelectionHandler)
	r.Post("/api/sc-currents", server.SCCurrentsHandler)
	r.Post("/api/em-sc-currents", server.EMSCCurrentsHandler)

	// Запуск сервера
	port := "8080"
	log.Printf("Сервер запущено на порту %s...", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Помилка при запуску сервера: %v", err)
	}
}
