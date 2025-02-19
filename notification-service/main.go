package main

import (
	"fmt"
	"net"
	"net/rpc"
)

// Define Notification Service
type NotificationService struct{}

// SendNotification method
func (n *NotificationService) SendNotification(request string, response *string) error {
	*response = "Notification sent: " + request
	return nil
}

func startNotificationService() {
	service := new(NotificationService)
	rpc.Register(service)

	listener, err := net.Listen("tcp", ":5003")
	if err != nil {
		fmt.Println("Error starting Notification Service:", err)
		return
	}
	fmt.Println("Notification Service running on port 5003")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}

func main() {
	startNotificationService()
}
