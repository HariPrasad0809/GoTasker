package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/harip/GoTasker/models"
)

func CreateTask(w http.ResponseWriter, r *http.Request) {
	if !IsDBInitialized() {
		log.Println("Error: Database not initialized")
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	userID, ok := r.Context().Value("user_id").(float64)
	if !ok {
		log.Println("Error: User ID not found in context")
		http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var input struct {
		Title       string     `json:"title"`
		Description string     `json:"description"`
		Status      string     `json:"status"`
		DueDate     *time.Time `json:"due_date"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, `{"error": "Invalid request body: `+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	if input.Title == "" {
		log.Printf("Invalid task data: Title is empty")
		http.Error(w, `{"error": "Title is required"}`, http.StatusBadRequest)
		return
	}
	if input.Status == "" {
		input.Status = "Pending"
	} else if !isValidStatus(input.Status) {
		log.Printf("Invalid task data: Status=%s", input.Status)
		http.Error(w, `{"error": "Status must be Pending, In Progress, or Completed"}`, http.StatusBadRequest)
		return
	}

	task := models.Task{
		Title:       input.Title,
		Description: input.Description,
		Status:      input.Status,
		DueDate:     input.DueDate,
		UserID:      int(userID),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := db.Create(&task).Error; err != nil {
		log.Printf("Error creating task for user_id %d: %v", int(userID), err)
		http.Error(w, `{"error": "Failed to create task: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	log.Printf("Task created successfully for user_id %d: ID=%d, Title=%s", int(userID), task.ID, task.Title)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	if !IsDBInitialized() {
		log.Println("Error: Database not initialized")
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	userID, ok := r.Context().Value("user_id").(float64)
	if !ok {
		log.Println("Error: User ID not found in context")
		http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit < 1 || limit > 50 {
		limit = 10
	}
	offset := (page - 1) * limit

	status := query.Get("status")
	dueDateAfter := query.Get("due_date_after")
	dueDateBefore := query.Get("due_date_before")
	sortBy := query.Get("sort_by")
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortOrder := strings.ToLower(query.Get("sort_order"))
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	var tasks []models.Task
	dbQuery := db.Model(&models.Task{}).Where("user_id = ?", int(userID))

	if status != "" && isValidStatus(status) {
		dbQuery = dbQuery.Where("status = ?", status)
	}
	if dueDateAfter != "" {
		if t, err := time.Parse(time.RFC3339, dueDateAfter); err == nil {
			dbQuery = dbQuery.Where("due_date > ?", t)
		}
	}
	if dueDateBefore != "" {
		if t, err := time.Parse(time.RFC3339, dueDateBefore); err == nil {
			dbQuery = dbQuery.Where("due_date < ?", t)
		}
	}

	var total int64
	dbQuery.Count(&total)

	dbQuery.Order(sortBy + " " + sortOrder).Offset(offset).Limit(limit).Find(&tasks)

	response := map[string]interface{}{
		"tasks": tasks,
		"page":  page,
		"limit": limit,
		"total": total,
	}

	log.Printf("Retrieved %d tasks for user_id %d (page=%d, limit=%d)", len(tasks), int(userID), page, limit)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetTaskByID(w http.ResponseWriter, r *http.Request) {
	if !IsDBInitialized() {
		log.Println("Error: Database not initialized")
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	userID, ok := r.Context().Value("user_id").(float64)
	if !ok {
		log.Println("Error: User ID not found in context")
		http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Invalid task ID: %v", err)
		http.Error(w, `{"error": "Invalid task ID"}`, http.StatusBadRequest)
		return
	}

	var task models.Task
	if err := db.Where("id = ? AND user_id = ?", id, int(userID)).First(&task).Error; err != nil {
		log.Printf("Task not found for user_id %d: ID=%d, error=%v", int(userID), id, err)
		http.Error(w, `{"error": "Task not found"}`, http.StatusNotFound)
		return
	}

	log.Printf("Retrieved task for user_id %d: ID=%d, Title=%s", int(userID), task.ID, task.Title)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	if !IsDBInitialized() {
		log.Println("Error: Database not initialized")
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	userID, ok := r.Context().Value("user_id").(float64)
	if !ok {
		log.Println("Error: User ID not found in context")
		http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Invalid task ID: %v", err)
		http.Error(w, `{"error": "Invalid task ID"}`, http.StatusBadRequest)
		return
	}

	var task models.Task
	if err := db.Where("id = ? AND user_id = ?", id, int(userID)).First(&task).Error; err != nil {
		log.Printf("Task not found for user_id %d: ID=%d, error=%v", int(userID), id, err)
		http.Error(w, `{"error": "Task not found"}`, http.StatusNotFound)
		return
	}

	var input struct {
		Title       string     `json:"title"`
		Description string     `json:"description"`
		Status      string     `json:"status"`
		DueDate     *time.Time `json:"due_date"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, `{"error": "Invalid request body: `+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	if input.Title == "" {
		log.Printf("Invalid task data: Title is empty")
		http.Error(w, `{"error": "Title is required"}`, http.StatusBadRequest)
		return
	}
	if input.Status == "" {
		input.Status = task.Status
	} else if !isValidStatus(input.Status) {
		log.Printf("Invalid task data: Status=%s", input.Status)
		http.Error(w, `{"error": "Status must be Pending, In Progress, or Completed"}`, http.StatusBadRequest)
		return
	}

	task.Title = input.Title
	task.Description = input.Description
	task.Status = input.Status
	task.DueDate = input.DueDate
	task.UpdatedAt = time.Now()

	if err := db.Save(&task).Error; err != nil {
		log.Printf("Error updating task for user_id %d: ID=%d, error=%v", int(userID), id, err)
		http.Error(w, `{"error": "Failed to update task: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	log.Printf("Task updated successfully for user_id %d: ID=%d, Title=%s", int(userID), task.ID, task.Title)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	if !IsDBInitialized() {
		log.Println("Error: Database not initialized")
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	userID, ok := r.Context().Value("user_id").(float64)
	if !ok {
		log.Println("Error: User ID not found in context")
		http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Invalid task ID: %v", err)
		http.Error(w, `{"error": "Invalid task ID"}`, http.StatusBadRequest)
		return
	}

	if err := db.Where("id = ? AND user_id = ?", id, int(userID)).Delete(&models.Task{}).Error; err != nil {
		log.Printf("Task not found for user_id %d: ID=%d, error=%v", int(userID), id, err)
		http.Error(w, `{"error": "Task not found"}`, http.StatusNotFound)
		return
	}

	log.Printf("Task deleted successfully for user_id %d: ID=%d", int(userID), id)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Task successfully deleted"})
}

func isValidStatus(status string) bool {
	return status == "Pending" || status == "In Progress" || status == "Completed"
}
