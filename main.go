package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	User    *User  `json:"user,omitempty"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

var db *sql.DB

func main() {
	var err error

	// Initialize database
	db, err = sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create users table
	createTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		email TEXT,
		role TEXT DEFAULT 'user'
	);`

	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}

	// Insert sample users
	insertUsers := `
	INSERT OR IGNORE INTO users (username, password, email, role) VALUES
	('admin', 'admin123', 'admin@example.com', 'admin'),
	('user', 'user123', 'user@example.com', 'user'),
	('john', 'john123', 'john@example.com', 'user');`

	_, err = db.Exec(insertUsers)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("POST /login", loginHandler)
	http.HandleFunc("POST /login-secure", loginSecureHandler)

	fmt.Println("Server started at http://localhost:8080")
	fmt.Println("\nEndpoints:")
	fmt.Println("  POST /login        - VULNERABLE to SQL Injection")
	fmt.Println("  POST /login-secure - SECURE version (prepared statements)")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var loginReq LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	// VULNERABLE: SQL Injection - menggunakan string concatenation langsung
	query := fmt.Sprintf("SELECT id, username, email, role FROM users WHERE username='%s' AND password='%s'",
		loginReq.Username, loginReq.Password)

	// Log query untuk debugging (menunjukkan vulnerability)
	log.Printf("Executing query: %s", query)

	row := db.QueryRow(query)

	var user User
	err = row.Scan(&user.ID, &user.Username, &user.Email, &user.Role)

	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid username or password",
		})
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Database error: " + err.Error(),
		})
		return
	}

	// Login successful
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "Login successful",
		User:    &user,
	})
}

func loginSecureHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var loginReq LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	// SECURE: Menggunakan prepared statement dengan parameterized query
	query := "SELECT id, username, email, role FROM users WHERE username=? AND password=?"

	log.Printf("Executing secure query with parameters: username=%s", loginReq.Username)

	row := db.QueryRow(query, loginReq.Username, loginReq.Password)

	var user User
	err = row.Scan(&user.ID, &user.Username, &user.Email, &user.Role)

	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid username or password",
		})
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Database error: " + err.Error(),
		})
		return
	}

	// Login successful
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "Login successful (SECURE)",
		User:    &user,
	})
}
