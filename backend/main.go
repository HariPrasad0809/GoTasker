package main

import (
	"fmt"
	"log"
	"net/http"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/harip/GoTasker/config"
	"github.com/harip/GoTasker/handlers"
	"github.com/harip/GoTasker/middleware"
	"github.com/harip/GoTasker/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func migrateUserTable(db *gorm.DB) bool {
	var columnExists bool
	db.Raw(`SELECT EXISTS (
		SELECT 1 
		FROM information_schema.columns 
		WHERE table_name = 'users' 
		AND column_name = 'password'
	)`).Scan(&columnExists)

	if !columnExists {
		if err := db.Exec("UPDATE users SET password_hash = 'default_password' WHERE password_hash IS NULL").Error; err != nil {
			log.Printf("Failed to set default password: %v", err)
			return false
		}
		if err := db.Exec("ALTER TABLE users RENAME COLUMN password_hash TO password").Error; err != nil {
			log.Printf("Failed to rename column: %v", err)
			return false
		}
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Printf("Auto-migrate failed: %v", err)
		return false
	}
	return true
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: Error loading .env file, using default environment variables")
	}
	config.LoadConfig()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.AppConfig.DBHost,
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBName,
		config.AppConfig.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Connected to the database")

	if err := db.AutoMigrate(&models.User{}, &models.Task{}); err != nil || !migrateUserTable(db) {
		log.Fatalf("Auto-migration failed: %v", err)
	}
	log.Println("Database schema migrated")

	h := &Handler{DB: db}
	handlers.SetDB(h.DB)
	if !handlers.IsDBInitialized() {
		log.Fatal("Failed to initialize handlers DB")
	}
	log.Println("Handlers DB initialized")

	r := mux.NewRouter()

	r.HandleFunc("/register", handlers.Register).Methods("POST")
	r.HandleFunc("/login", handlers.Login).Methods("POST")
	r.Handle("/tasks", middleware.JWTMiddleware(http.HandlerFunc(handlers.CreateTask))).Methods("POST")
	r.Handle("/tasks", middleware.JWTMiddleware(http.HandlerFunc(handlers.GetTasks))).Methods("GET")
	r.Handle("/tasks/{id}", middleware.JWTMiddleware(http.HandlerFunc(handlers.GetTaskByID))).Methods("GET")
	r.Handle("/tasks/{id}", middleware.JWTMiddleware(http.HandlerFunc(handlers.UpdateTask))).Methods("PUT")
	r.Handle("/tasks/{id}", middleware.JWTMiddleware(http.HandlerFunc(handlers.DeleteTask))).Methods("DELETE")

	cors := gorillaHandlers.CORS(
		gorillaHandlers.AllowedOrigins([]string{"http://localhost:3000"}),
		gorillaHandlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		gorillaHandlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", cors(r)))
}
