package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/harip/GoTasker/handlers"
	"github.com/harip/GoTasker/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupAuthTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Failed to auto-migrate test database: %v", err)
	}
	handlers.InitDB(db)
	if !handlers.IsDBInitialized() {
		log.Fatal("Failed to initialize test database")
	}
	log.Println("Test database initialized successfully")
	return db
}

func TestLogin(t *testing.T) {
	db := setupAuthTestDB()
	defer db.Migrator().DropTable(&models.User{})

	user := models.User{
		Username: "testuser",
		Password: "password123", // Store plain text password
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	log.Printf("Test user created: %s with password: %s", user.Username, user.Password)

	// Test case: Invalid credentials
	creds := map[string]string{
		"username": "testuser",
		"password": "password1213",
	}
	body, err := json.Marshal(creds)
	if err != nil {
		t.Fatalf("Failed to marshal invalid creds: %v", err)
	}

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(handlers.Login)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %v for invalid credentials, got %v", http.StatusUnauthorized, rr.Code)
	}
	log.Println("Invalid credentials test passed")

	// Test case: Valid credentials
	creds = map[string]string{
		"username": "testuser",
		"password": "password123",
	}
	body, err = json.Marshal(creds)
	if err != nil {
		t.Fatalf("Failed to marshal valid creds: %v", err)
	}

	req, err = http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(handlers.Login)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %v for valid credentials, got %v", http.StatusOK, rr.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if response["token"] == "" {
		t.Error("Expected a token, got none")
	}
	log.Println("Valid credentials test passed, token received")
}
