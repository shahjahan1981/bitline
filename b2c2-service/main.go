package main

import (
	"fmt"
	"net"
	"net/rpc"
)

// Data Layer Service struct
type DataLayerService struct{}

// StoreData method
func (d *DataLayerService) StoreData(request string, response *string) error {
	*response = "Data stored: " + request
	return nil
}

func startDataLayerService() {
	service := new(DataLayerService)
	rpc.Register(service)

	listener, err := net.Listen("tcp", ":5009")
	if err != nil {
		fmt.Println("Error starting Data Layer Service:", err)
		return
	}
	fmt.Println("Data Layer Service running on port 5009")

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
	startDataLayerService()
}
