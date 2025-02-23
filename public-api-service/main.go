package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"regexp"
)

// PublicAPIService struct
type PublicAPIService struct{}

// RegistrationRequest holds user registration data
type RegistrationRequest struct {
	Username        string
	Password        string
	ConfirmPassword string
	Email           string
	Age             int
}

// RegistrationResponse provides the result of registration
type RegistrationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// LoginRequest contains login credentials
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse contains the JWT token
type LoginResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
}

// validateRegistration performs input validation
func validateRegistration(req RegistrationRequest) string {
	if req.Username == "" {
		return "Username cannot be empty"
	}
	if req.Password == "" {
		return "Password cannot be empty"
	}
	if req.Password != req.ConfirmPassword {
		return "Passwords do not match"
	}
	if !isValidEmail(req.Email) {
		return "Invalid email format"
	}
	if req.Age <= 0 {
		return "Age must be a positive number"
	}
	return ""
}

// isValidEmail validates the email pattern
func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// RegisterUser calls the User Service to register a new user
func (p *PublicAPIService) RegisterUser(req RegistrationRequest, res *RegistrationResponse) error {
	// Perform input validation
	validationError := validateRegistration(req)
	if validationError != "" {
		res.Success = false
		res.Message = validationError
		return nil
	}

	// Connect to User Service
	client, err := rpc.Dial("tcp", "localhost:5000")
	if err != nil {
		return fmt.Errorf("Error connecting to User Service: %v", err)
	}
	defer client.Close()

	// Call User Service RegisterUser
	err = client.Call("UserService.RegisterUser", req, res)
	if err != nil {
		return fmt.Errorf("Error calling RegisterUser: %v", err)
	}

	return nil
}

// LoginUser calls the Authentication Service to validate user credentials and return a JWT token
func (p *PublicAPIService) LoginUser(req LoginRequest, res *LoginResponse) error {
	// Connect to Auth Service
	client, err := rpc.Dial("tcp", "localhost:5001")
	if err != nil {
		return fmt.Errorf("Error connecting to Auth Service: %v", err)
	}
	defer client.Close()

	// Call AuthService.Login
	err = client.Call("AuthService.Login", req, res)
	if err != nil {
		return fmt.Errorf("Error calling Login: %v", err)
	}

	return nil
}

func main() {
	publicAPI := new(PublicAPIService)
	err := rpc.Register(publicAPI)
	if err != nil {
		log.Fatal("Error registering PublicAPIService:", err)
	}

	listener, err := net.Listen("tcp", ":5004")
	if err != nil {
		log.Fatal("Error starting Public API Service:", err)
	}
	defer listener.Close()

	fmt.Println("Public API Service is running on port 5004...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Connection error:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
