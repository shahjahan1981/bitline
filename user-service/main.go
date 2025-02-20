package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

// UserService handles user-related operations
type UserService struct {
	users map[string]string // stores username-password pairs
}

// RegistrationRequest holds user registration data
type RegistrationRequest struct {
	Username string
	Password string
}

// RegistrationResponse provides the result of registration
type RegistrationResponse struct {
	Message string
}

// LoginRequest contains login credentials
type LoginRequest struct {
	Username string
	Password string
}

// ValidateUserResponse provides the validation result
type ValidateUserResponse struct {
	Message string
}

// RegisterUser stores user credentials
func (u *UserService) RegisterUser(req RegistrationRequest, res *RegistrationResponse) error {
	if _, exists := u.users[req.Username]; exists {
		res.Message = "User already exists"
		return nil
	}

	u.users[req.Username] = req.Password
	res.Message = "Registration successful"
	return nil
}

// ValidateUser checks if the username and password match
func (u *UserService) ValidateUser(req LoginRequest, res *ValidateUserResponse) error {
	if password, exists := u.users[req.Username]; exists && password == req.Password {
		res.Message = "Valid"
	} else {
		res.Message = "Invalid credentials"
	}
	return nil
}

func main() {
	userService := &UserService{users: make(map[string]string)}
	err := rpc.Register(userService)
	if err != nil {
		log.Fatal("Error registering UserService:", err)
	}

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
