package main

import (
	"fmt"
	"net"
	"net/rpc"
	"time"
)

// PublicAPIService struct
type PublicAPIService struct{}

// RequestData method interacts with the Logical Layer
func (p *PublicAPIService) RequestData(request string, response *string) error {
	fmt.Println("Public API → Authenticating request...")

	// Authenticate before processing
	authResp, err := callAuthService()
	if err != nil {
		return fmt.Errorf("Authentication failed: %v", err)
	}
	fmt.Println("Authentication Successful:", authResp)

	// Send request to Logical Layer Service
	fmt.Println("Public API → Connecting to Logical Layer Service...")
	client, err := rpc.Dial("tcp", "localhost:5008")
	if err != nil {
		return fmt.Errorf("Error connecting to Logical Layer Service: %v", err)
	}
	defer client.Close()
	fmt.Println("Public API → Connected to Logical Layer Service.")

	// Process request via Logical Layer
	err = client.Call("LogicalLayerService.ProcessData", request, response)
	if err != nil {
		return fmt.Errorf("Error calling LogicalLayerService.ProcessData: %v", err)
	}
	fmt.Println("Public API → Response from Logical Layer:", *response)

	return nil
}

// startPublicAPIService initializes and runs the Public API Service
func startPublicAPIService() {
	service := new(PublicAPIService)
	rpc.Register(service)

	listener, err := net.Listen("tcp", ":5004")
	if err != nil {
		fmt.Println("Error starting Public API Service:", err)
		return
	}
	fmt.Println("Public API Service running on port 5004...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}

// callAuthService verifies API authentication
func callAuthService() (string, error) {
	fmt.Println("Public API → Connecting to Authentication Service...")

	client, err := rpc.Dial("tcp", "localhost:5001")
	if err != nil {
		return "", fmt.Errorf("Error connecting to Auth Service: %v", err)
	}
	defer client.Close()
	fmt.Println("Public API → Connected to Authentication Service.")

	var response string
	err = client.Call("AuthService.Authenticate", "Public API Request", &response)
	if err != nil {
		return "", fmt.Errorf("Error calling AuthService: %v", err)
	}
	return response, nil
}

// callLogicalLayerService sends a request to Logical Layer for data processing
func callLogicalLayerService() {
	time.Sleep(3 * time.Second) // Wait for Logical Layer Service to start
	fmt.Println("Public API → Calling Logical Layer Service to process data...")

	var response string
	publicAPI := new(PublicAPIService)
	err := publicAPI.RequestData("Public API Request Data", &response)
	if err != nil {
		fmt.Println("Error processing data through Logical Layer:", err)
	} else {
		fmt.Println("Public API → Response from Logical Layer Service:", response)
	}
}

func main() {
	go startPublicAPIService() // Start Public API Service
	time.Sleep(2 * time.Second)

	// Call test connections
	callAuthService()          // Authenticate API request
	callLogicalLayerService()  // Send request to Logical Layer

	select {} // Keep the service running
}
