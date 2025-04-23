package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/harip/GoTasker/handlers"
	"github.com/harip/GoTasker/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database: " + err.Error())
	}
	db.AutoMigrate(&models.Task{})
	handlers.InitDB(db)
	return db
}

func TestCreateTask(t *testing.T) {
	db := setupTestDB()
	defer db.Migrator().DropTable(&models.Task{})

	task := models.Task{
		Title:       "Test Task",
		Description: "Test Description",
		Status:      "Pending",
	}
	body, _ := json.Marshal(task)

	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(handlers.CreateTask)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status %v, got %v", http.StatusCreated, rr.Code)
	}

	var createdTask models.Task
	json.Unmarshal(rr.Body.Bytes(), &createdTask)
	if createdTask.Title != task.Title {
		t.Errorf("Expected title %v, got %v", task.Title, createdTask.Title)
	}
}

func TestGetTasks(t *testing.T) {
	db := setupTestDB()
	defer db.Migrator().DropTable(&models.Task{})

	db.Create(&models.Task{Title: "Task 1", Status: "Pending"})
	db.Create(&models.Task{Title: "Task 2", Status: "In Progress"})

	req, _ := http.NewRequest("GET", "/tasks?page=1&limit=10", nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(handlers.GetTasks)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %v, got %v", http.StatusOK, rr.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &response)
	tasks := response["tasks"].([]interface{})
	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %v", len(tasks))
	}
}

func TestGetTaskByID(t *testing.T) {
	db := setupTestDB()
	defer db.Migrator().DropTable(&models.Task{})

	task := models.Task{Title: "Test Task", Status: "Pending"}
	db.Create(&task)

	req, _ := http.NewRequest("GET", "/tasks/1", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/tasks/{id}", handlers.GetTaskByID).Methods("GET")
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %v, got %v", http.StatusOK, rr.Code)
	}

	var fetchedTask models.Task
	json.Unmarshal(rr.Body.Bytes(), &fetchedTask)
	if fetchedTask.ID != task.ID {
		t.Errorf("Expected ID %v, got %v", task.ID, fetchedTask.ID)
	}
}
