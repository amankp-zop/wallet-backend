package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/amankp-zop/wallet/internal/api/handler"
	"github.com/amankp-zop/wallet/internal/config"
	"github.com/amankp-zop/wallet/internal/database"
	"github.com/amankp-zop/wallet/internal/repository"
	"github.com/amankp-zop/wallet/internal/service"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	cfg, err := config.LoadConfig("./configs")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	db, err := database.NewDatabase(cfg.Database.DSN)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	fmt.Println("Database connected Successfully.")
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	//router.Use(middleware.Recoverer)

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		response := struct {
			Status string `json:"status"`
		}{
			Status: "ok",
		}

		w.Header().Set("Content-Type", "application/json")

		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Printf("Error encoding health check response: %v", err)
		}
	})

	router.Route("/users", func(r chi.Router) {
		fmt.Printf("Calling users route\n")
		r.Post("/signup", userHandler.Signup)
	})

	log.Printf("Starting server on port %s", cfg.Server.Port)

	err = http.ListenAndServe(":"+cfg.Server.Port, router)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
