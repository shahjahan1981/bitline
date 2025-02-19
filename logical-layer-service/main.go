package main

import (
	"fmt"
	"net"
	"net/rpc"
	"time"
)

// LogicalLayerService struct
type LogicalLayerService struct{}

// ProcessData interacts with Authentication & Data Layer
func (l *LogicalLayerService) ProcessData(request string, response *string) error {
	fmt.Println("Logical Layer → Authenticating user before processing data...")

	// Authenticate user
	authResp, err := callAuthService("user123")
	if err != nil {
		return fmt.Errorf("Authentication failed: %v", err)
	}
	fmt.Println("Authentication Successful:", authResp)

	// Connect to Data Layer
	fmt.Println("Logical Layer → Connecting to Data Layer Service...")
	client, err := rpc.Dial("tcp", "localhost:5009")
	if err != nil {
		return fmt.Errorf("Error connecting to Data Layer Service: %v", err)
	}
	defer client.Close()
	fmt.Println("Logical Layer → Connected to Data Layer Service.")

	// Store Data
	err = client.Call("DataLayerService.StoreData", request, response)
	if err != nil {
		return fmt.Errorf("Error calling DataLayerService.StoreData: %v", err)
	}
	fmt.Println("Data Successfully Stored:", *response)

	return nil
}

// ExecuteTrade interacts with B2C2 Service to execute a trade
func (l *LogicalLayerService) ExecuteTrade(request string, response *string) error {
	fmt.Println("Logical Layer → Authenticating user before trade execution...")

	// Authenticate user
	authResp, err := callAuthService("user123")
	if err != nil {
		return fmt.Errorf("Authentication failed: %v", err)
	}
	fmt.Println("Authentication Successful:", authResp)

	// Connect to B2C2 Service
	fmt.Println("Logical Layer → Connecting to B2C2 Service for trade execution...")
	client, err := rpc.Dial("tcp", "localhost:5007")
	if err != nil {
		return fmt.Errorf("Error connecting to B2C2 Service: %v", err)
	}
	defer client.Close()
	fmt.Println("Logical Layer → Connected to B2C2 Service.")

	// Execute Trade
	err = client.Call("B2C2Service.ExecuteTrade", request, response)
	if err != nil {
		return fmt.Errorf("Error calling B2C2Service.ExecuteTrade: %v", err)
	}
	fmt.Println("Trade Executed:", *response)
	return nil
}

// callAuthService verifies user authentication
func callAuthService(userID string) (string, error) {
	fmt.Println("Logical Layer → Connecting to Authentication Service...")

	client, err := rpc.Dial("tcp", "localhost:5001") // Connect to Authentication Service
	if err != nil {
		return "", fmt.Errorf("Error connecting to Auth Service: %v", err)
	}
	defer client.Close()
	fmt.Println("Logical Layer → Connected to Authentication Service.")

	var response string
	err = client.Call("AuthService.Authenticate", userID, &response)
	if err != nil {
		return "", fmt.Errorf("Error calling AuthService: %v", err)
	}
	return response, nil
}

// startLogicalLayerService initializes and runs the Logical Layer Service
func startLogicalLayerService() {
	service := new(LogicalLayerService)
	rpc.Register(service)

	listener, err := net.Listen("tcp", ":5008")
	if err != nil {
		fmt.Println("Error starting Logical Layer Service:", err)
		return
	}
	fmt.Println("Logical Layer Service is running on port 5008...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}

// Call Data Layer Service to store data
func callDataLayerService() {
	time.Sleep(3 * time.Second) // Wait for services to start
	fmt.Println("Logical Layer → Calling Data Layer Service to store data...")

	var response string
	logicalService := new(LogicalLayerService)
	err := logicalService.ProcessData("Sample Data", &response)
	if err != nil {
		fmt.Println("Error Processing Data:", err)
	} else {
		fmt.Println("Logical Layer Process Completed:", response)
	}
}

// Call B2C2 Service to execute a trade
func callB2C2Service() {
	time.Sleep(2 * time.Second) // Wait for services to start
	fmt.Println("Logical Layer → Calling B2C2 Service for trade execution...")

	var response string
	logicalService := new(LogicalLayerService)
	err := logicalService.ExecuteTrade("Buy BTC 1.5", &response)
	if err != nil {
		fmt.Println("Error Executing Trade:", err)
	} else {
		fmt.Println("Trade Execution Completed:", response)
	}
}

func main() {
	go startLogicalLayerService() // Start Logical Layer Service
	time.Sleep(2 * time.Second)

	// Test connections
	callDataLayerService() // Store Data
	callB2C2Service()      // Execute Trade

	select {} // Keep the service running
}
