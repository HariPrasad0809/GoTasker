GoTasker
GoTasker is a RESTful API built in Go for managing tasks with CRUD operations, JWT authentication, rate-limiting, caching, and PostgreSQL storage. It includes a simple React.js frontend for testing.
Prerequisites

Go 1.21+
PostgreSQL
Docker (optional)
Node.js (for frontend)

Setup Instructions
Backend

Clone the repository:git clone <repository-url>
cd GoTasker


Install dependencies:go mod tidy


Set up PostgreSQL:
Create a database: psql -U postgres -c "CREATE DATABASE taskdb;"
Run migrations:export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING="host=localhost port=5432 user=postgres password=yourpassword dbname=taskdb sslmode=disable"
goose -dir migrations up




Set environment variables:
DB_URL: PostgreSQL connection string (e.g., postgresql://postgres:yourpassword@localhost:5432/taskdb)
JWT_SECRET: Secret key for JWT (e.g., your-secret-key)


Run the backend:go run main.go



Frontend

Navigate to the frontend directory:cd GoTasker-frontend


Install dependencies:npm install


Run the frontend:npm start



API Endpoints
Authentication

POST /login
Request: {"username": "string", "password": "string"}
Response: {"token": "jwt-token"}



Tasks

POST /tasks (Requires JWT)
Request: {"title": "string", "description": "string", "status": "Pending|In Progress|Completed", "due_date": "2025-02-25T00:00:00Z"}
Response: Task object with ID, Title, Description, Status, DueDate, CreatedAt, UpdatedAt


GET /tasks (Requires JWT)
Query Params: page, limit, status, due_date_after, due_date_before, sort_by, sort_order
Response: {"tasks": [], "page": int, "limit": int, "total": int}


GET /tasks/{id} (Requires JWT)
Response: Task object or 404


PUT /tasks/{id} (Requires JWT)
Request: Same as POST /tasks
Response: Updated task object


DELETE /tasks/{id} (Requires JWT)
Response: {"message": "Task successfully deleted"} or 404



Running Tests
go test ./tests -v

Docker Setup

Build the Docker image:docker build -t gotasker:latest .


Run the container:docker run -d -p 8080:8080 --name gotasker -e DB_URL=postgresql://postgres:yourpassword@host:5432/taskdb -e JWT_SECRET=your-secret-key gotasker:latest



Notes

Ensure PostgreSQL is running and accessible.
Use Postman or the frontend to test API endpoints.
Rate-limiting: 60 requests per minute per IP.
Caching: Task list cached for 5 minutes.

