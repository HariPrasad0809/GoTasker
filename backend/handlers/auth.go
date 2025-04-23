package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/harip/GoTasker/config"
	"github.com/harip/GoTasker/models"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	if !IsDBInitialized() {
		log.Println("Error: Database not initialized")
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}
	log.Printf("Login attempt for username: %s", creds.Username)

	var user models.User
	if err := db.Where("username = ?", creds.Username).First(&user).Error; err != nil {
		log.Printf("User not found or query error for username '%s': %v", creds.Username, err)
		http.Error(w, `{"error": "Invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	if creds.Password != user.Password {
		hash := sha256.Sum256([]byte(creds.Password))
		inputHash := hex.EncodeToString(hash[:])
		storedHash := sha256.Sum256([]byte(user.Password))
		storedHashStr := hex.EncodeToString(storedHash[:])
		log.Printf("Password verification failed for username: %s, stored password (hashed for log): %s, input password (hashed for log): %s",
			creds.Username, storedHashStr, inputHash)
		http.Error(w, `{"error": "Invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  float64(user.ID),
		"username": creds.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		log.Printf("Error generating token: %v", err)
		http.Error(w, `{"error": "Failed to generate token"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
	log.Printf("Login successful for username: %s", creds.Username)
}

func Register(w http.ResponseWriter, r *http.Request) {
	if !IsDBInitialized() {
		log.Println("Error: Database not initialized")
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		log.Printf("Error decoding register request body: %v", err)
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}
	log.Printf("Register attempt for username: %s", creds.Username)

	if creds.Username == "" {
		log.Printf("Username is empty")
		http.Error(w, `{"error": "Username is required"}`, http.StatusBadRequest)
		return
	}
	if creds.Password == "" {
		log.Printf("Password is empty")
		http.Error(w, `{"error": "Password is required"}`, http.StatusBadRequest)
		return
	}
	if creds.Email == "" {
		log.Printf("Email is empty")
		http.Error(w, `{"error": "Email is required"}`, http.StatusBadRequest)
		return
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(creds.Email) {
		log.Printf("Invalid email format: %s", creds.Email)
		http.Error(w, `{"error": "Invalid email format"}`, http.StatusBadRequest)
		return
	}

	var existingUser models.User
	if err := db.Where("username = ?", creds.Username).First(&existingUser).Error; err == nil {
		log.Printf("Username %s already exists", creds.Username)
		http.Error(w, `{"error": "Username already taken"}`, http.StatusConflict)
		return
	}

	if err := db.Where("email = ?", creds.Email).First(&existingUser).Error; err == nil {
		log.Printf("Email %s already exists", creds.Email)
		http.Error(w, `{"error": "Email already taken"}`, http.StatusConflict)
		return
	}

	user := models.User{
		Username: creds.Username,
		Password: creds.Password,
		Email:    creds.Email,
	}
	if err := db.Create(&user).Error; err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, `{"error": "Failed to create user"}`, http.StatusInternalServerError)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  float64(user.ID),
		"username": creds.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		log.Printf("Error generating token: %v", err)
		http.Error(w, `{"error": "Failed to generate token"}`, http.StatusInternalServerError)
		return
	}

	hash := sha256.Sum256([]byte(user.Password))
	hashedPassword := hex.EncodeToString(hash[:])
	log.Printf("User %s created successfully with password (hashed for log): %s", creds.Username, hashedPassword)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
