package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"regexp"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// Database connection string
const dbConnStr = "host=localhost port=5433 user=postgres password=198181018 dbname=go_api_db sslmode=disable"

// UserService handles user-related operations
type UserService struct {
	db *sql.DB
}

// RegistrationRequest holds user registration data
type RegistrationRequest struct {
	FullName        string `json:"full_name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

// RegistrationResponse provides the result of registration
type RegistrationResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// Email validation regex pattern
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// RegisterUser stores user credentials and sends a welcome email
func (u *UserService) RegisterUser(req RegistrationRequest, res *RegistrationResponse) error {
	// Validate required fields
	if req.FullName == "" || req.Password == "" || req.ConfirmPassword == "" || req.Email == "" {
		res.Status = "error"
		res.Message = "All fields are required"
		return nil
	}

	// Validate password match
	if req.Password != req.ConfirmPassword {
		res.Status = "error"
		res.Message = "Passwords do not match"
		return nil
	}

	// Validate email pattern
	if !emailRegex.MatchString(req.Email) {
		res.Status = "error"
		res.Message = "Invalid email format"
		return nil
	}

	// Check if user already exists
	var exists bool
	err := u.db.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE email=$1)", req.Email).Scan(&exists)
	if err != nil {
		res.Status = "error"
		res.Message = "Database error"
		return err
	}
	if exists {
		res.Status = "error"
		res.Message = "User already exists"
		return nil
	}

	// Hash the password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		res.Status = "error"
		res.Message = "Failed to hash password"
		return err
	}

	// Insert user into database
	_, err = u.db.Exec("INSERT INTO users (full_name, email, password_hash, created_at) VALUES ($1, $2, $3, $4)",
		req.FullName, req.Email, string(hashedPassword), time.Now())
	if err != nil {
		res.Status = "error"
		res.Message = "Failed to register user"
		return err
	}

	// Generate welcome email message
	emailMessage := fmt.Sprintf("Welcome %s! Your account has been created successfully.", req.FullName)

	// Save welcome email to database
	_, err = u.db.Exec("INSERT INTO emails (full_name, email, message, timestamp) VALUES ($1, $2, $3, $4)",
		req.FullName, req.Email, emailMessage, time.Now())
	if err != nil {
		res.Status = "error"
		res.Message = "Failed to save welcome email"
		return err
	}

	// Successful response
	res.Status = "success"
	res.Message = "Registration successful"
	return nil
}

func main() {
	// Connect to PostgreSQL
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatal("Database connection error:", err)
	}
	defer db.Close()

	// Verify DB connection
	if err = db.Ping(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Create tables if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			full_name VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
			password_hash TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS emails (
			id SERIAL PRIMARY KEY,
			full_name VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL,
			message TEXT NOT NULL,
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		log.Fatal("Error creating tables:", err)
	}

	// Register UserService with database
	userService := &UserService{db: db}
	err = rpc.Register(userService)
	if err != nil {
		log.Fatal("Error registering UserService:", err)
	}

	// Start User Service on port 5000
	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatal("Error starting User Service:", err)
	}
	defer listener.Close()

	fmt.Println("User Service is running on port 5000...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Connection error:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
