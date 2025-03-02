package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AuthService handles authentication
type AuthService struct{}

// LoginRequest contains login credentials
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse contains the JWT token or error message
type LoginResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token,omitempty"`
	Error   string `json:"error,omitempty"`
}

// ValidateUserResponse from User Service
type ValidateUserResponse struct {
	Message string `json:"message"`
}

// Secret key for JWT signing
var jwtSecret = []byte("supersecretkey")

// GenerateJWT creates a JWT token for authentication
func GenerateJWT(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 1).Unix(), // Token expires in 1 hour
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// Login validates credentials using User Service and returns JSON response
func (a *AuthService) Login(req LoginRequest, res *LoginResponse) error {
	client, err := rpc.Dial("tcp", "localhost:5000") // Connect to User Service
	if err != nil {
		res.Success = false
		res.Error = "Error connecting to User Service"
		return nil
	}
	defer client.Close()

	// Call UserService.ValidateUser
	var validationRes ValidateUserResponse
	err = client.Call("UserService.ValidateUser", req, &validationRes)
	if err != nil {
		res.Success = false
		res.Error = "Error calling ValidateUser"
		return nil
	}

	// Check if credentials are valid
	if validationRes.Message == "Valid" {
		token, err := GenerateJWT(req.Username)
		if err != nil {
			res.Success = false
			res.Error = "Error generating token"
			return nil
		}
		res.Success = true
		res.Token = token
	} else {
		res.Success = false
		res.Error = "Invalid credentials"
	}
	return nil
}

func main() {
	authService := new(AuthService)
	err := rpc.Register(authService)
	if err != nil {
		log.Fatal("Error registering AuthService:", err)
	}

	listener, err := net.Listen("tcp", ":5001")
	if err != nil {
		log.Fatal("Error starting Authentication Service:", err)
	}
	defer listener.Close()

	fmt.Println("Authentication Service is running on port 5001...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Connection error:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
