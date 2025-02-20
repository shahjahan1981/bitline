package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

// PublicAPIService struct
type PublicAPIService struct{}

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

// LoginResponse contains the JWT token
type LoginResponse struct {
	Token string
}

// RegisterUser calls the User Service to register a new user
func (p *PublicAPIService) RegisterUser(req RegistrationRequest, res *RegistrationResponse) error {
	client, err := rpc.Dial("tcp", "localhost:5000")
	if err != nil {
		return fmt.Errorf("Error connecting to User Service: %v", err)
	}
	defer client.Close()

	err = client.Call("UserService.RegisterUser", req, res)
	if err != nil {
		return fmt.Errorf("Error calling RegisterUser: %v", err)
	}

	return nil
}

// LoginUser calls the Authentication Service to validate user credentials and return a JWT token
func (p *PublicAPIService) LoginUser(req LoginRequest, res *LoginResponse) error {
	client, err := rpc.Dial("tcp", "localhost:5001")
	if err != nil {
		return fmt.Errorf("Error connecting to Auth Service: %v", err)
	}
	defer client.Close()

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
