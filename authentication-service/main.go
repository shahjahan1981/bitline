package main

import (
	"fmt"
	"net"
	"net/rpc"
	"time"
)

// AuthService defines authentication-related methods
type AuthService struct{}

// Authenticate validates user credentials
func (a *AuthService) Authenticate(request string, response *string) error {
	*response = "Authentication successful for user: " + request
	return nil
}

// startServer initializes and starts the Authentication Service
func startServer() {
	service := new(AuthService)
	rpc.Register(service)

	listener, err := net.Listen("tcp", ":5001")
	if err != nil {
		fmt.Println("Error starting Authentication Service:", err)
		return
	}
	fmt.Println("Authentication Service is running on port 5001...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}

// callUserService calls the User Service to fetch user details
func callUserService() {
	time.Sleep(3 * time.Second) // Give User Service time to start
	client, err := rpc.Dial("tcp", "localhost:5000") // Connect to User Service
	if err != nil {
		fmt.Println("Error connecting to User Service:", err)
		return
	}
	defer client.Close()

	var response string
	err = client.Call("UserService.GetUser", "shahjahan", &response)
	if err != nil {
		fmt.Println("Error calling UserService:", err)
	} else {
		fmt.Println("Response from User Service:", response)
	}
}

func main() {
	go startServer() // Start the Authentication Service in a goroutine
	time.Sleep(2 * time.Second)
	callUserService() // Make an RPC call to User Service
	select {}         // Keep the service running
}
