package main

import (
	"fmt"
	"net"
	"net/rpc"
	"time"
)

// Define Webhook Service
type WebhookService struct{}

// TriggerWebhook notifies an event and forwards it to Logical Layer
func (w *WebhookService) TriggerWebhook(request string, response *string) error {
	fmt.Println("Webhook Service → Triggering Webhook Event...")

	// Call Logical Layer Service for processing
	logicalResp, err := callLogicalLayerService(request)
	if err != nil {
		return fmt.Errorf("Error processing event in Logical Layer: %v", err)
	}
	fmt.Println("Webhook Service → Processed by Logical Layer:", logicalResp)

	*response = "Webhook triggered: " + request
	return nil
}

// Start Webhook Service
func startWebhookService() {
	service := new(WebhookService)
	rpc.Register(service)

	listener, err := net.Listen("tcp", ":5002")
	if err != nil {
		fmt.Println("Error starting Webhook Service:", err)
		return
	}
	fmt.Println("Webhook Service running on port 5002")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}

// Call Notification Service (port 5003)
func callNotificationService() {
	time.Sleep(3 * time.Second) // Wait for Notification Service to start
	fmt.Println("Webhook Service → Connecting to Notification Service...")

	client, err := rpc.Dial("tcp", "localhost:5003") // Connect to Notification Service
	if err != nil {
		fmt.Println("Error connecting to Notification Service:", err)
		return
	}
	defer client.Close()
	fmt.Println("Webhook Service → Connected to Notification Service.")

	var response string
	err = client.Call("NotificationService.SendNotification", "Webhook Event", &response)
	if err != nil {
		fmt.Println("Error calling Notification Service:", err)
	} else {
		fmt.Println("Webhook Service → Response from Notification Service:", response)
	}
}

// Call Logical Layer Service (port 5008)
func callLogicalLayerService(eventData string) (string, error) {
	time.Sleep(2 * time.Second) // Wait for Logical Layer Service to start
	fmt.Println("Webhook Service → Connecting to Logical Layer Service...")

	client, err := rpc.Dial("tcp", "localhost:5008") // Connect to Logical Layer Service
	if err != nil {
		return "", fmt.Errorf("Error connecting to Logical Layer Service: %v", err)
	}
	defer client.Close()
	fmt.Println("Webhook Service → Connected to Logical Layer Service.")

	var response string
	err = client.Call("LogicalLayerService.ProcessData", eventData, &response)
	if err != nil {
		return "", fmt.Errorf("Error calling LogicalLayerService.ProcessData: %v", err)
	}
	return response, nil
}

func main() {
	go startWebhookService() // Start Webhook Service in a goroutine
	time.Sleep(2 * time.Second)

	// Call test connections
	callNotificationService()   // Notify external system
	callLogicalLayerService("Sample Webhook Event") // Process Webhook through Logical Layer

	select {} // Keep service running
}
