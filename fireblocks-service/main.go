package main

import (
	"fmt"
	"net"
	"net/rpc"
)

// Fireblocks Service struct
type FireblocksService struct{}

// ProcessTransaction method
func (f *FireblocksService) ProcessTransaction(request string, response *string) error {
	*response = "Transaction processed: " + request
	return nil
}

func startFireblocksService() {
	service := new(FireblocksService)
	rpc.Register(service)

	listener, err := net.Listen("tcp", ":5005")
	if err != nil {
		fmt.Println("Error starting Fireblocks Service:", err)
		return
	}
	fmt.Println("Fireblocks Service running on port 5005")

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
	startFireblocksService()
}
